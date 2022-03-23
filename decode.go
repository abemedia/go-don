package don

import (
	"context"

	"github.com/valyala/fasthttp"
)

type (
	Unmarshaler        = func(data []byte, v interface{}) error
	ContextUnmarshaler = func(ctx context.Context, data []byte, v interface{}) error
	RequestParser      = func(ctx *fasthttp.RequestCtx, v interface{}) error
)

type DecoderConstraint interface {
	Unmarshaler | ContextUnmarshaler | RequestParser
}

// RegisterDecoder registers a request decoder.
func RegisterDecoder[T DecoderConstraint](contentType string, dec T, aliases ...string) {
	switch d := any(dec).(type) {
	case Unmarshaler:
		decoders[contentType] = func(ctx *fasthttp.RequestCtx, v interface{}) error {
			return d(ctx.Request.Body(), v)
		}

	case ContextUnmarshaler:
		decoders[contentType] = func(ctx *fasthttp.RequestCtx, v interface{}) error {
			return d(ctx, ctx.Request.Body(), v)
		}

	case RequestParser:
		decoders[contentType] = d
	}

	for _, alias := range aliases {
		decoderAliases[alias] = contentType
	}
}

func getDecoder(mime string) (RequestParser, error) {
	if enc := decoders[mime]; enc != nil {
		return enc, nil
	}

	if name := decoderAliases[mime]; name != "" {
		return decoders[name], nil
	}

	return nil, ErrUnsupportedMediaType
}

var (
	decoders       = map[string]RequestParser{}
	decoderAliases = map[string]string{}
)
