package yaml_test

import (
	"testing"

	"github.com/abemedia/go-don/internal/testutil"
)

type item struct {
	Foo string `yaml:"foo"`
}

var opt = testutil.EncodingOptions[item]{
	Mime:   "application/yaml",
	Raw:    "foo: bar\n",
	Parsed: item{Foo: "bar"},
}

func TestYAML(t *testing.T) {
	testutil.TestEncoding(t, opt)
}

func BenchmarkYAML(b *testing.B) {
	testutil.BenchmarkEncoding(b, opt)
}
