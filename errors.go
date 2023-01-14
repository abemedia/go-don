package don

import (
	"bytes"
	"context"
	"encoding"
	"encoding/xml"
	"errors"
	"net/http"
	"strconv"

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

	return []byte(e.Error()), nil
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

func (e *HTTPError) MarshalYAML() (interface{}, error) {
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
