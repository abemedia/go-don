package don

import (
	"sync"
	"unsafe"

	"github.com/goccy/go-reflect"
)

type tagChecker struct {
	tag   string
	cache sync.Map
}

func newTagCheck(tag string) *tagChecker {
	return &tagChecker{tag: tag}
}

func (tc *tagChecker) Check(v interface{}) bool {
	typ := reflect.TypeOf(v).Elem()
	if typ.Kind() != reflect.Struct {
		return false
	}

	id := uintptr(unsafe.Pointer(typ))
	if res, ok := tc.cache.Load(id); ok {
		return res.(bool)
	}

	res := tc.has(typ)
	tc.cache.Store(id, res)

	return res
}

func (tc *tagChecker) has(typ reflect.Type) bool {
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)

		if f.PkgPath != "" {
			continue // skip unexported fields
		}

		if _, ok := f.Tag.Lookup(tc.tag); ok {
			return true
		}

		ft := f.Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}

		if ft.Kind() == reflect.Struct {
			if tc.has(ft) {
				return true
			}
		}
	}

	return false
}
