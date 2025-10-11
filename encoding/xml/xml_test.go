package xml_test

import (
	"testing"

	"github.com/abemedia/go-don/internal/testutil"
)

type item struct {
	Foo string `xml:"foo"`
}

var opt = testutil.EncodingOptions[item]{
	Mime:   "application/xml",
	Raw:    "<item><foo>bar</foo></item>",
	Parsed: item{Foo: "bar"},
}

func TestXML(t *testing.T) {
	testutil.TestEncoding(t, opt)
}

func BenchmarkXML(b *testing.B) {
	testutil.BenchmarkEncoding(b, opt)
}
