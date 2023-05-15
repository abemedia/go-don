package text_test

import (
	"errors"
	"testing"

	"github.com/abemedia/go-don/encoding"
	"github.com/abemedia/go-don/pkg/httptest"
	"github.com/google/go-cmp/cmp"
	"github.com/valyala/fasthttp"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		in   any
		want string
	}{
		{"test", "test"},
		{[]byte("test"), "test"},
		{int(5), "5"},
		{int8(5), "5"},
		{int16(5), "5"},
		{int32(5), "5"},
		{int64(5), "5"},
		{uint(5), "5"},
		{uint8(5), "5"},
		{uint16(5), "5"},
		{uint32(5), "5"},
		{uint64(5), "5"},
		{float32(5.1), "5.1"},
		{float64(5.1), "5.1"},
		{true, "true"},
		{errors.New("test"), "test"},
		{marshaler{}, "test"},
		{&marshaler{}, "test"},
		{stringer{}, "test"},
		{&stringer{}, "test"},
	}

	for _, test := range tests {
		enc := encoding.GetEncoder("text/plain")
		if enc == nil {
			t.Fatal("encoder not found")
		}

		ctx := httptest.NewRequest(fasthttp.MethodGet, "/", "", nil)
		if err := enc(ctx, test.in); err != nil {
			t.Error(err)
		} else {
			if diff := cmp.Diff(test.want, string(ctx.Response.Body())); diff != "" {
				t.Errorf("%T: %s", test.in, diff)
			}
		}
	}
}

type marshaler struct{}

func (m marshaler) MarshalText() ([]byte, error) {
	return []byte("test"), nil
}

type stringer struct{}

func (m stringer) String() string {
	return "test"
}
