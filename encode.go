package don

import (
	"context"

	"github.com/valyala/fasthttp"
)

type (
	Marshaler        = func(v interface{}) ([]byte, error)
	ContextMarshaler = func(ctx context.Context, v interface{}) ([]byte, error)
	ResponseEncoder  = func(ctx *fasthttp.RequestCtx, v interface{}) error
)

type EncoderConstraint interface {
	Marshaler | ContextMarshaler | ResponseEncoder
}

// RegisterEncoder registers a response encoder.
func RegisterEncoder[T EncoderConstraint](contentType string, enc T, aliases ...string) {
	switch e := any(enc).(type) {
	case Marshaler:
		encoders[contentType] = func(ctx *fasthttp.RequestCtx, v interface{}) error {
			b, err := e(v)
			if err != nil {
				return err
			}

			_, err = ctx.Write(b)
			return err
		}

	case ContextMarshaler:
		encoders[contentType] = func(ctx *fasthttp.RequestCtx, v interface{}) error {
			b, err := e(ctx, v)
			if err != nil {
				return err
			}

			_, err = ctx.Write(b)
			return err
		}

	case ResponseEncoder:
		encoders[contentType] = e
	}

	for _, alias := range aliases {
		encoderAliases[alias] = contentType
	}
}

func getEncoder(mime string) (ResponseEncoder, error) {
	if enc := encoders[mime]; enc != nil {
		return enc, nil
	}

	if name := encoderAliases[mime]; name != "" {
		return encoders[name], nil
	}

	return nil, ErrNotAcceptable
}

var (
	encoders       = map[string]ResponseEncoder{}
	encoderAliases = map[string]string{}
)
