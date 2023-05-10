package msgpack_test

import (
	"testing"

	"github.com/abemedia/go-don/internal/test"
)

func TestMsgpack(t *testing.T) {
	type item struct {
		Foo string `msgpack:"foo"`
	}

	test.Encoding(t, test.EncodingOptions[item]{
		Mime:   "application/x-msgpack",
		Raw:    "\x81\xa3foo\xa3bar",
		Parsed: item{Foo: "bar"},
	})
}
