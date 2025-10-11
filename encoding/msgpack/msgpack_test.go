package msgpack_test

import (
	"testing"

	"github.com/abemedia/go-don/internal/testutil"
)

type item struct {
	Foo string `msgpack:"foo"`
}

var opt = testutil.EncodingOptions[item]{
	Mime:   "application/x-msgpack",
	Raw:    "\x81\xa3foo\xa3bar",
	Parsed: item{Foo: "bar"},
}

func TestMsgpack(t *testing.T) {
	testutil.TestEncoding(t, opt)
}

func BenchmarkMsgpack(b *testing.B) {
	testutil.BenchmarkEncoding(b, opt)
}
