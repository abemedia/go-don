package don

import (
	"testing"
)

func TestNilCheck(t *testing.T) {
	tests := []struct {
		typ      any
		data     any
		message  string
		expected bool
	}{
		{
			message:  "nil interface",
			typ:      nil,
			data:     nil,
			expected: true,
		},
		{
			message:  "empty string",
			typ:      "",
			data:     "",
			expected: true,
		},
		{
			message:  "nil struct",
			typ:      (*Empty)(nil),
			data:     (*Empty)(nil),
			expected: true,
		},
		{
			message:  "zero struct",
			typ:      struct{}{},
			data:     struct{}{},
			expected: false,
		},
		{
			message:  "don.Empty",
			typ:      Empty{},
			data:     Empty{},
			expected: true,
		},
		{
			message:  "nil map",
			typ:      (map[string]string)(nil),
			data:     (map[string]string)(nil),
			expected: true,
		},
		{
			message:  "zero map",
			typ:      (map[string]string)(nil),
			data:     map[string]string{},
			expected: false,
		},
		{
			message:  "non-zero map",
			typ:      (map[string]string)(nil),
			data:     map[string]string{"foo": "bar"},
			expected: false,
		},
		{
			message:  "nil slice",
			typ:      ([]string)(nil),
			data:     ([]string)(nil),
			expected: true,
		},
		{
			message:  "zero slice",
			typ:      ([]string)(nil),
			data:     []string{},
			expected: false,
		},
		{
			message:  "non-zero slice",
			typ:      ([]string)(nil),
			data:     []string{"aa"},
			expected: false,
		},
		{
			message:  "boolean",
			typ:      false,
			data:     false,
			expected: false,
		},
		{
			message:  "integer",
			typ:      0,
			data:     0,
			expected: false,
		},
		{
			message: "non-zero slice pointer",
			typ:     (*[]string)(nil),
			data: func() any {
				m := []string{"aa"}
				return &m
			}(),
			expected: false,
		},
	}

	for _, test := range tests {
		isNil := makeNilCheck(test.typ)
		if isNil(test.data) != test.expected {
			t.Errorf("%s should be %t", test.message, test.expected)
		}
	}
}
