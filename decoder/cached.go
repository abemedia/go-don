package decoder

import (
	"reflect"
)

type Getter interface {
	Get(string) string
	Values(string) []string
}

type CachedDecoder struct {
	dec decoder
}

func NewCachedDecoder(v any, tag string) (*CachedDecoder, error) {
	t, k, ptr := typeKind(reflect.TypeOf(v))
	if k != reflect.Struct {
		return nil, ErrUnsupportedType
	}

	dec, err := compile(t, tag, ptr)
	if err != nil {
		return nil, err
	}

	return &CachedDecoder{dec}, nil
}

func (d *CachedDecoder) Decode(data Getter, v any) error {
	return d.dec(reflect.ValueOf(v).Elem(), data)
}
