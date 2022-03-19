package don

import (
	"bytes"
	"context"
	"encoding"
	"encoding/xml"
	"net/http"
	"strconv"

	"github.com/goccy/go-json"
	"gopkg.in/yaml.v2"
)

func E(err error) Handler {
	return H(func(context.Context, *Empty) (*Empty, error) {
		return nil, err
	})
}

type HTTPError struct {
	err  error
	code int
}

func Error(err error, code int) *HTTPError {
	return &HTTPError{err, code}
}

func (e *HTTPError) Error() string {
	return e.err.Error()
}

func (e *HTTPError) IsPrivate() bool {
	return e.StatusCode() == http.StatusInternalServerError
}

func (e *HTTPError) StatusCode() int {
	if sc, ok := e.err.(StatusCoder); ok {
		return sc.StatusCode()
	}

	if e.code == 0 {
		return http.StatusInternalServerError
	}

	return e.code
}

func (e *HTTPError) MarshalText() ([]byte, error) {
	return []byte(e.Error()), nil
}

func (e *HTTPError) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString(`{"message":`)
	buf.WriteString(strconv.Quote(e.Error()))
	buf.WriteRune('}')

	return buf.Bytes(), nil
}

func (e *HTTPError) MarshalXML(enc *xml.Encoder, _ xml.StartElement) error {
	start := xml.StartElement{Name: xml.Name{Local: "message"}}
	return enc.EncodeElement(e.Error(), start)
}

func (e *HTTPError) MarshalYAML() (interface{}, error) {
	return map[string]string{"message": e.Error()}, nil
}

var (
	_ error                  = (*HTTPError)(nil)
	_ encoding.TextMarshaler = (*HTTPError)(nil)
	_ json.Marshaler         = (*HTTPError)(nil)
	_ xml.Marshaler          = (*HTTPError)(nil)
	_ yaml.Marshaler         = (*HTTPError)(nil)
)
