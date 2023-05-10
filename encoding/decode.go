package encoding

import (
	"context"

	"github.com/valyala/fasthttp"
)

type (
	Unmarshaler        = func(data []byte, v any) error
	ContextUnmarshaler = func(ctx context.Context, data []byte, v any) error
	RequestParser      = func(ctx *fasthttp.RequestCtx, v any) error
)

type DecoderConstraint interface {
	Unmarshaler | ContextUnmarshaler | RequestParser
}

// RegisterDecoder registers a request decoder for a given media type.
func RegisterDecoder[T DecoderConstraint](dec T, mime string, aliases ...string) {
	switch d := any(dec).(type) {
	case Unmarshaler:
		decoders[mime] = func(ctx *fasthttp.RequestCtx, v any) error {
			return d(ctx.Request.Body(), v)
		}

	case ContextUnmarshaler:
		decoders[mime] = func(ctx *fasthttp.RequestCtx, v any) error {
			return d(ctx, ctx.Request.Body(), v)
		}

	case RequestParser:
		decoders[mime] = d
	}

	for _, alias := range aliases {
		decoders[alias] = decoders[mime]
	}
}

// GetDecoder returns the request decoder for a given media type.
func GetDecoder(mime string) RequestParser {
	return decoders[mime]
}

var decoders = map[string]RequestParser{}
