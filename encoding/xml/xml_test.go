package xml_test

import (
	"testing"

	"github.com/abemedia/go-don/internal/test"
)

func TestXML(t *testing.T) {
	type item struct {
		Foo string `xml:"foo"`
	}

	test.Encoding(t, test.EncodingTest[item]{
		Mime:   "application/xml",
		Raw:    "<item><foo>bar</foo></item>",
		Parsed: item{Foo: "bar"},
	})
}
