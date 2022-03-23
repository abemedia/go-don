package don

import (
	"unsafe"

	"github.com/goccy/go-reflect"
)

func makeNilCheck(zero any) func(v any) bool {
	// Return true for nil interfaces.
	if zero == *new(any) {
		return func(v any) bool { return v == nil }
	}

	// Return true for don.Empty.
	if _, ok := zero.(Empty); ok {
		return func(v any) bool { return true }
	}

	switch reflect.TypeOf(zero).Kind() {
	case reflect.String, reflect.Ptr:
		// Return true for empty strings and nil pointer.
		return func(v any) bool { return v == zero }
	case reflect.Map:
		// Return true for and nil map.
		return func(v any) bool {
			return (*emptyInterface)(unsafe.Pointer(&v)).ptr == nil
		}
	case reflect.Slice:
		// Return true for and nil slice.
		return func(v any) bool {
			header := (*reflect.SliceHeader)((*emptyInterface)(unsafe.Pointer(&v)).ptr)
			return header.Data == 0
		}
	default:
		// Return false for all others.
		return func(v any) bool { return false }
	}
}

// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
	typ unsafe.Pointer
	ptr unsafe.Pointer
}
