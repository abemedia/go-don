package text_test

import (
	"errors"
	"testing"

	"github.com/abemedia/go-don/encoding/text"
	"github.com/abemedia/go-don/pkg/httptest"
	"github.com/google/go-cmp/cmp"
	"github.com/valyala/fasthttp"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		in   any
		want string
	}{
		{"test", "test\n"},
		{[]byte("test"), "test\n"},
		{int(5), "5\n"},
		{int8(5), "5\n"},
		{int16(5), "5\n"},
		{int32(5), "5\n"},
		{int64(5), "5\n"},
		{uint(5), "5\n"},
		{uint8(5), "5\n"},
		{uint16(5), "5\n"},
		{uint32(5), "5\n"},
		{uint64(5), "5\n"},
		{float32(5.1), "5.1\n"},
		{float64(5.1), "5.1\n"},
		{true, "true\n"},
		{errors.New("test"), "test\n"},
		{marshaler{}, "test\n"},
		{&marshaler{}, "test\n"},
		{stringer{}, "test\n"},
		{&stringer{}, "test\n"},
	}

	for _, test := range tests {
		ctx := httptest.NewRequest(fasthttp.MethodGet, "/", "", nil)
		if err := text.Encode(ctx, test.in); err != nil {
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
