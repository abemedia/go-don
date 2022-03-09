package don

import (
	"context"
	"io"
	"net/http"
)

type (
	Marshaler        = func(v interface{}) ([]byte, error)
	ContextMarshaler = func(ctx context.Context, v interface{}) ([]byte, error)
	EncoderFactory   = func(io.Writer) interface{ Encode(interface{}) error }
	ResponseEncoder  = func(w http.ResponseWriter, v interface{}) error
)

type EncoderConstraint interface {
	Marshaler | ContextMarshaler | EncoderFactory | ResponseEncoder
}

// RegisterEncoder registers a response encoder.
func RegisterEncoder[T EncoderConstraint](contentType string, enc T, aliases ...string) {
	switch e := any(enc).(type) {
	case Marshaler:
		encoders[contentType] = func(w http.ResponseWriter, v interface{}) error {
			b, err := e(v)
			if err != nil {
				return err
			}
			_, err = w.Write(b)
			return err
		}

	case ContextMarshaler:
		encoders[contentType] = func(w http.ResponseWriter, v interface{}) error {
			b, err := e(context.TODO(), v)
			if err != nil {
				return err
			}
			_, err = w.Write(b)
			return err
		}

	case EncoderFactory:
		encoders[contentType] = func(w http.ResponseWriter, v interface{}) error {
			return e(w).Encode(v)
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
