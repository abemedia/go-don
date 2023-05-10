package xml_test

import (
	"testing"

	"github.com/abemedia/go-don/internal/test"
)

type item struct {
	Foo string `xml:"foo"`
}

var opt = test.EncodingOptions[item]{
	Mime:   "application/xml",
	Raw:    "<item><foo>bar</foo></item>",
	Parsed: item{Foo: "bar"},
}

func TestXML(t *testing.T) {
	test.Encoding(t, opt)
}

func BenchmarkXML(b *testing.B) {
	test.BenchmarkEncoding(b, opt)
}
