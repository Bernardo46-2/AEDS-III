package main

import (
    "encoding/binary"
    "math"
)

// Funcao para converter int32 para []byte
func IntToBytes(n int32) []byte {
	var buf []byte
	return binary.LittleEndian.AppendUint32(buf, uint32(n))
}

// Funcao para converter float32 para []byte
func FloatToBytes(f float32) []byte {
	b := make([]byte, 4)
	bits := math.Float32bits(f)
	binary.LittleEndian.PutUint32(b, bits)

	return b
}