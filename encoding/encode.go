package encoding

import (
	"context"
	"strings"

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

// RegisterEncoder registers a response encoder on a given media type.
func RegisterEncoder[T EncoderConstraint](enc T, mime string, aliases ...string) {
	switch e := any(enc).(type) {
	case Marshaler:
		encoders[mime] = func(ctx *fasthttp.RequestCtx, v any) error {
			b, err := e(v)
			if err != nil {
				return err
			}
			ctx.Response.SetBodyRaw(b)
			return nil
		}

	case ContextMarshaler:
		encoders[mime] = func(ctx *fasthttp.RequestCtx, v any) error {
			b, err := e(ctx, v)
			if err != nil {
				return err
			}
			ctx.Response.SetBodyRaw(b)
			return nil
		}

	case ResponseEncoder:
		encoders[mime] = e
	}

	for _, alias := range aliases {
		encoders[alias] = encoders[mime]
	}
}

// GetEncoder returns the response encoder for a given media type.
func GetEncoder(mime string) ResponseEncoder {
	mimeParts := strings.SplitSeq(mime, ",")
	for part := range mimeParts {
		if enc, ok := encoders[part]; ok {
			return enc
		}
	}
	return nil
}

var encoders = map[string]ResponseEncoder{}
