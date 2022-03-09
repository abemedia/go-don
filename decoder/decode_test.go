package decoder_test

import (
	"testing"

	"github.com/abemedia/go-don/decoder"
	"github.com/google/go-cmp/cmp"
)

func TestDecode(t *testing.T) {
	type child struct {
		String string `header:"string"`
	}

	type test struct {
		String    string   `header:"string"`
		Int       int      `header:"int"`
		Int8      int8     `header:"int8"`
		Int16     int16    `header:"int16"`
		Int32     int32    `header:"int32"`
		Int64     int64    `header:"int64"`
		Uint      uint     `header:"uint"`
		Uint8     uint8    `header:"uint8"`
		Uint16    uint16   `header:"uint16"`
		Uint32    uint32   `header:"uint32"`
		Uint64    uint64   `header:"uint64"`
		Float32   float32  `header:"float32"`
		Float64   float64  `header:"float64"`
		Bool      bool     `header:"bool"`
		Bytes     []byte   `header:"bytes"`
		Strings   []string `header:"strings"`
		Nested    child
		NestedPtr *child
	}

	in := decoder.MapGetter{
		"string":  {"string"},
		"strings": {"string", "string"},
		"int":     {"1"},
		"int8":    {"1"},
		"int16":   {"1"},
		"int32":   {"1"},
		"int64":   {"1"},
		"uint":    {"1"},
		"uint8":   {"1"},
		"uint16":  {"1"},
		"uint32":  {"1"},
		"uint64":  {"1"},
		"float32": {"1"},
		"float64": {"1"},
		"bool":    {"true"},
		"bytes":   {"bytes"},
	}

	expected := &test{
		String:  "string",
		Int:     1,
		Int8:    1,
		Int16:   1,
		Int32:   1,
		Int64:   1,
		Uint:    1,
		Uint8:   1,
		Uint16:  1,
		Uint32:  1,
		Uint64:  1,
		Float32: 1,
		Float64: 1,
		Bool:    true,
		Bytes:   []byte("bytes"),
		Strings: []string{"string", "string"},
		Nested: child{
			String: "string",
		},
		NestedPtr: &child{
			String: "string",
		},
	}

	dec, err := decoder.NewDecoder(test{}, "header")
	if err != nil {
		t.Fatal(err)
	}

	actual := &test{}
	if err = dec.Decode(in, actual); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf(diff)
	}
}
