package test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/abemedia/go-don"
	"github.com/abemedia/go-don/pkg/httptest"
	"github.com/abemedia/httprouter"
	"github.com/valyala/fasthttp"
)

func Router(t *testing.T, r don.Router, handler fasthttp.RequestHandler, basePath string) {
	t.Helper()

	tests := []struct {
		desc   string
		method string
		fn     func(path string, handle httprouter.Handle)
	}{
		{"Get", fasthttp.MethodGet, r.Get},
		{"Post", fasthttp.MethodPost, r.Post},
		{"Put", fasthttp.MethodPut, r.Put},
		{"Patch", fasthttp.MethodPatch, r.Patch},
		{"Delete", fasthttp.MethodDelete, r.Delete},
		{"Handle", fasthttp.MethodGet, func(path string, handle httprouter.Handle) {
			r.Handle(fasthttp.MethodGet, path, handle)
		}},
		{"Handler", fasthttp.MethodGet, func(path string, handle httprouter.Handle) {
			r.Handler(fasthttp.MethodGet, path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("Handler"))
			}))
		}},
		{"HandleFunc", fasthttp.MethodGet, func(path string, handle httprouter.Handle) {
			r.HandleFunc(fasthttp.MethodGet, path, func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("HandleFunc"))
			})
		}},
	}

	for _, test := range tests {
		path := "/" + strings.ToLower(test.desc)

		test.fn(path, func(ctx *fasthttp.RequestCtx, p httprouter.Params) {
			_, _ = ctx.WriteString(test.desc)
		})

		ctx := httptest.NewRequest(test.method, basePath+path, "", nil)
		handler(ctx)

		if code := ctx.Response.StatusCode(); code != fasthttp.StatusOK {
			t.Errorf("%s request should return success status: %s", test.desc, fasthttp.StatusMessage(code))
		}

		if string(ctx.Response.Body()) != test.desc {
			t.Errorf("%s request should reach handler", test.desc)
		}
	}
}
