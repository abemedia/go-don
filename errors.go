package don

import (
	"bytes"
	"context"
	"encoding"
	"encoding/xml"
	"errors"
	"net/http"
	"strconv"

	"github.com/abemedia/go-don/internal/byteconv"
	"github.com/goccy/go-json"
	"github.com/valyala/fasthttp"
	"gopkg.in/yaml.v3"
)

func E(err error) fasthttp.RequestHandler {
	h := H(func(context.Context, Empty) (*Empty, error) { return nil, err })
	return func(ctx *fasthttp.RequestCtx) { h(ctx, nil) }
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

func (e *HTTPError) StatusCode() int {
	var sc StatusCoder
	if errors.As(e.err, &sc) {
		return sc.StatusCode()
	}

	if e.code == 0 {
		return http.StatusInternalServerError
	}

	return e.code
}

func (e *HTTPError) MarshalText() ([]byte, error) {
	var m encoding.TextMarshaler
	if errors.As(e.err, &m) {
		return m.MarshalText()
	}

	return byteconv.Atob(e.Error()), nil
}

func (e *HTTPError) MarshalJSON() ([]byte, error) {
	var m json.Marshaler
	if errors.As(e.err, &m) {
		return m.MarshalJSON()
	}

	var buf bytes.Buffer
	buf.WriteString(`{"message":`)
	buf.WriteString(strconv.Quote(e.Error()))
	buf.WriteRune('}')

	return buf.Bytes(), nil
}

func (e *HTTPError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	var m xml.Marshaler
	if errors.As(e.err, &m) {
		return m.MarshalXML(enc, start)
	}

	start = xml.StartElement{Name: xml.Name{Local: "message"}}
	return enc.EncodeElement(e.Error(), start)
}

func (e *HTTPError) MarshalYAML() (any, error) {
	var m yaml.Marshaler
	if errors.As(e.err, &m) {
		return m.MarshalYAML()
	}

	return map[string]string{"message": e.Error()}, nil
}

var (
	_ error                  = (*HTTPError)(nil)
	_ encoding.TextMarshaler = (*HTTPError)(nil)
	_ json.Marshaler         = (*HTTPError)(nil)
	_ xml.Marshaler          = (*HTTPError)(nil)
	_ yaml.Marshaler         = (*HTTPError)(nil)
)

// StatusError creates an error from an HTTP status code.
type StatusError int

const (
	ErrBadRequest           = StatusError(fasthttp.StatusBadRequest)
	ErrUnauthorized         = StatusError(fasthttp.StatusUnauthorized)
	ErrForbidden            = StatusError(fasthttp.StatusForbidden)
	ErrNotFound             = StatusError(fasthttp.StatusNotFound)
	ErrMethodNotAllowed     = StatusError(fasthttp.StatusMethodNotAllowed)
	ErrNotAcceptable        = StatusError(fasthttp.StatusNotAcceptable)
	ErrUnsupportedMediaType = StatusError(fasthttp.StatusUnsupportedMediaType)
	ErrInternalServerError  = StatusError(fasthttp.StatusInternalServerError)
)

func (e StatusError) Error() string {
	return fasthttp.StatusMessage(int(e))
}

func (e StatusError) StatusCode() int {
	return int(e)
}

func (e StatusError) MarshalText() ([]byte, error) {
	return byteconv.Atob(e.Error()), nil
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

func (e StatusError) MarshalYAML() (any, error) {
	return map[string]string{"message": e.Error()}, nil
}

var (
	_ error                  = (*StatusError)(nil)
	_ encoding.TextMarshaler = (*StatusError)(nil)
	_ json.Marshaler         = (*StatusError)(nil)
	_ xml.Marshaler          = (*StatusError)(nil)
	_ yaml.Marshaler         = (*StatusError)(nil)
)
