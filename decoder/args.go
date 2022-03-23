package decoder

import (
	"github.com/abemedia/go-don/internal"
	"github.com/valyala/fasthttp"
)

type ArgsDecoder struct {
	dec *CachedDecoder
}

func NewArgsDecoder(v interface{}, tag string) (*ArgsDecoder, error) {
	dec, err := NewCachedDecoder(v, tag)
	if err != nil {
		return nil, err
	}

	return &ArgsDecoder{dec}, nil
}

func (d *ArgsDecoder) Decode(data *fasthttp.Args, v interface{}) error {
	return d.dec.Decode((*ArgsGetter)(data), v)
}

type ArgsGetter fasthttp.Args

func (ps *ArgsGetter) Get(key string) string {
	return internal.Btoa((*fasthttp.Args)(ps).Peek(key))
}

func (ps *ArgsGetter) Values(key string) []string {
	args := (*fasthttp.Args)(ps).PeekMulti(key)
	if len(args) == 0 {
		return nil
	}

	res := make([]string, len(args))
	for i := range args {
		res[i] = internal.Btoa(args[i])
	}

	return res
}
