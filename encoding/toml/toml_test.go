package toml_test

import (
	"testing"

	"github.com/abemedia/go-don/internal/test"
)

type item struct {
	Foo string `toml:"foo"`
}

var opt = test.EncodingOptions[item]{
	Mime:   "application/toml",
	Raw:    `foo = "bar"` + "\n",
	Parsed: item{Foo: "bar"},
}

func TestTOML(t *testing.T) {
	test.Encoding(t, opt)
}

func BenchmarkTOML(b *testing.B) {
	test.BenchmarkEncoding(b, opt)
}
