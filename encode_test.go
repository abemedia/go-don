package don_test

import (
	"context"
	"io"
	"testing"

	"github.com/abemedia/go-don"
	"github.com/abemedia/go-don/pkg/httptest"
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

func testRegisterEncoder[T don.EncoderConstraint](t *testing.T, dec T, contentType, alias string) {
	t.Helper()

	don.RegisterEncoder(contentType, dec, alias)

	for _, v := range []string{contentType, alias} {
		encode, err := don.GetEncoder(v)
		if err != nil {
			t.Error(err)
			continue
		}

		req := httptest.NewRequest("", "", v, nil)

		if err = encode(req, []byte(v)); err != nil {
			t.Error(err)
		} else if string(req.Response.Body()) != v {
			t.Error("should encode response")
		}

		if err = encode(req, []byte{}); err == nil {
			t.Error("should return error")
		}
	}
}
