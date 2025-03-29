package don_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/abemedia/go-don"
	_ "github.com/abemedia/go-don/encoding/text"
	"github.com/abemedia/go-don/internal/test"
	"github.com/abemedia/httprouter"
	"github.com/valyala/fasthttp"
)

func TestAPI(t *testing.T) {
	t.Run("Router", func(t *testing.T) {
		api := don.New(nil)
		test.Router(t, api, api.RequestHandler(), "")
	})

	t.Run("ListenAndServe", func(t *testing.T) {
		testAPI(t, func(api *don.API, ln net.Listener) error {
			ln.Close() //nolint:errcheck
			return api.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", ln.Addr().(*net.TCPAddr).Port))
		})
	})

	t.Run("Serve", func(t *testing.T) {
		testAPI(t, func(api *don.API, ln net.Listener) error { return api.Serve(ln) })
	})
}

func testAPI(t *testing.T, serve func(api *don.API, ln net.Listener) error) {
	t.Helper()

	api := don.New(nil)
	api.Get("/", func(ctx *fasthttp.RequestCtx, _ httprouter.Params) {
		ctx.WriteString("OK") //nolint:errcheck
	})

	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		if err := serve(api, ln); err != nil {
			t.Error(err)
		}
	}()

	time.Sleep(time.Millisecond) // Wait for server to be listening.

	if t.Failed() {
		t.FailNow() // Fail fast if calling serve failed.
	}

	code, body, err := fasthttp.Get(nil, fmt.Sprintf("http://localhost:%d/", ln.Addr().(*net.TCPAddr).Port))
	if err != nil {
		t.Fatal(err)
	}

	if code != fasthttp.StatusOK {
		t.Fatal("should return success status")
	}

	if string(body) != "OK" {
		t.Fatalf(`should return body "OK": %q`, body)
	}
}
