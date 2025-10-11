package benchmarks_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/abemedia/go-don"
	_ "github.com/abemedia/go-don/encoding/text"
	"github.com/gin-gonic/gin"
	"github.com/valyala/fasthttp"
)

func BenchmarkDon_HTTP(b *testing.B) {
	api := don.New(nil)
	api.Get("/:path", don.H(func(ctx context.Context, req any) (string, error) {
		return "foo", nil
	}))

	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		b.Fatal(err)
	}

	srv := fasthttp.Server{Handler: api.RequestHandler()}
	go srv.Serve(ln)

	url := fmt.Sprintf("http://%s/path", ln.Addr())

	for b.Loop() {
		fasthttp.Get(nil, url)
	}

	srv.Shutdown()
}

func BenchmarkGin_HTTP(b *testing.B) {
	gin.SetMode("release")

	router := gin.New()
	router.GET("/:path", func(c *gin.Context) {
		c.String(200, "foo")
	})

	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		b.Fatal(err)
	}

	srv := http.Server{Handler: router}
	go srv.Serve(ln)

	url := fmt.Sprintf("http://%s/path", ln.Addr())

	for b.Loop() {
		fasthttp.Get(nil, url)
	}

	srv.Shutdown(context.Background())
}
