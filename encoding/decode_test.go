package encoding_test

import (
	"context"
	"io"
	"reflect"
	"testing"

	"github.com/abemedia/go-don/encoding"
	"github.com/abemedia/go-don/internal/testutil"
	"github.com/valyala/fasthttp"
)

func TestRegisterDecoder(t *testing.T) {
	t.Run("Unmarshaler", func(t *testing.T) {
		testRegisterDecoder(t, func(data []byte, v any) error {
			if len(data) == 0 {
				return io.EOF
			}
			reflect.ValueOf(v).Elem().SetBytes(data)
			return nil
		}, "unmarshaler", "unmarshaler-alias")
	})

	t.Run("ContextUnmarshaler", func(t *testing.T) {
		testRegisterDecoder(t, func(ctx context.Context, data []byte, v any) error {
			if len(data) == 0 {
				return io.EOF
			}
			reflect.ValueOf(v).Elem().SetBytes(data)
			return nil
		}, "context-unmarshaler", "context-unmarshaler-alias")
	})

	t.Run("RequestParser", func(t *testing.T) {
		testRegisterDecoder(t, func(ctx *fasthttp.RequestCtx, v any) error {
			b := ctx.Request.Body()
			if len(b) == 0 {
				return io.EOF
			}
			reflect.ValueOf(v).Elem().SetBytes(b)
			return nil
		}, "request-parser", "request-parser-alias")
	})
}

func testRegisterDecoder[T encoding.DecoderConstraint](t *testing.T, dec T, contentType, alias string) {
	t.Helper()

	encoding.RegisterDecoder(dec, contentType, alias)

	for _, v := range []string{contentType, alias} {
		decode := encoding.GetDecoder(v)
		if decode == nil {
			t.Error("decoder not found")
			continue
		}

		req := testutil.NewRequest("", "", v, nil)

		var b []byte
		if err := decode(req, &b); err != nil {
			t.Error(err)
		} else if string(b) != v {
			t.Error("should decode request")
		}

		req.Request.SetBody(nil)
		if err := decode(req, &b); err == nil {
			t.Error("should return error")
		}
	}
}
