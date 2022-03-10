package decoder

import (
	"errors"
	"reflect"
	"strconv"
	"unsafe"
)

type Getter interface {
	Get(string) string
	Values(string) []string
}

type Decoder struct {
	dec decoder
}

func NewDecoder(v interface{}, tag string) (*Decoder, error) {
	typ := reflect.TypeOf(v)
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("not struct")
	}

	dec, err := compile(typ, tag)
	if err != nil {
		return nil, err
	}

	return &Decoder{dec}, nil
}

func (d *Decoder) Decode(data Getter, v interface{}) error {
	typ := reflect.ValueOf(v).Elem()
	if typ.Kind() != reflect.Struct {
		return errors.New("not struct")
	}

	return d.dec(typ, data)
}

type decoder func(reflect.Value, Getter) error

func compile(typ reflect.Type, tagKey string) (decoder, error) {
	decoders := []decoder{}

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		t, k, ptr := typeKind(f.Type)

		tag, ok := f.Tag.Lookup(tagKey)
		if !ok && k != reflect.Struct {
			continue
		}

		switch k {
		case reflect.Struct:
			dec, err := compile(t, tagKey)
			if err != nil {
				return nil, err
			}
			index := i

			if ptr {
				decoders = append(decoders, func(v reflect.Value, m Getter) error {
					ff := v.Field(index)
					if ff.IsNil() {
						ff.Set(reflect.New(ff.Type().Elem()))
					}
					return dec(ff.Elem(), m)
				})
			} else {
				decoders = append(decoders, func(v reflect.Value, m Getter) error {
					return dec(v.Field(index), m)
				})
			}
		case reflect.String:
			decoders = append(decoders, decodeString(set[string](ptr, i, t), tag))
		case reflect.Int:
			decoders = append(decoders, decodeInt(set[int](ptr, i, t), tag))
		case reflect.Int8:
			decoders = append(decoders, decodeInt8(set[int8](ptr, i, t), tag))
		case reflect.Int16:
			decoders = append(decoders, decodeInt16(set[int16](ptr, i, t), tag))
		case reflect.Int32:
			decoders = append(decoders, decodeInt32(set[int32](ptr, i, t), tag))
		case reflect.Int64:
			decoders = append(decoders, decodeInt64(set[int64](ptr, i, t), tag))
		case reflect.Uint:
			decoders = append(decoders, decodeUint(set[uint](ptr, i, t), tag))
		case reflect.Uint8:
			decoders = append(decoders, decodeUint8(set[uint8](ptr, i, t), tag))
		case reflect.Uint16:
			decoders = append(decoders, decodeUint16(set[uint16](ptr, i, t), tag))
		case reflect.Uint32:
			decoders = append(decoders, decodeUint32(set[uint32](ptr, i, t), tag))
		case reflect.Uint64:
			decoders = append(decoders, decodeUint64(set[uint64](ptr, i, t), tag))
		case reflect.Float32:
			decoders = append(decoders, decodeFloat32(set[float32](ptr, i, t), tag))
		case reflect.Float64:
			decoders = append(decoders, decodeFloat64(set[float64](ptr, i, t), tag))
		case reflect.Bool:
			decoders = append(decoders, decodeBool(set[bool](ptr, i, t), tag))
		case reflect.Slice:
			_, sk, _ := typeKind(t.Elem())
			switch sk {
			case reflect.String:
				decoders = append(decoders, decodeStrings(set[[]string](ptr, i, t), tag))
			case reflect.Uint8:
				decoders = append(decoders, decodeBytes(set[[]byte](ptr, i, t), tag))
			}
		default:
			return nil, errors.New("unsupported type")
		}
	}

	return func(v reflect.Value, d Getter) error {
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

func decodeString(set func(reflect.Value, string), k string) decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Get(k); s != "" {
			set(v, s)
		}
		return nil
	}
}

func decodeInt(set func(reflect.Value, int), k string) decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Get(k); s != "" {
			n, err := strconv.Atoi(s)
			if err != nil {
				return err
			}
			set(v, n)
		}
		return nil
	}
}

func decodeInt8(set func(reflect.Value, int8), k string) decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Get(k); s != "" {
			n, err := strconv.Atoi(s)
			if err != nil {
				return err
			}
			set(v, int8(n))
		}
		return nil
	}
}

func decodeInt16(set func(reflect.Value, int16), k string) decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Get(k); s != "" {
			n, err := strconv.Atoi(s)
			if err != nil {
				return err
			}
			set(v, int16(n))
		}
		return nil
	}
}

func decodeInt32(set func(reflect.Value, int32), k string) decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Get(k); s != "" {
			n, err := strconv.Atoi(s)
			if err != nil {
				return err
			}
			set(v, int32(n))
		}
		return nil
	}
}

func decodeInt64(set func(reflect.Value, int64), k string) decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Get(k); s != "" {
			n, err := strconv.Atoi(s)
			if err != nil {
				return err
			}
			set(v, int64(n))
		}
		return nil
	}
}

func decodeFloat32(set func(reflect.Value, float32), k string) decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Get(k); s != "" {
			f, err := strconv.ParseFloat(s, 32)
			if err != nil {
				return err
			}
			set(v, float32(f))
		}
		return nil
	}
}

func decodeFloat64(set func(reflect.Value, float64), k string) decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Get(k); s != "" {
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return err
			}
			set(v, f)
		}
		return nil
	}
}

func decodeUint(set func(reflect.Value, uint), k string) decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Get(k); s != "" {
			n, err := strconv.ParseUint(s, 10, strconv.IntSize)
			if err != nil {
				return err
			}
			set(v, uint(n))
		}
		return nil
	}
}

func decodeUint8(set func(reflect.Value, uint8), k string) decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Get(k); s != "" {
			n, err := strconv.ParseUint(s, 10, 8)
			if err != nil {
				return err
			}
			set(v, uint8(n))
		}
		return nil
	}
}

func decodeUint16(set func(reflect.Value, uint16), k string) decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Get(k); s != "" {
			n, err := strconv.ParseUint(s, 10, 16)
			if err != nil {
				return err
			}
			set(v, uint16(n))
		}
		return nil
	}
}

func decodeUint32(set func(reflect.Value, uint32), k string) decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Get(k); s != "" {
			n, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				return err
			}
			set(v, uint32(n))
		}
		return nil
	}
}

func decodeUint64(set func(reflect.Value, uint64), k string) decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Get(k); s != "" {
			n, err := strconv.ParseUint(s, 10, 64)
			if err != nil {
				return err
			}
			set(v, n)
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
			sp := unsafe.Pointer(&s)
			b := *(*[]byte)(sp)
			(*reflect.SliceHeader)(unsafe.Pointer(&b)).Cap = (*reflect.StringHeader)(sp).Len
			set(v, b)
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
