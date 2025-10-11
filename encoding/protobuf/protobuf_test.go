package protobuf_test

import (
	"errors"
	"testing"

	"github.com/abemedia/go-don"
	"github.com/abemedia/go-don/encoding"
	"github.com/abemedia/go-don/encoding/protobuf/testdata"
	"github.com/abemedia/go-don/internal/testutil"
)

//go:generate protoc testdata/test.proto --go_out=. --go_opt=paths=source_relative

var opt = testutil.EncodingOptions[*testdata.Item]{
	Mime:   "application/protobuf",
	Raw:    "\n\x03bar",
	Parsed: &testdata.Item{Foo: "bar"},
}

func TestProtobuf(t *testing.T) {
	testutil.TestEncoding(t, opt)
}

func TestProtobufError(t *testing.T) {
	ctx := testutil.NewRequest("", "", "", nil)
	v := "test"

	t.Run("Decode", func(t *testing.T) {
		dec := encoding.GetDecoder("application/protobuf")
		err := dec(ctx, &v)
		if !errors.Is(err, don.ErrUnsupportedMediaType) {
			t.Fatal("should fail")
		}
	})

	t.Run("Encode", func(t *testing.T) {
		enc := encoding.GetEncoder("application/protobuf")
		err := enc(ctx, &v)
		if !errors.Is(err, don.ErrNotAcceptable) {
			t.Fatal("should fail")
		}
	})
}

func BenchmarkProtobuf(b *testing.B) {
	testutil.BenchmarkEncoding(b, opt)
}
