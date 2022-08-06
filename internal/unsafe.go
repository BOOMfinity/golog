package internal

import (
	"reflect"
	"unsafe"
)

func GetBytes(s string) []byte {
	return (*[0x7fff0000]byte)(unsafe.Pointer(
		(*reflect.StringHeader)(unsafe.Pointer(&s)).Data),
	)[:len(s):len(s)]
}

func ToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
