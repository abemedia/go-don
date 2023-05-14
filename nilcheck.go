package don

import (
	"reflect"
	"unsafe"
)

func newNilCheck(zero any) func(v any) bool {
	typ := reflect.TypeOf(zero)

	// Return true for nil interface.
	if typ == nil {
		return func(v any) bool { return v == nil }
	}

	switch typ.Kind() {
	case reflect.String, reflect.Ptr, reflect.Interface:
		// Return true for empty string and nil pointer.
		return func(v any) bool { return v == zero }
	case reflect.Map:
		// Return true for and nil map.
		return func(v any) bool {
			return (*emptyInterface)(unsafe.Pointer(&v)).ptr == nil
		}
	case reflect.Slice:
		// Return true for nil slice.
		return func(v any) bool {
			return (*reflect.SliceHeader)((*emptyInterface)(unsafe.Pointer(&v)).ptr).Data == 0
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
