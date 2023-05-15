package text_test

import (
	"reflect"
	"testing"

	"github.com/abemedia/go-don/encoding"
	"github.com/abemedia/go-don/pkg/httptest"
	"github.com/google/go-cmp/cmp"
	"github.com/valyala/fasthttp"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		in   string
		want any
	}{
		{"test\n", "test"},
		{"test\n", []byte("test")},
		{"5\n", int(5)},
		{"5\n", int8(5)},
		{"5\n", int16(5)},
		{"5\n", int32(5)},
		{"5\n", int64(5)},
		{"5\n", uint(5)},
		{"5\n", uint8(5)},
		{"5\n", uint16(5)},
		{"5\n", uint32(5)},
		{"5\n", uint64(5)},
		{"5.1\n", float32(5.1)},
		{"5.1\n", float64(5.1)},
		{"true\n", true},
		{"test\n", unmarshaler{"test"}},
		{"test\n", &unmarshaler{"test"}},
	}

	for _, test := range tests {
		dec := encoding.GetDecoder("text/plain")
		if dec == nil {
			t.Fatal("decoder not found")
		}

		ctx := httptest.NewRequest(fasthttp.MethodGet, "/", test.in, nil)
		v := reflect.New(reflect.TypeOf(test.want)).Interface()
		if err := dec(ctx, v); err != nil {
			t.Error(err)
		} else {
			if diff := cmp.Diff(test.want, reflect.ValueOf(v).Elem().Interface()); diff != "" {
				t.Errorf("%T: %s", test.in, diff)
			}
		}
	}
}

type unmarshaler struct {
	S string
}

func (m *unmarshaler) UnmarshalText(text []byte) error {
	m.S = string(text)
	return nil
}
