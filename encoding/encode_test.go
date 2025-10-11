package encoding_test

import (
	"context"
	"io"
	"testing"

	"github.com/abemedia/go-don/encoding"
	"github.com/abemedia/go-don/internal/testutil"
	"github.com/valyala/fasthttp"
)

func TestRegisterEncoder(t *testing.T) {
	t.Run("Marshaler", func(t *testing.T) {
		testRegisterEncoder(t, func(v any) ([]byte, error) {
			b := v.([]byte)
			if len(b) == 0 {
				return nil, io.EOF
			}
			return b, nil
		}, "unmarshaler", "marshaler-alias")
	})

	t.Run("ContextMarshaler", func(t *testing.T) {
		testRegisterEncoder(t, func(ctx context.Context, v any) ([]byte, error) {
			b := v.([]byte)
			if len(b) == 0 {
				return nil, io.EOF
			}
			return b, nil
		}, "context-marshaler", "context-marshaler-alias")
	})

	t.Run("ResponseEncoder", func(t *testing.T) {
		testRegisterEncoder(t, func(ctx *fasthttp.RequestCtx, v any) error {
			b := v.([]byte)
			if len(b) == 0 {
				return io.EOF
			}
			ctx.Response.SetBodyRaw(b)
			return nil
		}, "response-encoder", "response-encoder-alias")
	})
}

func testRegisterEncoder[T encoding.EncoderConstraint](t *testing.T, dec T, contentType, alias string) {
	t.Helper()

	encoding.RegisterEncoder(dec, contentType, alias)

	for _, v := range []string{contentType, alias} {
		encode := encoding.GetEncoder(v)
		if encode == nil {
			t.Error("encoder not found")
			continue
		}

		req := testutil.NewRequest("", "", v, nil)

		if err := encode(req, []byte(v)); err != nil {
			t.Error(err)
		} else if string(req.Response.Body()) != v {
			t.Error("should encode response")
		}

		if err := encode(req, []byte{}); err == nil {
			t.Error("should return error")
		}
	}
}

func TestGetEncoderMultipleContentTypes(t *testing.T) {
	encFn := func(ctx *fasthttp.RequestCtx, v any) error {
		return nil
	}

	encoding.RegisterEncoder(encFn, "application/xml")

	enc := encoding.GetEncoder("text/html,application/xhtml+xml,application/xml")
	if enc == nil {
		t.Fatal("encoder not found")
	}

	enc = encoding.GetEncoder("application/xhtml+xml")
	if enc != nil {
		t.Fatal("encoder should not be found")
	}
}
