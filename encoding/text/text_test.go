package text_test

import (
	"testing"

	"github.com/abemedia/go-don/internal/testutil"
)

var opt = testutil.EncodingOptions[string]{
	Mime:   "text/plain",
	Raw:    "foo",
	Parsed: "foo",
}

func TestText(t *testing.T) {
	testutil.TestEncoding(t, opt)
}

func BenchmarkText(b *testing.B) {
	testutil.BenchmarkEncoding(b, opt)
}
