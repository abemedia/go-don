package don_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/abemedia/go-don"
	_ "github.com/abemedia/go-don/encoding/text"
	"github.com/abemedia/go-don/pkg/httptest"
	"github.com/valyala/fasthttp"
)

func TestAPI(t *testing.T) {
	type GreetRequest struct {
		Name string `query:"name"`
	}

	api := don.New(nil)
	api.Get("/", don.H(func(ctx context.Context, req GreetRequest) (string, error) {
		if req.Name == "" {
			return "", don.ErrBadRequest
		}

		return fmt.Sprintf("Hello %s.", req.Name), nil
	}))

	h := api.RequestHandler()

	ctx := httptest.NewRequest(fasthttp.MethodGet, "/?name=mike", "", nil)

	h(ctx)

	if ctx.Response.StatusCode() >= 300 {
		t.Error(ctx.Response.Header.StatusMessage())
	}

	actual := string(ctx.Response.Body())
	expected := "Hello mike.\n"

	if expected != actual {
		t.Errorf("expected: '%s' actual: '%s'", expected, actual)
	}
}
