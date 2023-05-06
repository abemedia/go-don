package don_test

import (
	"testing"

	"github.com/abemedia/go-don"
	_ "github.com/abemedia/go-don/encoding/text"
	"github.com/abemedia/go-don/internal/test"
)

func TestAPI(t *testing.T) {
	api := don.New(nil)
	test.Router(t, api, api.RequestHandler(), "")
}
