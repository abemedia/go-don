package decoder_test

import (
	"testing"

	"github.com/abemedia/go-don/decoder"
	"github.com/google/go-cmp/cmp"
)

type unmarshaler string

func (h *unmarshaler) UnmarshalText(b []byte) error {
	*h = unmarshaler(":" + string(b) + ":")
	return nil
}

func TestDecode(t *testing.T) {
	type child struct {
		String string `header:"string"`
	}

	type test struct {
		Unmarshaler    unmarshaler  `header:"string"`
		UnmarshalerPtr *unmarshaler `header:"string"`
		String         string       `header:"string"`
		StringPtr      *string      `header:"string"`
		Int            int          `header:"number"`
		Int8           int8         `header:"number"`
		Int16          int16        `header:"number"`
		Int32          int32        `header:"number"`
		Int64          int64        `header:"number"`
		Uint           uint         `header:"number"`
		Uint8          uint8        `header:"number"`
		Uint16         uint16       `header:"number"`
		Uint32         uint32       `header:"number"`
		Uint64         uint64       `header:"number"`
		Float32        float32      `header:"number"`
		Float64        float64      `header:"number"`
		Bool           bool         `header:"bool"`
		Bytes          []byte       `header:"string"`
		Strings        []string     `header:"strings"`
		Nested         child
		NestedPtr      *child
	}

	in := decoder.MapGetter{
		"string":  {"string"},
		"strings": {"string", "string"},
		"number":  {"1"},
		"bool":    {"true"},
	}

	s := "string"
	u := unmarshaler(":string:")
	expected := &test{
		Unmarshaler:    ":string:",
		UnmarshalerPtr: &u,
		String:         "string",
		StringPtr:      &s,
		Int:            1,
		Int8:           1,
		Int16:          1,
		Int32:          1,
		Int64:          1,
		Uint:           1,
		Uint8:          1,
		Uint16:         1,
		Uint32:         1,
		Uint64:         1,
		Float32:        1,
		Float64:        1,
		Bool:           true,
		Bytes:          []byte("string"),
		Strings:        []string{"string", "string"},
		Nested: child{
			String: "string",
		},
		NestedPtr: &child{
			String: "string",
		},
	}

	dec, err := decoder.NewCachedDecoder(test{}, "header")
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

func TestDecodeNil(t *testing.T) {
	type test struct {
		String string `header:"string"`
	}

	in := decoder.MapGetter{
		"string": {"string"},
	}

	expected := &test{
		String: "string",
	}

	var actual *test

	dec, err := decoder.NewCachedDecoder(actual, "header")
	if err != nil {
		t.Fatal(err)
	}

	if err = dec.Decode(in, &actual); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf(diff)
	}
}
