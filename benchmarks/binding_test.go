package benchmarks_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/abemedia/go-don"
	_ "github.com/abemedia/go-don/encoding/text"
	"github.com/gin-gonic/gin"
	"github.com/valyala/fasthttp"
)

func BenchmarkDon_BindRequest(b *testing.B) {
	type request struct {
		Path   string `path:"path"`
		Header string `header:"Header"`
		Query  string `query:"query"`
	}

	api := don.New(nil)
	api.Post("/:path", don.H(func(ctx context.Context, req request) (any, error) {
		return nil, nil
	}))

	h := api.RequestHandler()

	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("POST")
	ctx.Request.SetRequestURI("/path?query=query")
	ctx.Request.Header.Set("Header", "header")

	for b.Loop() {
		h(ctx)
	}
}

func BenchmarkGin_BindRequest(b *testing.B) {
	type request struct {
		Path   string `uri:"path"`
		Header string `header:"Header"`
		Query  string `form:"query"`
	}

	gin.SetMode("release")
	router := gin.New()
	router.POST("/:path", func(c *gin.Context) {
		req := &request{}
		c.ShouldBindHeader(req)
		c.ShouldBindQuery(req)
		c.ShouldBindUri(req)
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/path?query=query", nil)
	r.Header.Add("Header", "header")

	for b.Loop() {
		router.ServeHTTP(w, r)
	}
}
