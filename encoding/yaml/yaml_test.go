package yaml_test

import (
	"testing"

	"github.com/abemedia/go-don/internal/test"
)

type item struct {
	Foo string `yaml:"foo"`
}

var opt = test.EncodingOptions[item]{
	Mime:   "application/x-yaml",
	Raw:    "foo: bar\n",
	Parsed: item{Foo: "bar"},
}

func TestYAML(t *testing.T) {
	test.Encoding(t, opt)
}

func BenchmarkYAML(b *testing.B) {
	test.BenchmarkEncoding(b, opt)
}
