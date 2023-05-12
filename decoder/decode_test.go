package decoder_test

import (
	"errors"
	"reflect"
	"strconv"
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
		String string `field:"string"`
	}

	type test struct {
		Unmarshaler    unmarshaler  `field:"string"`
		UnmarshalerPtr *unmarshaler `field:"string"`
		String         string       `field:"string"`
		StringPtr      *string      `field:"string"`
		Int            int          `field:"number"`
		Int8           int8         `field:"number"`
		Int16          int16        `field:"number"`
		Int32          int32        `field:"number"`
		Int64          int64        `field:"number"`
		Uint           uint         `field:"number"`
		Uint8          uint8        `field:"number"`
		Uint16         uint16       `field:"number"`
		Uint32         uint32       `field:"number"`
		Uint64         uint64       `field:"number"`
		Float32        float32      `field:"number"`
		Float64        float64      `field:"number"`
		Bool           bool         `field:"bool"`
		Bytes          []byte       `field:"string"`
		Strings        []string     `field:"strings"`
		Nested         child
		NestedPtr      *child
		unexported     string `field:"string"` //nolint:unused
	}

	in := decoder.Map{
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

	exportAll := cmp.Exporter(func(t reflect.Type) bool { return true })

	t.Run("Decoder", func(t *testing.T) {
		dec := decoder.New("field")
		actual := &test{}
		if err := dec.Decode(in, actual); err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(expected, actual, exportAll); diff != "" {
			t.Errorf(diff)
		}
	})

	t.Run("CachedDecoder", func(t *testing.T) {
		dec, err := decoder.NewCached(test{}, "field")
		if err != nil {
			t.Fatal(err)
		}

		actual := &test{}
		if err := dec.Decode(in, actual); err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(expected, actual, exportAll); diff != "" {
			t.Errorf(diff)
		}

		actual = &test{}
		val := reflect.ValueOf(actual).Elem()
		if err = dec.DecodeValue(in, val); err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(expected, actual, exportAll); diff != "" {
			t.Errorf(diff)
		}
	})

	t.Run("CachedDecoder_NilPointer", func(t *testing.T) {
		dec, err := decoder.NewCached(&test{}, "field")
		if err != nil {
			t.Fatal(err)
		}

		var actual *test
		if err := dec.Decode(in, &actual); err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(expected, actual, exportAll); diff != "" {
			t.Errorf(diff)
		}
	})
}

func TestDecodeError(t *testing.T) {
	data := decoder.Map{"test": {"test"}}

	tests := []struct {
		target any
		error  error
	}{
		{
			target: "",
			error:  decoder.ErrUnsupportedType,
		},
		{
			target: new(string),
			error:  decoder.ErrUnsupportedType,
		},
		{
			target: new(int),
			error:  decoder.ErrUnsupportedType,
		},
		{
			target: &struct {
				Test string `json:"test"`
			}{},
			error: decoder.ErrTagNotFound,
		},
		{
			target: &struct {
				Test chan string `field:"test"`
			}{},
			error: decoder.ErrUnsupportedType,
		},
		{
			target: &struct {
				Test  string `field:"test"`
				Child struct {
					Test chan string `field:"test"`
				}
			}{},
			error: decoder.ErrUnsupportedType,
		},
		{
			target: &struct {
				Test int `field:"test"`
			}{},
			error: strconv.ErrSyntax,
		},
		{
			target: &struct {
				Test uint `field:"test"`
			}{},
			error: strconv.ErrSyntax,
		},
		{
			target: &struct {
				Test float64 `field:"test"`
			}{},
			error: strconv.ErrSyntax,
		},
		{
			target: &struct {
				Test bool `field:"test"`
			}{},
			error: strconv.ErrSyntax,
		},
	}

	t.Run("Decoder", func(t *testing.T) {
		for _, test := range tests {
			dec := decoder.New("field")
			err := dec.Decode(data, test.target)
			if errors.Is(test.error, decoder.ErrTagNotFound) {
				if err != nil {
					t.Errorf("should silently ignore error %q for %T", test.error, test.target)
				}
			} else {
				if !errors.Is(err, test.error) {
					t.Errorf("should return %q for %T: %q", test.error, test.target, err)
				}
			}
		}
	})

	t.Run("CachedDecoder", func(t *testing.T) {
		for _, test := range tests {
			dec, err := decoder.NewCached(test.target, "field")
			if err != nil {
				if !errors.Is(err, test.error) {
					t.Errorf("should return %q for %T: %q", test.error, test.target, err)
				}
				continue
			}
			err = dec.Decode(data, &test.target)
			if !errors.Is(err, test.error) {
				t.Errorf("should return %q for %T: %q", test.error, test.target, err)
			}
		}
	})
}
