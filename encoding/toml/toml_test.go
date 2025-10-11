package toml_test

import (
	"testing"

	"github.com/abemedia/go-don/internal/testutil"
)

type item struct {
	Foo string `toml:"foo"`
}

var opt = testutil.EncodingOptions[item]{
	Mime:   "application/toml",
	Raw:    `foo = "bar"` + "\n",
	Parsed: item{Foo: "bar"},
}

func TestTOML(t *testing.T) {
	testutil.TestEncoding(t, opt)
}

func BenchmarkTOML(b *testing.B) {
	testutil.BenchmarkEncoding(b, opt)
}
