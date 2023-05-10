package form_test

import (
	"testing"

	"github.com/abemedia/go-don/internal/test"
)

func TestForm(t *testing.T) {
	type item struct {
		Foo string `form:"foo"`
	}

	t.Run("URLEncoded", func(t *testing.T) {
		test.Decode(t, test.EncodingOptions[item]{
			Mime:   "application/x-www-form-urlencoded",
			Raw:    "foo=bar",
			Parsed: item{Foo: "bar"},
		})
	})

	t.Run("Multipart", func(t *testing.T) {
		test.Decode(t, test.EncodingOptions[item]{
			Mime:   `multipart/form-data;boundary="boundary"`,
			Raw:    "--boundary\nContent-Disposition: form-data; name=\"foo\"\n\nbar\n--boundary\n",
			Parsed: item{Foo: "bar"},
		})
	})
}
