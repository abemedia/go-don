package don_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/abemedia/go-don"
	_ "github.com/abemedia/go-don/encoding/text"
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

	h := api.Router()

	r := httptest.NewRequest("GET", "/?name=mike", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	if w.Result().StatusCode >= 300 {
		t.Error(w.Result().Status)
	}

	buf, _ := ioutil.ReadAll(w.Result().Body)
	actual := string(buf)
	expected := "Hello mike.\n"

	if expected != actual {
		t.Errorf("expected: '%s' actual: '%s'", expected, actual)
	}
}
