package don_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/abemedia/go-don"
	_ "github.com/abemedia/go-don/encoding/json"
	_ "github.com/abemedia/go-don/encoding/text"
	"github.com/abemedia/go-don/internal/httptest"
	"github.com/abemedia/httprouter"
	"github.com/google/go-cmp/cmp"
	"github.com/valyala/fasthttp"
)

func TestHandler(t *testing.T) {
	type response struct {
		Code   int
		Body   string
		Header map[string]string
	}

	tests := []struct {
		message  string
		expected response
		config   *don.Config
		handler  httprouter.Handle
		body     string
		header   map[string]string
	}{
		{
			message: "should return no content",
			expected: response{
				Code: http.StatusNoContent,
				Body: "",
				Header: map[string]string{
					"Content-Length": "0",
					"Content-Type":   "text/plain; charset=utf-8",
				},
			},
			config: &don.Config{DefaultEncoding: "text/plain"},
			handler: don.H(func(ctx context.Context, req don.Empty) (any, error) {
				return nil, nil
			}),
		},
		{
			message: "should return null",
			expected: response{
				Code: http.StatusOK,
				Body: "null\n",
				Header: map[string]string{
					"Content-Length": "0",
					"Content-Type":   "application/json; charset=utf-8",
				},
			},
			config: &don.Config{DefaultEncoding: "application/json", DisableNoContent: true},
			handler: don.H(func(ctx context.Context, req don.Empty) (any, error) {
				return nil, nil
			}),
		},
		{
			message: "should JSON encode map",
			expected: response{
				Code:   http.StatusOK,
				Body:   `{"foo":"bar"}` + "\n",
				Header: map[string]string{"Content-Type": "application/json; charset=utf-8"},
			},
			config: &don.Config{DefaultEncoding: "application/json"},
			handler: don.H(func(ctx context.Context, req don.Empty) (map[string]string, error) {
				return map[string]string{"foo": "bar"}, nil
			}),
		},
		{
			message: "should return error on unprocessable request",
			expected: response{
				Code:   http.StatusUnsupportedMediaType,
				Body:   "Unsupported Media Type\n",
				Header: map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			},
			config: &don.Config{DefaultEncoding: "text/plain"},
			handler: don.H(func(ctx context.Context, req struct{ Hello string }) (any, error) {
				return nil, nil
			}),
			body: `{"foo":"bar"}`,
		},
		{
			message: "should return error on unacceptable",
			expected: response{
				Code:   http.StatusNotAcceptable,
				Body:   "Not Acceptable",
				Header: map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			},
			config: &don.Config{DefaultEncoding: "text/plain"},
			handler: don.H(func(ctx context.Context, req don.Empty) ([]string, error) {
				return []string{"foo", "bar"}, nil
			}),
			header: map[string]string{"Accept": "text/plain; charset=utf-8"},
		},
		{
			message: "should return error on unsupported accept",
			expected: response{
				Code:   http.StatusNotAcceptable,
				Body:   "Not Acceptable",
				Header: map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			},
			config: &don.Config{DefaultEncoding: "text/plain"},
			handler: don.H(func(ctx context.Context, req don.Empty) (any, error) {
				return nil, nil
			}),
			header: map[string]string{"Accept": "application/msword"},
		},
		{
			message: "should return error on unsupported content type",
			expected: response{
				Code:   http.StatusUnsupportedMediaType,
				Body:   "Unsupported Media Type\n",
				Header: map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			},
			config: &don.Config{DefaultEncoding: "text/plain"},
			handler: don.H(func(ctx context.Context, req don.Empty) (any, error) {
				return nil, nil
			}),
			body:   `foo`,
			header: map[string]string{"Content-Type": "application/msword"},
		},
		{
			message: "should read text request",
			expected: response{
				Code:   http.StatusOK,
				Body:   "foo\n",
				Header: map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			},
			config: &don.Config{DefaultEncoding: "text/plain"},
			handler: don.H(func(ctx context.Context, req string) (string, error) {
				return req, nil
			}),
			body: `foo`,
		},
	}

	for _, test := range tests {
		// Pass config to handler.
		api := don.New(test.config)
		api.Get("/", test.handler)

		ctx := httptest.NewRequest(fasthttp.MethodGet, "/", test.body, test.header)
		api.RequestHandler()(ctx)

		res := response{ctx.Response.StatusCode(), string(ctx.Response.Body()), map[string]string{}}
		ctx.Response.Header.VisitAll(func(key, value []byte) { res.Header[string(key)] = string(value) })

		if diff := cmp.Diff(test.expected, res); diff != "" {
			t.Errorf("%s:\n%s", test.message, diff)
		}
	}
}

func TestHandlerRequest(t *testing.T) {
	type request struct {
		Path   string `path:"path"`
		Header string `header:"header"`
		Query  string `query:"query"`
		JSON   string `json:"json"`
	}

	api := don.New(&don.Config{
		DefaultEncoding: "application/json",
	})

	api.Post("/:path", don.H(func(ctx context.Context, req request) (any, error) {
		if req.Path != "path" {
			t.Error("no path")
		}
		if req.Header != "header" {
			t.Error("no header")
		}
		if req.Query != "query" {
			t.Error("no query")
		}
		if req.JSON != "json" {
			t.Error("no JSON", req)
		}
		return nil, nil
	}))

	h := api.RequestHandler()

	ctx := httptest.NewRequest(
		fasthttp.MethodPost,
		"/path?query=query",
		`{"json":"json"}`,
		map[string]string{"header": "header"},
	)

	h(ctx)
}
