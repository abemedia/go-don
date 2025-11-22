package decoder

import (
	"encoding"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/abemedia/go-don/internal/byteconv"
)

type decoder = func(reflect.Value, Getter) error

func noopDecoder(reflect.Value, Getter) error { return nil }

var unmarshalerType = reflect.TypeFor[encoding.TextUnmarshaler]()

//nolint:cyclop,funlen
func compile(typ reflect.Type, tagKey string, isPtr bool) (decoder, error) {
	decoders := []decoder{}

	for i := range typ.NumField() {
		f := typ.Field(i)
		if f.PkgPath != "" {
			continue // skip unexported fields
		}

		t, k, ptr := typeKind(f.Type)

		tag, ok := f.Tag.Lookup(tagKey)
		if !ok && k != reflect.Struct {
			continue
		}

		if reflect.PointerTo(t).Implements(unmarshalerType) {
			decoders = append(decoders, decodeTextUnmarshaler(get(ptr, i, t), tag))
			continue
		}

		switch k {
		case reflect.Struct:
			dec, err := compile(t, tagKey, ptr)
			if err != nil {
				return nil, err
			}
			index := i
			decoders = append(decoders, func(v reflect.Value, m Getter) error {
				return dec(v.Field(index), m)
			})
		case reflect.String:
			decoders = append(decoders, decodeString(set[string](ptr, i, t), tag))
		case reflect.Int:
			decoders = append(decoders, decodeInt(set[int](ptr, i, t), tag, strconv.IntSize))
		case reflect.Int8:
			decoders = append(decoders, decodeInt(set[int8](ptr, i, t), tag, 8))
		case reflect.Int16:
			decoders = append(decoders, decodeInt(set[int16](ptr, i, t), tag, 16))
		case reflect.Int32:
			decoders = append(decoders, decodeInt(set[int32](ptr, i, t), tag, 32))
		case reflect.Int64:
			decoders = append(decoders, decodeInt(set[int64](ptr, i, t), tag, 64))
		case reflect.Uint:
			decoders = append(decoders, decodeUint(set[uint](ptr, i, t), tag, strconv.IntSize))
		case reflect.Uint8:
			decoders = append(decoders, decodeUint(set[uint8](ptr, i, t), tag, 8))
		case reflect.Uint16:
			decoders = append(decoders, decodeUint(set[uint16](ptr, i, t), tag, 16))
		case reflect.Uint32:
			decoders = append(decoders, decodeUint(set[uint32](ptr, i, t), tag, 32))
		case reflect.Uint64:
			decoders = append(decoders, decodeUint(set[uint64](ptr, i, t), tag, 64))
		case reflect.Float32:
			decoders = append(decoders, decodeFloat(set[float32](ptr, i, t), tag, 32))
		case reflect.Float64:
			decoders = append(decoders, decodeFloat(set[float64](ptr, i, t), tag, 64))
		case reflect.Bool:
			decoders = append(decoders, decodeBool(set[bool](ptr, i, t), tag))
		case reflect.Slice:
			switch t.Elem().Kind() {
			case reflect.String:
				decoders = append(decoders, decodeStrings(set[[]string](ptr, i, t), tag))
			case reflect.Uint8:
				decoders = append(decoders, decodeBytes(set[[]byte](ptr, i, t), tag))
			}
		default:
			return nil, ErrUnsupportedType
		}
	}

	if len(decoders) == 0 {
		return nil, ErrTagNotFound
	}

	return func(v reflect.Value, d Getter) error {
		if isPtr {
			if v.IsNil() {
				v.Set(reflect.New(typ))
			}
			v = v.Elem()
		}

		for _, dec := range decoders {
			if err := dec(v, d); err != nil {
				return err
			}
		}

		return nil
	}, nil
}

func typeKind(t reflect.Type) (reflect.Type, reflect.Kind, bool) {
	var isPtr bool

	k := t.Kind()
	if k == reflect.Pointer {
		t = t.Elem()
		k = t.Kind()
		isPtr = true
	}

	return t, k, isPtr
}

func set[T any](ptr bool, i int, t reflect.Type) func(reflect.Value, T) {
	if ptr {
		return func(v reflect.Value, d T) {
			f := v.Field(i)
			if f.IsNil() {
				f.Set(reflect.New(t))
			}
			*(*T)(unsafe.Pointer(f.Elem().UnsafeAddr())) = d
		}
	}

	return func(v reflect.Value, d T) {
		*(*T)(unsafe.Pointer(v.Field(i).UnsafeAddr())) = d
	}
}

func get(ptr bool, i int, t reflect.Type) func(v reflect.Value) reflect.Value {
	if ptr {
		return func(v reflect.Value) reflect.Value {
			f := v.Field(i)
			if f.IsNil() {
				f.Set(reflect.New(t))
			}
			return f
		}
	}

	return func(v reflect.Value) reflect.Value {
		return v.Field(i).Addr()
	}
}

func decodeTextUnmarshaler(get func(reflect.Value) reflect.Value, k string) decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Get(k); s != "" {
			return get(v).Interface().(encoding.TextUnmarshaler).UnmarshalText(byteconv.Atob(s)) //nolint:forcetypeassert
		}
		return nil
	}
}

func decodeString(set func(reflect.Value, string), k string) decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Get(k); s != "" {
			set(v, s)
		}
		return nil
	}
}

func decodeInt[T int | int8 | int16 | int32 | int64](set func(reflect.Value, T), k string, bits int) decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Get(k); s != "" {
			n, err := strconv.ParseInt(s, 10, bits)
			if err != nil {
				return err
			}
			set(v, T(n))
		}
		return nil
	}
}

func decodeFloat[T float32 | float64](set func(reflect.Value, T), k string, bits int) decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Get(k); s != "" {
			f, err := strconv.ParseFloat(s, bits)
			if err != nil {
				return err
			}
			set(v, T(f))
		}
		return nil
	}
}

func decodeUint[T uint | uint8 | uint16 | uint32 | uint64](set func(reflect.Value, T), k string, bits int) decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Get(k); s != "" {
			n, err := strconv.ParseUint(s, 10, bits)
			if err != nil {
				return err
			}
			set(v, T(n))
		}
		return nil
	}
}

func decodeBool(set func(reflect.Value, bool), k string) decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Get(k); s != "" {
			b, err := strconv.ParseBool(s)
			if err != nil {
				return err
			}
			set(v, b)
		}
		return nil
	}
}

func decodeBytes(set func(reflect.Value, []byte), k string) decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Get(k); s != "" {
			set(v, byteconv.Atob(s))
		}
		return nil
	}
}

func decodeStrings(set func(reflect.Value, []string), k string) decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Values(k); s != nil {
			set(v, s)
		}
		return nil
	}
}
