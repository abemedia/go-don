package test

import (
	"context"
	"net/http"
	"testing"

	"github.com/abemedia/go-don"
	_ "github.com/abemedia/go-don/encoding/text" // default encoding
	"github.com/abemedia/go-don/pkg/httptest"
	"github.com/google/go-cmp/cmp"
)

type EncodingTest[T any] struct {
	Mime   string
	Parsed T
	Raw    string
}

func Encoding[T any](t *testing.T, test EncodingTest[T]) {
	t.Helper()
	t.Run("Encode", func(t *testing.T) {
		t.Helper()
		Encode(t, test)
	})
	t.Run("Decode", func(t *testing.T) {
		t.Helper()
		Decode(t, test)
	})
}

func Encode[T any](t *testing.T, test EncodingTest[T]) {
	t.Helper()

	api := don.New(nil)
	api.Post("/", don.H(func(ctx context.Context, req don.Empty) (T, error) {
		return test.Parsed, nil
	}))

	ctx := httptest.NewRequest(http.MethodPost, "/", "", map[string]string{"Accept": test.Mime})
	api.RequestHandler()(ctx)

	if diff := cmp.Diff(test.Raw, string(ctx.Response.Body())); diff != "" {
		t.Fatal(diff)
	}

	if ctx.Response.StatusCode() != http.StatusOK {
		t.Fatalf("expected success status: %v", &ctx.Response)
	}
}

func Decode[T any](t *testing.T, test EncodingTest[T]) {
	t.Helper()

	var got T

	api := don.New(nil)
	api.Post("/", don.H(func(ctx context.Context, req T) (don.Empty, error) {
		got = req
		return don.Empty{}, nil
	}))

	ctx := httptest.NewRequest(http.MethodPost, "/", test.Raw, map[string]string{"Content-Type": test.Mime})
	api.RequestHandler()(ctx)

	if diff := cmp.Diff(test.Parsed, got); diff != "" {
		t.Fatal(diff)
	}

	if ctx.Response.StatusCode() != http.StatusNoContent {
		t.Fatalf("expected success status: %v", &ctx.Response)
	}
}
