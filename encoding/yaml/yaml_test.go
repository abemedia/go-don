package yaml_test

import (
	"testing"

	"github.com/abemedia/go-don/internal/test"
)

func TestYAML(t *testing.T) {
	type item struct {
		Foo string `yaml:"foo"`
	}

	test.Encoding(t, test.EncodingOptions[item]{
		Mime:   "application/x-yaml",
		Raw:    "foo: bar\n",
		Parsed: item{Foo: "bar"},
	})
}
