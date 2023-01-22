package decoder

import (
	"github.com/abemedia/go-don/internal/byteconv"
	"github.com/valyala/fasthttp"
)

type HeaderDecoder struct {
	dec *CachedDecoder
}

func NewHeaderDecoder(v any, tag string) (*HeaderDecoder, error) {
	dec, err := NewCachedDecoder(v, tag)
	if err != nil {
		return nil, err
	}

	return &HeaderDecoder{dec}, nil
}

func (d *HeaderDecoder) Decode(data *fasthttp.RequestHeader, v any) error {
	return d.dec.Decode((*HeaderGetter)(data), v)
}

type HeaderGetter fasthttp.RequestHeader

func (ps *HeaderGetter) Get(key string) string {
	return byteconv.Btoa((*fasthttp.RequestHeader)(ps).Peek(key))
}

func (ps *HeaderGetter) Values(key string) []string {
	arg := (*fasthttp.RequestHeader)(ps).Peek(key)
	if len(arg) == 0 {
		return nil
	}

	return []string{byteconv.Btoa(arg)}
}
