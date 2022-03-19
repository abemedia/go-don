package decoder_test

import (
	"testing"

	"github.com/abemedia/go-don/decoder"
	"github.com/gorilla/schema"
)

func BenchmarkDecoder(b *testing.B) {
	type child struct {
		String string `header:"string"`
	}

	type test struct {
		String    string   `header:"string"`
		StringPtr *string  `header:"string"`
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
	}

	b.Run("Gorilla", func(b *testing.B) {
		dec := schema.NewDecoder()

		for i := 0; i < b.N; i++ {
			out := &test{}
			if err := dec.Decode(out, in); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("DonCached", func(b *testing.B) {
		dec, err := decoder.NewCachedDecoder(test{}, "schema")
		if err != nil {
			b.Fatal(err)
		}

		for i := 0; i < b.N; i++ {
			out := &test{}
			if err := dec.Decode(in, out); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Don", func(b *testing.B) {
		dec := decoder.NewDecoder("schema")

		for i := 0; i < b.N; i++ {
			out := &test{}
			if err := dec.Decode(in, out); err != nil {
				b.Fatal(err)
			}
		}
	})
}
