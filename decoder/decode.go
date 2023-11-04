package decoder

import (
	"errors"
	"reflect"
	"sync"
)

var (
	ErrUnsupportedType = errors.New("decoder: unsupported type")
	ErrTagNotFound     = errors.New("decoder: tag not found")
)

type Getter interface {
	Get(key string) string
	Values(key string) []string
}

type Decoder struct {
	tag   string
	cache sync.Map
}

func New(tag string) *Decoder {
	return &Decoder{tag: tag}
}

func (d *Decoder) Decode(data Getter, v any) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Pointer {
		return ErrUnsupportedType
	}

	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return ErrUnsupportedType
	}

	t := val.Type()

	dec, ok := d.cache.Load(t)
	if !ok {
		var err error
		dec, err = compile(t, d.tag, t.Kind() == reflect.Ptr)
		if err != nil {
			if err != ErrTagNotFound { //nolint:errorlint,goerr113
				return err
			}
			dec = noopDecoder
		}

		d.cache.Store(t, dec)
	}

	return dec.(decoder)(val, data) //nolint:forcetypeassert
}

type CachedDecoder[V any] struct {
	dec decoder
}

func NewCached[V any](v V, tag string) (*CachedDecoder[V], error) {
	t := reflect.TypeOf(v)
	if t == nil {
		return nil, ErrUnsupportedType
	}

	t, k, ptr := typeKind(t)
	if k != reflect.Struct {
		return nil, ErrUnsupportedType
	}

	dec, err := compile(t, tag, ptr)
	if err != nil {
		return nil, err
	}

	return &CachedDecoder[V]{dec}, nil
}

func (d *CachedDecoder[V]) Decode(data Getter, v *V) error {
	return d.dec(reflect.ValueOf(v).Elem(), data)
}

func (d *CachedDecoder[V]) DecodeValue(data Getter, v reflect.Value) error {
	return d.dec(v, data)
}
