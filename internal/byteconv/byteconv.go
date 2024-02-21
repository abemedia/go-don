// Package byteconv provides fast and efficient conversion functions for byte slices and strings.
package byteconv

import (
	"reflect"
	"unsafe"
)

// Btoa returns a string from a byte slice without memory allocation.
func Btoa(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// Atob returns a byte slice from a string without memory allocation.
func Atob(s string) []byte {
	sp := unsafe.Pointer(&s)
	b := *(*[]byte)(sp)
	(*reflect.SliceHeader)(unsafe.Pointer(&b)).Cap = (*reflect.StringHeader)(sp).Len
	return b
}
