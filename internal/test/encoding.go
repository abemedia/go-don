package test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/abemedia/go-don"
	"github.com/abemedia/go-don/encoding"
	_ "github.com/abemedia/go-don/encoding/text" // default encoding
	"github.com/abemedia/go-don/pkg/httptest"
	"github.com/google/go-cmp/cmp"
)

type EncodingOptions[T any] struct {
	Mime   string
	Parsed T
	Raw    string
}

func Encoding[T any](t *testing.T, opt EncodingOptions[T]) {
	t.Helper()
	t.Run("Encode", func(t *testing.T) {
		t.Helper()
		Encode(t, opt)
	})
	t.Run("Decode", func(t *testing.T) {
		t.Helper()
		Decode(t, opt)
	})
}

func Encode[T any](t *testing.T, opt EncodingOptions[T]) {
	t.Helper()

	api := don.New(nil)
	api.Post("/", don.H(func(ctx context.Context, req don.Empty) (T, error) {
		return opt.Parsed, nil
	}))

	ctx := httptest.NewRequest(http.MethodPost, "/", "", map[string]string{"Accept": opt.Mime})
	api.RequestHandler()(ctx)

	if diff := cmp.Diff(opt.Raw, string(ctx.Response.Body())); diff != "" {
		t.Fatal(diff)
	}

	if ctx.Response.StatusCode() != http.StatusOK {
		t.Fatalf("expected success status: %v", &ctx.Response)
	}
}

func Decode[T any](t *testing.T, opt EncodingOptions[T]) {
	t.Helper()

	var got T

	api := don.New(nil)
	api.Post("/", don.H(func(ctx context.Context, req T) (don.Empty, error) {
		got = req
		return don.Empty{}, nil
	}))

	ctx := httptest.NewRequest(http.MethodPost, "/", opt.Raw, map[string]string{"Content-Type": opt.Mime})
	api.RequestHandler()(ctx)

	if diff := cmp.Diff(opt.Parsed, got); diff != "" {
		t.Fatal(diff)
	}

	if ctx.Response.StatusCode() != http.StatusNoContent {
		t.Fatalf("expected success status: %v", &ctx.Response)
	}
}

func BenchmarkEncoding[T any](b *testing.B, opt EncodingOptions[T]) {
	b.Run("Encode", func(b *testing.B) {
		BenchmarkEncode(b, opt)
	})
	b.Run("Decode", func(b *testing.B) {
		BenchmarkDecode(b, opt)
	})
}

func BenchmarkEncode[T any](b *testing.B, opt EncodingOptions[T]) {
	enc := encoding.GetEncoder(opt.Mime)
	if enc == nil {
		b.Fatal("encoder not found")
	}

	ctx := httptest.NewRequest("POST", "/", "", nil)

	for i := 0; i < b.N; i++ {
		ctx.Response.ResetBody()
		enc(ctx, opt.Parsed) //nolint:errcheck
	}
}

func BenchmarkDecode[T any](b *testing.B, opt EncodingOptions[T]) {
	dec := encoding.GetDecoder(opt.Mime)
	if dec == nil {
		b.Fatal("decoder not found")
	}

	rd := strings.NewReader(opt.Raw)
	ctx := httptest.NewRequest("POST", "/", "", nil)
	ctx.Request.SetBodyStream(rd, len(opt.Raw))

	for i := 0; i < b.N; i++ {
		rd.Seek(0, io.SeekStart) //nolint:errcheck
		dec(ctx, new(T))         //nolint:errcheck
	}
}
