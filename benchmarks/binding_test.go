package benchmarks_test

import (
	"context"
	stdhttptest "net/http/httptest"
	"testing"

	"github.com/abemedia/go-don"
	_ "github.com/abemedia/go-don/encoding/text"
	"github.com/abemedia/go-don/pkg/httptest"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
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
	ctx := httptest.NewRequest("POST", "/path?query=query", "", map[string]string{"header": "header"})

	for i := 0; i < b.N; i++ {
		h(ctx)
	}
}

func BenchmarkFiber_BindRequest(b *testing.B) {
	type request struct {
		Path   string `params:"path"`
		Header string `reqHeader:"Header"`
		Query  string `query:"query"`
	}

	app := fiber.New()
	app.Post("/:path", func(c *fiber.Ctx) error {
		req := &request{}
		_ = c.ParamsParser(req)
		_ = c.ReqHeaderParser(req)
		_ = c.QueryParser(req)
		return nil
	})

	h := app.Handler()
	ctx := httptest.NewRequest("POST", "/path?query=query", "", map[string]string{"header": "header"})

	for i := 0; i < b.N; i++ {
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
		_ = c.ShouldBindHeader(req)
		_ = c.ShouldBindQuery(req)
		_ = c.ShouldBindUri(req)
	})

	w := stdhttptest.NewRecorder()
	r := stdhttptest.NewRequest("POST", "/path?query=query", nil)
	r.Header.Add("header", "header")

	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, r)
	}
}
