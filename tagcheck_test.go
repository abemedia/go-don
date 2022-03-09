package don

import (
	"testing"
)

func TestCheck(t *testing.T) {
	tests := []struct {
		in       interface{}
		expected bool
	}{
		{
			&struct {
				foo string `foo:"1"`
			}{},
			false,
		},
		{
			&struct {
				Foo string `foo:"1"`
			}{},
			true,
		},
		{
			&struct {
				Foo struct {
					Foo string `foo:"1"`
				}
			}{},
			true,
		},
		{
			&struct {
				Foo *struct {
					Foo string `foo:"1"`
				}
			}{},
			true,
		},
	}

	for _, test := range tests {
		actual := hasTag(test.in, "foo")
		if actual != test.expected {
			t.Errorf("expected %t for %T", test.expected, test.in)
		}
	}
}
