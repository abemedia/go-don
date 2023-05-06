package don_test

import (
	"strings"
	"testing"

	"github.com/abemedia/go-don"
	_ "github.com/abemedia/go-don/encoding/text"
	"github.com/abemedia/go-don/internal/test"
	"github.com/abemedia/httprouter"
	"github.com/valyala/fasthttp"
)

func TestGroup(t *testing.T) {
	api := don.New(nil)
	group := api.Group("/group")
	test.Router(t, group, api.RequestHandler(), "/group")
}

func TestGroup_Use(t *testing.T) {
	mwCalled := false

	api := don.New(nil)
	api.Get("/", func(*fasthttp.RequestCtx, httprouter.Params) {})

	group := api.Group("/group")
	group.Get("/foo", func(*fasthttp.RequestCtx, httprouter.Params) {})
	group.Use(func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			if strings.HasPrefix(string(ctx.Path()), "/group/") {
				mwCalled = true
			} else {
				t.Error("middleware called outside of group")
			}
		}
	})

	h := api.RequestHandler()

	urls := []string{"/", "/group/foo"}
	for _, url := range urls {
		ctx := &fasthttp.RequestCtx{}
		ctx.Request.Header.SetMethod(fasthttp.MethodGet)
		ctx.Request.SetRequestURI(url)

		h(ctx)

		if ctx.Response.StatusCode() >= 300 {
			t.Errorf("expected success status got %d", ctx.Response.Header.StatusCode())
		}
	}

	if !mwCalled {
		t.Error("group middleware wasn't called")
	}
}
