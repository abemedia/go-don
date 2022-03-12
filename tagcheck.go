package don

import (
	"github.com/goccy/go-reflect"
)

func hasTag(v interface{}, tag string) bool {
	t := reflect.TypeOf(v)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return false
	}

	return typeHasTag(t, tag)
}

func typeHasTag(t reflect.Type, tag string) bool {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if f.PkgPath != "" {
			continue // skip unexported fields
		}

		if _, ok := f.Tag.Lookup(tag); ok {
			return true
		}

		ft := f.Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}

		if ft.Kind() == reflect.Struct {
			if typeHasTag(ft, tag) {
				return true
			}
		}
	}

	return false
}
