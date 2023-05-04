package json_test

import (
	"testing"

	"github.com/abemedia/go-don/internal/test"
)

func TestJSON(t *testing.T) {
	type item struct {
		Foo string `json:"foo"`
	}

	test.Encoding(t, test.EncodingTest[item]{
		Mime:   "application/json",
		Raw:    `{"foo":"bar"}` + "\n",
		Parsed: item{Foo: "bar"},
	})
}
