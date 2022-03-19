package don

import (
	"bytes"
	"encoding"
	"encoding/xml"
	"net/http"
	"strconv"

	"github.com/goccy/go-json"
	"gopkg.in/yaml.v2"
)

// StatusError creates an error from an HTTP status code.
type StatusError int

const (
	ErrBadRequest           = StatusError(http.StatusBadRequest)
	ErrUnauthorized         = StatusError(http.StatusUnauthorized)
	ErrForbidden            = StatusError(http.StatusForbidden)
	ErrNotFound             = StatusError(http.StatusNotFound)
	ErrMethodNotAllowed     = StatusError(http.StatusMethodNotAllowed)
	ErrNotAcceptable        = StatusError(http.StatusNotAcceptable)
	ErrUnsupportedMediaType = StatusError(http.StatusUnsupportedMediaType)
	ErrInternalServerError  = StatusError(http.StatusInternalServerError)
)

func (e StatusError) Error() string {
	return http.StatusText(int(e))
}

func (e StatusError) StatusCode() int {
	return int(e)
}

func (e StatusError) MarshalText() ([]byte, error) {
	return []byte(e.Error()), nil
}

func (e StatusError) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString(`{"message":`)
	buf.WriteString(strconv.Quote(e.Error()))
	buf.WriteRune('}')

	return buf.Bytes(), nil
}

func (e StatusError) MarshalXML(enc *xml.Encoder, _ xml.StartElement) error {
	start := xml.StartElement{Name: xml.Name{Local: "message"}}
	return enc.EncodeElement(e.Error(), start)
}

func (e StatusError) MarshalYAML() (interface{}, error) {
	return map[string]string{"message": e.Error()}, nil
}

var (
	_ error                  = (*StatusError)(nil)
	_ encoding.TextMarshaler = (*StatusError)(nil)
	_ json.Marshaler         = (*StatusError)(nil)
	_ xml.Marshaler          = (*StatusError)(nil)
	_ yaml.Marshaler         = (*StatusError)(nil)
)
