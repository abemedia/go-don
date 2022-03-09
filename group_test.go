package don_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/abemedia/go-don"
	_ "github.com/abemedia/go-don/encoding/text"
)

func TestGroup(t *testing.T) {
	mwCalled := false

	api := don.New(nil)
	api.Get("/", don.E(nil))

	group := api.Group("/group")
	group.Get("/foo", don.E(nil))
	group.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/group/") {
				mwCalled = true
			} else {
				t.Error("middleware called outside of group")
			}
		})
	})

	h := api.Router()

	urls := []string{"/", "/group/foo"}

	for _, url := range urls {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", url, nil)

		h.ServeHTTP(w, r)

		if w.Result().StatusCode >= 300 {
			t.Error(w.Result().Status)
		}
	}

	if !mwCalled {
		t.Error("group middleware wasn't called")
	}
}
