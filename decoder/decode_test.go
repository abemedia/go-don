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
		String string `header:"String"`
	}

	type test struct {
		Unmarshaler    unmarshaler  `header:"String"`
		UnmarshalerPtr *unmarshaler `header:"String"`
		String         string       `header:"String"`
		StringPtr      *string      `header:"String"`
		Int            int          `header:"Number"`
		Int8           int8         `header:"Number"`
		Int16          int16        `header:"Number"`
		Int32          int32        `header:"Number"`
		Int64          int64        `header:"Number"`
		Uint           uint         `header:"Number"`
		Uint8          uint8        `header:"Number"`
		Uint16         uint16       `header:"Number"`
		Uint32         uint32       `header:"Number"`
		Uint64         uint64       `header:"Number"`
		Float32        float32      `header:"Number"`
		Float64        float64      `header:"Number"`
		Bool           bool         `header:"Bool"`
		Bytes          []byte       `header:"String"`
		Strings        []string     `header:"Strings"`
		Nested         child
		NestedPtr      *child
	}

	in := decoder.MapGetter{
		"String":  {"string"},
		"Strings": {"string", "string"},
		"Number":  {"1"},
		"Bool":    {"true"},
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
		String string `header:"String"`
	}

	in := decoder.MapGetter{
		"String": {"string"},
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
