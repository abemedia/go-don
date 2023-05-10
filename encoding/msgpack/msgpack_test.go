package msgpack_test

import (
	"testing"

	"github.com/abemedia/go-don/internal/test"
)

type item struct {
	Foo string `msgpack:"foo"`
}

var opt = test.EncodingOptions[item]{
	Mime:   "application/x-msgpack",
	Raw:    "\x81\xa3foo\xa3bar",
	Parsed: item{Foo: "bar"},
}

func TestMsgpack(t *testing.T) {
	test.Encoding(t, opt)
}

func BenchmarkMsgpack(b *testing.B) {
	test.BenchmarkEncoding(b, opt)
}
