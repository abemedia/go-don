package text_test

import (
	"testing"

	"github.com/abemedia/go-don/internal/test"
)

var opt = test.EncodingOptions[string]{
	Mime:   "text/plain",
	Raw:    "foo",
	Parsed: "foo",
}

func TestText(t *testing.T) {
	test.Encoding(t, opt)
}

func BenchmarkText(b *testing.B) {
	test.BenchmarkEncoding(b, opt)
}
