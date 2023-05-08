package don

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

// RegisterDecoder registers a request decoder.
func RegisterDecoder[T DecoderConstraint](contentType string, dec T, aliases ...string) {
	switch d := any(dec).(type) {
	case Unmarshaler:
		decoders[contentType] = func(ctx *fasthttp.RequestCtx, v any) error {
			return d(ctx.Request.Body(), v)
		}

	case ContextUnmarshaler:
		decoders[contentType] = func(ctx *fasthttp.RequestCtx, v any) error {
			return d(ctx, ctx.Request.Body(), v)
		}

	case RequestParser:
		decoders[contentType] = d
	}

	for _, alias := range aliases {
		decoders[alias] = decoders[contentType]
	}
}

func getDecoder(mime string) (RequestParser, error) {
	if enc := decoders[mime]; enc != nil {
		return enc, nil
	}
	return nil, ErrUnsupportedMediaType
}

var decoders = map[string]RequestParser{}
