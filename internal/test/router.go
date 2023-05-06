package test

import (
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
	}

	for _, test := range tests {
		path := "/" + strings.ToLower(test.desc)

		var ok bool
		test.fn(path, func(ctx *fasthttp.RequestCtx, p httprouter.Params) {
			ok = true
		})

		ctx := httptest.NewRequest(test.method, basePath+path, "", nil)
		handler(ctx)

		if code := ctx.Response.StatusCode(); code != fasthttp.StatusOK {
			t.Errorf("%s request should return success status: %s", test.desc, fasthttp.StatusMessage(code))
		}

		if !ok {
			t.Errorf("%s request should reach handler", test.desc)
		}
	}
}
