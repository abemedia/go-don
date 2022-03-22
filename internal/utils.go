package internal

import "unsafe"

func Btoa(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
