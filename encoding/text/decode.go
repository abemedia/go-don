package text

import (
	"bytes"
	"encoding"
	"reflect"
	"strconv"
	"sync"

	"github.com/abemedia/go-don"
	"github.com/abemedia/go-don/internal/byteconv"
	"github.com/valyala/fasthttp"
)

//nolint:cyclop
func decode(ctx *fasthttp.RequestCtx, v any) error {
	b := bytes.TrimSpace(ctx.Request.Body())
	if len(b) == 0 {
		return nil
	}

	var err error

	switch v := v.(type) {
	case *string:
		*v = byteconv.Btoa(b)
	case *[]byte:
		*v = b
	case *int:
		*v, err = strconv.Atoi(byteconv.Btoa(b))
	case *int8:
		return decodeInt(b, v, 8)
	case *int16:
		return decodeInt(b, v, 16)
	case *int32:
		return decodeInt(b, v, 32)
	case *int64:
		return decodeInt(b, v, 64)
	case *uint:
		return decodeUint(b, v, 0)
	case *uint8:
		return decodeUint(b, v, 8)
	case *uint16:
		return decodeUint(b, v, 16)
	case *uint32:
		return decodeUint(b, v, 32)
	case *uint64:
		return decodeUint(b, v, 64)
	case *float32:
		return decodeFloat(b, v, 32)
	case *float64:
		return decodeFloat(b, v, 64)
	case *bool:
		*v, err = strconv.ParseBool(byteconv.Btoa(b))
	default:
		return unmarshal(b, v)
	}

	return err
}

func decodeInt[T int | int8 | int16 | int32 | int64](b []byte, v *T, bits int) error {
	d, err := strconv.ParseInt(byteconv.Btoa(b), 10, bits)
	*v = T(d)
	return err
}

func decodeUint[T uint | uint8 | uint16 | uint32 | uint64](b []byte, v *T, bits int) error {
	d, err := strconv.ParseUint(byteconv.Btoa(b), 10, bits)
	*v = T(d)
	return err
}

func decodeFloat[T float32 | float64](b []byte, v *T, bits int) error {
	d, err := strconv.ParseFloat(byteconv.Btoa(b), bits)
	*v = T(d)
	return err
}

func unmarshal(b []byte, v any) error {
	val := reflect.ValueOf(v)
	typ := val.Type()
	if dec, ok := unmarshalers.Load(typ); ok {
		return dec.(func([]byte, reflect.Value) error)(b, val) //nolint:forcetypeassert
	}
	dec, err := newUnmarshaler(typ)
	if err != nil {
		return err
	}
	unmarshalers.Store(typ, dec)
	return dec(b, val)
}

func newUnmarshaler(typ reflect.Type) (func([]byte, reflect.Value) error, error) {
	if typ.Implements(unmarshalerType) {
		isPtr := typ.Kind() == reflect.Pointer
		typ = typ.Elem()
		return func(b []byte, v reflect.Value) error {
			if len(b) == 0 {
				return nil
			}
			if isPtr && v.IsNil() {
				v.Set(reflect.New(typ))
			}
			return v.Interface().(encoding.TextUnmarshaler).UnmarshalText(b) //nolint:forcetypeassert
		}, nil
	}

	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
		dec, err := newUnmarshaler(typ)
		if err != nil {
			return nil, err
		}
		return func(b []byte, v reflect.Value) error {
			if v.IsNil() {
				v.Set(reflect.New(typ))
			}
			return dec(b, v.Elem())
		}, nil
	}

	return nil, don.ErrUnsupportedMediaType
}

var (
	unmarshalers    sync.Map
	unmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
)
