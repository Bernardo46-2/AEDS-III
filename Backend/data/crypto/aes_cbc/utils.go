package aescbc

import (
    "crypto/rand"
    "fmt"    
)

func randBytes(size int) (bytes []byte, err error) {
    bytes = make([]byte, size)
    _, err = rand.Read(bytes)
    return
}

func RotateLeft(array []byte){
    tmp := array[0]
    for i := 0; i < len(array) - 1; i++ {
        array[i] = array[i + 1]
    }
    array[len(array)-1] = tmp
}

func xorBlock(dest []byte, src []byte) {
    for i := 0; i < len(dest); i++ {
        dest[i] ^= src[i]
    }
}

func transposeMatrix(matrix []byte, lines, cols int) {
    for l := 0; l < lines; l++ {
        for c := l+1; c < cols; c++ {
            i := l * cols + c
            j := c * cols + l
            matrix[i], matrix[j] = matrix[j], matrix[i]
        }
    }
}

func prependSlice(slice, prefix []byte) []byte {
	result := make([]byte, 0, len(prefix)+len(slice))
	result = append(result, prefix...)
	return append(result, slice...)
}

func printVecHex(bytes []byte) {
    for _, b := range bytes {
        fmt.Printf("%02x ", b)
    }
    fmt.Println()
}