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

			decoders = append(decoders, func(v reflect.Value, m Getter) error {
				ff := v.Field(index)
				if ptr && ff.IsNil() {
					ff.Set(reflect.New(ff.Type().Elem()))
				}
				return dec(reflect.Indirect(ff), m)
			})
		case reflect.String:
			decoders = append(decoders, decodeString(i, tag))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			decoders = append(decoders, decodeInt(i, tag))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			decoders = append(decoders, decodeUint(i, tag))
		case reflect.Float32, reflect.Float64:
			decoders = append(decoders, decodeFloat(i, tag))
		case reflect.Bool:
			decoders = append(decoders, decodeBool(i, tag))
		case reflect.Slice:
			_, sk, _ := typeKind(t.Elem())
			switch sk {
			case reflect.String:
				decoders = append(decoders, decodeStrings(i, tag))
			case reflect.Uint8:
				decoders = append(decoders, decodeBytes(i, tag))
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
	if k == reflect.Ptr {
		t = t.Elem()
		k = t.Kind()
		isPtr = true
	}
	return t, k, isPtr
}

func decodeString(i int, k string) decoder {
	return func(v reflect.Value, m Getter) error {
		if d := m.Get(k); d != "" {
			v.Field(i).SetString(d)
		}
		return nil
	}
}

func decodeInt(i int, k string) decoder {
	return func(v reflect.Value, m Getter) error {
		if d := m.Get(k); d != "" {
			n, err := strconv.ParseInt(d, 10, 64)
			if err != nil {
				return err
			}
			v.Field(i).SetInt(n)
		}
		return nil
	}
}

func decodeFloat(i int, k string) decoder {
	return func(v reflect.Value, m Getter) error {
		if d := m.Get(k); d != "" {
			n, err := strconv.ParseFloat(d, 64)
			if err != nil {
				return err
			}
			v.Field(i).SetFloat(n)
		}
		return nil
	}
}

func decodeUint(i int, k string) decoder {
	return func(v reflect.Value, m Getter) error {
		if d := m.Get(k); d != "" {
			n, err := strconv.ParseUint(d, 10, 64)
			if err != nil {
				return err
			}
			v.Field(i).SetUint(n)
		}
		return nil
	}
}

func decodeBool(i int, k string) decoder {
	return func(v reflect.Value, m Getter) error {
		if d := m.Get(k); d != "" {
			n, err := strconv.ParseBool(d)
			if err != nil {
				return err
			}
			v.Field(i).SetBool(n)
		}
		return nil
	}
}

func decodeBytes(i int, k string) decoder {
	return func(v reflect.Value, m Getter) error {
		if d := m.Get(k); d != "" {
			*(*[]byte)(unsafe.Pointer(v.Field(i).UnsafeAddr())) = []byte(d)
		}
		return nil
	}
}

func decodeStrings(i int, k string) decoder {
	return func(v reflect.Value, m Getter) error {
		if d := m.Values(k); d != nil {
			*(*[]string)(unsafe.Pointer(v.Field(i).UnsafeAddr())) = d
		}
		return nil
	}
}
