package tools

import (
	"encoding/binary"
	"unsafe"
)

func Uint32ToByteArray(num uint32) []byte {
	bytes := make([]byte, int(unsafe.Sizeof(num)))

	binary.LittleEndian.PutUint32(bytes, num)
	return bytes
}

func ByteArrayToUint32(bytes []byte) uint32 {
	return binary.LittleEndian.Uint32(bytes)
}
