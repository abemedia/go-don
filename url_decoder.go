package don

import (
	"sync"
	"unsafe"

	"github.com/goccy/go-reflect"
	"github.com/gorilla/schema"
)

var (
	decodeHeader = newDecoder("header")
	decodeQuery  = newDecoder("query")
	decodePath   = newDecoder("path")
)

type decoder struct {
	tag   string
	cache sync.Map
	d     *schema.Decoder
}

func newDecoder(tag string) *decoder {
	d := schema.NewDecoder()
	d.SetAliasTag(tag)
	d.IgnoreUnknownKeys(true)

	return &decoder{tag: tag, d: d}
}

func (d *decoder) Decode(dst interface{}, src map[string][]string) error {
	return d.d.Decode(dst, src)
}

func (d *decoder) Check(v interface{}) bool {
	typ := reflect.TypeOf(v).Elem()
	if typ.Kind() != reflect.Struct {
		return false
	}

	id := uintptr(unsafe.Pointer(typ))
	if res, ok := d.cache.Load(id); ok {
		return res.(bool)
	}

	res := d.has(typ)
	d.cache.Store(id, res)

	return res
}

func (d *decoder) has(typ reflect.Type) bool {
	var res bool
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)

		if f.PkgPath != "" {
			continue // skip unexported fields
		}

		_, ok := f.Tag.Lookup(d.tag)
		if ok {
			res = true
			break
		}

		if f.Type.Kind() == reflect.Ptr {
			f.Type = f.Type.Elem()
		}

		if f.Type.Kind() == reflect.Struct {
			res = d.has(f.Type)
			if res {
				return true
			}
		}
	}

	return res
}
