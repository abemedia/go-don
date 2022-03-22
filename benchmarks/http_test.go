package don_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/abemedia/go-don"
	"github.com/gin-gonic/gin"
	"github.com/valyala/fasthttp"
)

func BenchmarkDon_HTTP(b *testing.B) {
	api := don.New(nil)

	api.Post("/:path", don.H(func(ctx context.Context, req don.Empty) (string, error) {
		return "foo", nil
	}))

	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		b.Fatal(err)
	}

	srv := fasthttp.Server{Handler: api.RequestHandler()}
	go srv.Serve(ln)

	url := fmt.Sprintf("http://localhost:%d/path", ln.Addr().(*net.TCPAddr).Port)

	for i := 0; i < b.N; i++ {
		fasthttp.Get(nil, url)
	}

	srv.Shutdown()
}

func BenchmarkGin_HTTP(b *testing.B) {
	h := func(c *gin.Context) {
		c.String(200, "foo")
	}

	gin.SetMode("release")
	router := gin.New()
	router.POST("/:path", h)

	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		b.Fatal(err)
	}

	srv := http.Server{Handler: router}
	go srv.Serve(ln)

	url := fmt.Sprintf("http://localhost:%d/path", ln.Addr().(*net.TCPAddr).Port)

	for i := 0; i < b.N; i++ {
		fasthttp.Get(nil, url)
	}

	srv.Shutdown(context.Background())
}
