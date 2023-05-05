package don_test

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/abemedia/go-don"
	_ "github.com/abemedia/go-don/encoding/json"
	_ "github.com/abemedia/go-don/encoding/text"
	_ "github.com/abemedia/go-don/encoding/xml"
	_ "github.com/abemedia/go-don/encoding/yaml"
	"github.com/abemedia/go-don/internal/byteconv"
	"github.com/abemedia/go-don/pkg/httptest"
	"github.com/google/go-cmp/cmp"
)

func TestError(t *testing.T) {
	tests := []struct {
		err  error
		mime string
		body string
		code int
	}{
		{
			err:  don.ErrBadRequest,
			mime: "text/plain; charset=utf-8",
			body: "Bad Request\n",
			code: http.StatusBadRequest,
		},
		{
			err:  don.Error(errors.New("test"), http.StatusBadRequest),
			mime: "text/plain; charset=utf-8",
			body: "test\n",
			code: http.StatusBadRequest,
		},
		{
			err:  errors.New("test"),
			mime: "text/plain; charset=utf-8",
			body: "test\n",
			code: http.StatusInternalServerError,
		},
		{
			err:  errors.New("test"),
			mime: "application/json; charset=utf-8",
			body: `{"message":"test"}` + "\n",
			code: http.StatusInternalServerError,
		},
		{
			err:  errors.New("test"),
			mime: "application/x-yaml; charset=utf-8",
			body: "message: test\n",
			code: http.StatusInternalServerError,
		},
		{
			err:  errors.New("test"),
			mime: "application/xml; charset=utf-8",
			body: "<message>test</message>",
			code: http.StatusInternalServerError,
		},
		{
			err:  &testError{},
			mime: "text/plain; charset=utf-8",
			body: "test\n",
			code: http.StatusInternalServerError,
		},
		{
			err:  &testError{},
			mime: "application/json; charset=utf-8",
			body: `{"custom":"test"}` + "\n",
			code: http.StatusInternalServerError,
		},
		{
			err:  &testError{},
			mime: "application/x-yaml; charset=utf-8",
			body: "custom: test\n",
			code: http.StatusInternalServerError,
		},
		{
			err:  &testError{},
			mime: "application/xml; charset=utf-8",
			body: "<custom>test</custom>",
			code: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		ctx := httptest.NewRequest("", "/", "", map[string]string{"Accept": tc.mime})

		api := don.New(nil)
		api.Get("/", don.H(func(ctx context.Context, req don.Empty) (any, error) { return nil, tc.err }))
		api.RequestHandler()(ctx)

		type response struct {
			Code   int
			Body   string
			Header map[string]string
		}

		res := response{ctx.Response.StatusCode(), string(ctx.Response.Body()), map[string]string{}}
		ctx.Response.Header.VisitAll(func(key, value []byte) { res.Header[string(key)] = string(value) })

		want := response{tc.code, tc.body, map[string]string{"Content-Type": tc.mime}}
		if diff := cmp.Diff(want, res); diff != "" {
			t.Error(diff)
		}
	}
}

type testError struct{}

func (e *testError) Error() string {
	return "test"
}

func (e *testError) MarshalText() ([]byte, error) {
	return byteconv.Atob(e.Error()), nil
}

func (e *testError) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"custom":%q}`, e.Error())), nil
}

func (e *testError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Local: "custom"}
	return enc.EncodeElement(e.Error(), start)
}

func (e *testError) MarshalYAML() (any, error) {
	return map[string]string{"custom": e.Error()}, nil
}
