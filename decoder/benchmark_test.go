package decoder_test

import (
	"testing"

	"github.com/abemedia/go-don/decoder"
	"github.com/gorilla/schema"
)

func BenchmarkDecoder(b *testing.B) {
	type test struct {
		String  string  `schema:"string"`
		Int     int     `schema:"int"`
		Int32   int32   `schema:"int32"`
		Int64   int64   `schema:"int64"`
		Float32 float32 `schema:"float32"`
		Float64 float64 `schema:"float64"`
		Bool    bool    `schema:"bool"`
	}

	in := decoder.MapGetter{
		"int":     {"1"},
		"int32":   {"1"},
		"int64":   {"1"},
		"float32": {"1"},
		"float64": {"1"},
		"bool":    {"true"},
	}

	b.Run("Gorilla", func(b *testing.B) {
		dec := schema.NewDecoder()

		for i := 0; i < b.N; i++ {
			actual := &test{}
			if err := dec.Decode(actual, in); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Don", func(b *testing.B) {
		dec, err := decoder.NewDecoder(test{}, "schema")
		if err != nil {
			b.Fatal(err)
		}

		for i := 0; i < b.N; i++ {
			actual := &test{}
			if err = dec.Decode(in, actual); err != nil {
				b.Fatal(err)
			}
		}
	})
}
