package don

import "unsafe"

//go:linkname unsafe_New reflect.unsafe_New
func unsafe_New(unsafe.Pointer) unsafe.Pointer //nolint:revive

//go:linkname typedmemmove reflect.typedmemmove
func typedmemmove(t, dst, src unsafe.Pointer)

// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
	typ unsafe.Pointer
	ptr unsafe.Pointer
}

func dataOf(v any) unsafe.Pointer {
	return (*emptyInterface)(unsafe.Pointer(&v)).ptr
}

func packEface(typ, ptr unsafe.Pointer) any {
	var i any
	e := (*emptyInterface)(unsafe.Pointer(&i))
	e.typ = typ
	e.ptr = ptr
	return i
}
