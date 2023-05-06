package decoder_test

import (
	"testing"

	"github.com/abemedia/go-don/decoder"
	"github.com/abemedia/httprouter"
	"github.com/google/go-cmp/cmp"
	"github.com/valyala/fasthttp"
)

func TestAdapters(t *testing.T) {
	t.Run("Map", func(t *testing.T) {
		in := map[string][]string{"string": {"string"}}
		testAdapter(t, decoder.Map(in))
	})

	t.Run("Args", func(t *testing.T) {
		in := &fasthttp.Args{}
		in.Add("string", "string")
		testAdapter(t, (*decoder.Args)(in))
	})

	t.Run("Header", func(t *testing.T) {
		in := &fasthttp.RequestHeader{}
		in.Add("string", "string")
		testAdapter(t, (*decoder.Header)(in))
	})

	t.Run("Params", func(t *testing.T) {
		in := httprouter.Params{{Key: "string", Value: "string"}}
		testAdapter(t, decoder.Params(in))
	})
}

func testAdapter(t *testing.T, in decoder.Getter) {
	t.Helper()

	type item struct {
		Zero    string   `field:"empty"`
		Nil     []string `field:"empty"`
		String  string   `field:"string"`
		Strings []string `field:"string"`
	}

	want := &item{
		String:  "string",
		Strings: []string{"string"},
	}

	dec, err := decoder.NewCached(item{}, "field")
	if err != nil {
		t.Fatal(err)
	}

	got := &item{}
	if err = dec.Decode(in, got); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf(diff)
	}
}
