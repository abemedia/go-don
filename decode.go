package don

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
)

type (
	Unmarshaler        = func(data []byte, v interface{}) error
	ContextUnmarshaler = func(ctx context.Context, data []byte, v interface{}) error
	DecoderFactory     = func(io.Reader) interface{ Decode(interface{}) error }
	RequestParser      = func(r *http.Request, v interface{}) error
)

type DecoderConstraint interface {
	Unmarshaler | ContextUnmarshaler | DecoderFactory | RequestParser
}

func RegisterDecoder[T DecoderConstraint](contentType string, dec T, aliases ...string) {
	switch d := any(dec).(type) {
	case Unmarshaler:
		decoders[contentType] = func(r *http.Request, v interface{}) error {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return err
			}
			return d(b, v)
		}

	case ContextUnmarshaler:
		decoders[contentType] = func(r *http.Request, v interface{}) error {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return err
			}
			return d(r.Context(), b, v)
		}

	case DecoderFactory:
		decoders[contentType] = func(r *http.Request, v interface{}) error {
			return d(r.Body).Decode(v)
		}

	case RequestParser:
		decoders[contentType] = d
	}

	for _, alias := range aliases {
		decoderAliases[alias] = contentType
	}
}

func getDecoder(contentType, fallback string) (RequestParser, error) {
	mime := parseMime(contentType, fallback)

	if enc, ok := decoders[mime]; ok {
		return enc, nil
	}

	if name, ok := decoderAliases[mime]; ok {
		return decoders[name], nil
	}

	return nil, ErrUnsupportedMediaType
}

var (
	decoders       = map[string]RequestParser{}
	decoderAliases = map[string]string{}
)
