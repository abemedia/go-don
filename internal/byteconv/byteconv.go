// Package byteconv provides fast and efficient conversion functions for byte slices and strings.
package byteconv

import "unsafe"

// Btoa returns a string from a byte slice without memory allocation.
func Btoa(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// Atob returns a byte slice from a string without memory allocation.
func Atob(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
