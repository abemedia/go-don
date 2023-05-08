package don

import (
	"context"

	"github.com/valyala/fasthttp"
)

type (
	Marshaler        = func(v any) ([]byte, error)
	ContextMarshaler = func(ctx context.Context, v any) ([]byte, error)
	ResponseEncoder  = func(ctx *fasthttp.RequestCtx, v any) error
)

type EncoderConstraint interface {
	Marshaler | ContextMarshaler | ResponseEncoder
}

// RegisterEncoder registers a response encoder.
func RegisterEncoder[T EncoderConstraint](contentType string, enc T, aliases ...string) {
	switch e := any(enc).(type) {
	case Marshaler:
		encoders[contentType] = func(ctx *fasthttp.RequestCtx, v any) error {
			b, err := e(v)
			if err != nil {
				return err
			}
			ctx.Response.SetBodyRaw(b)
			return nil
		}

	case ContextMarshaler:
		encoders[contentType] = func(ctx *fasthttp.RequestCtx, v any) error {
			b, err := e(ctx, v)
			if err != nil {
				return err
			}
			ctx.Response.SetBodyRaw(b)
			return nil
		}

	case ResponseEncoder:
		encoders[contentType] = e
	}

	for _, alias := range aliases {
		encoders[alias] = encoders[contentType]
	}
}

func getEncoder(mime string) (ResponseEncoder, error) {
	if enc := encoders[mime]; enc != nil {
		return enc, nil
	}
	return nil, ErrNotAcceptable
}

var encoders = map[string]ResponseEncoder{}
