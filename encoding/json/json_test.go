package json_test

import (
	"testing"

	"github.com/abemedia/go-don/internal/testutil"
)

type item struct {
	Foo string `json:"foo"`
}

var opt = testutil.EncodingOptions[item]{
	Mime:   "application/json",
	Raw:    `{"foo":"bar"}`,
	Parsed: item{Foo: "bar"},
}

func TestJSON(t *testing.T) {
	testutil.TestEncoding(t, opt)
}

func BenchmarkJSON(b *testing.B) {
	testutil.BenchmarkEncoding(b, opt)
}
