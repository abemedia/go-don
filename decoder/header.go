package decoder

import (
	"github.com/abemedia/go-don/internal"
	"github.com/valyala/fasthttp"
)

type HeaderDecoder struct {
	dec *CachedDecoder
}

func NewHeaderDecoder(v interface{}, tag string) (*HeaderDecoder, error) {
	dec, err := NewCachedDecoder(v, tag)
	if err != nil {
		return nil, err
	}

	return &HeaderDecoder{dec}, nil
}

func (d *HeaderDecoder) Decode(data *fasthttp.RequestHeader, v interface{}) error {
	return d.dec.Decode((*HeaderGetter)(data), v)
}

type HeaderGetter fasthttp.RequestHeader

func (ps *HeaderGetter) Get(key string) string {
	return internal.Btoa((*fasthttp.RequestHeader)(ps).Peek(key))
}

func (ps *HeaderGetter) Values(key string) []string {
	arg := (*fasthttp.RequestHeader)(ps).Peek(key)
	if len(arg) == 0 {
		return nil
	}

	return []string{internal.Btoa(arg)}
}
