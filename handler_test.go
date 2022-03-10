package don_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/abemedia/go-don"
	_ "github.com/abemedia/go-don/encoding/json"
	_ "github.com/abemedia/go-don/encoding/text"
	"github.com/google/go-cmp/cmp"
)

func TestHandler(t *testing.T) {
	type response struct {
		Code   int
		Body   string
		Header http.Header
	}

	tests := []struct {
		message  string
		expected response
		config   *don.Config
		handler  http.Handler
		body     string
		header   http.Header
	}{
		{
			message: "should return no content",
			expected: response{
				Code:   http.StatusNoContent,
				Body:   "",
				Header: http.Header{"Content-Type": {"text/plain; charset=utf-8"}},
			},
			config: &don.Config{DefaultEncoding: "text/plain"},
			handler: don.H(func(ctx context.Context, req *don.Empty) (any, error) {
				return nil, nil
			}),
		},
		{
			message: "should return null",
			expected: response{
				Code:   http.StatusOK,
				Body:   "null\n",
				Header: http.Header{"Content-Type": {"application/json; charset=utf-8"}},
			},
			config: &don.Config{DefaultEncoding: "application/json", DisableNoContent: true},
			handler: don.H(func(ctx context.Context, req *don.Empty) (any, error) {
				return nil, nil
			}),
		},
		{
			message: "should JSON encode map",
			expected: response{
				Code:   http.StatusOK,
				Body:   `{"foo":"bar"}` + "\n",
				Header: http.Header{"Content-Type": {"application/json; charset=utf-8"}},
			},
			config: &don.Config{DefaultEncoding: "application/json"},
			handler: don.H(func(ctx context.Context, req *don.Empty) (map[string]string, error) {
				return map[string]string{"foo": "bar"}, nil
			}),
		},
		{
			message: "should return error on unprocessable request",
			expected: response{
				Code:   http.StatusBadRequest,
				Body:   "Bad Request\n",
				Header: http.Header{"Content-Type": {"text/plain; charset=utf-8"}},
			},
			config: &don.Config{DefaultEncoding: "text/plain"},
			handler: don.H(func(ctx context.Context, req *struct{ Hello string }) (any, error) {
				return nil, nil
			}),
			body: `{"foo":"bar"}`,
		},
		{
			message: "should return error on unacceptable",
			expected: response{
				Code: http.StatusNotAcceptable,
				Body: "Not Acceptable\n",
				Header: http.Header{
					"Content-Type":           {"text/plain; charset=utf-8"},
					"X-Content-Type-Options": {"nosniff"},
				},
			},
			config: &don.Config{DefaultEncoding: "text/plain"},
			handler: don.H(func(ctx context.Context, req *don.Empty) ([]string, error) {
				return []string{"foo", "bar"}, nil
			}),
			header: http.Header{"Accept": {"text/plain; charset=utf-8"}},
		},
		{
			message: "should return error on unsupported accept",
			expected: response{
				Code: http.StatusNotAcceptable,
				Body: "Not Acceptable\n",
				Header: http.Header{
					"Content-Type":           {"text/plain; charset=utf-8"},
					"X-Content-Type-Options": {"nosniff"},
				},
			},
			config: &don.Config{DefaultEncoding: "text/plain"},
			handler: don.H(func(ctx context.Context, req *don.Empty) (any, error) {
				return nil, nil
			}),
			header: http.Header{"Accept": {"application/msword"}},
		},
		{
			message: "should return error on unsupported content type",
			expected: response{
				Code:   http.StatusNotAcceptable,
				Body:   "Not Acceptable\n",
				Header: http.Header{"Content-Type": {"text/plain; charset=utf-8"}},
			},
			config: &don.Config{DefaultEncoding: "text/plain"},
			handler: don.H(func(ctx context.Context, req *don.Empty) (any, error) {
				return nil, nil
			}),
			body:   `foo`,
			header: http.Header{"Content-Type": {"application/msword"}},
		},
		{
			message: "should read text request",
			expected: response{
				Code:   http.StatusOK,
				Body:   "foo\n",
				Header: http.Header{"Content-Type": {"text/plain; charset=utf-8"}},
			},
			config: &don.Config{DefaultEncoding: "text/plain"},
			handler: don.H(func(ctx context.Context, req *string) (string, error) {
				return *req, nil
			}),
			body: `foo`,
		},
	}

	for _, test := range tests {
		// Pass config to handler.
		api := don.New(test.config)
		api.Get("/", test.handler)

		r := httptest.NewRequest("", "/", strings.NewReader(test.body))
		r.Header = test.header

		w := httptest.NewRecorder()

		test.handler.ServeHTTP(w, r)
		res := response{w.Code, w.Body.String(), w.HeaderMap}
		if diff := cmp.Diff(test.expected, res); diff != "" {
			t.Errorf("%s:\n%s", test.message, diff)
		}
	}
}
