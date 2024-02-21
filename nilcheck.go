package don

import "reflect"

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
		return func(v any) bool { return dataOf(v) == nil }
	case reflect.Slice:
		// Return true for nil slice.
		return func(v any) bool { return (*reflect.SliceHeader)(dataOf(v)).Data == 0 }
	default:
		// Return false for all others.
		return func(any) bool { return false }
	}
}
