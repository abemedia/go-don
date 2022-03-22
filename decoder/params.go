package decoder

import "github.com/abemedia/httprouter"

type ParamsDecoder struct {
	dec *CachedDecoder
}

func NewParamsDecoder(v interface{}, tag string) (*ParamsDecoder, error) {
	dec, err := NewCachedDecoder(v, tag)
	if err != nil {
		return nil, err
	}

	return &ParamsDecoder{dec}, nil
}

func (d *ParamsDecoder) Decode(data []httprouter.Param, v interface{}) error {
	return d.dec.Decode(ParamsGetter(data), v)
}

type ParamsGetter []httprouter.Param

func (ps ParamsGetter) Get(key string) string {
	for i := range ps {
		if ps[i].Key == key {
			return ps[i].Value
		}
	}

	return ""
}

func (ps ParamsGetter) Values(key string) []string {
	for i := range ps {
		if ps[i].Key == key {
			return []string{ps[i].Value}
		}
	}

	return nil
}
