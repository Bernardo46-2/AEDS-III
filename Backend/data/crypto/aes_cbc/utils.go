package aescbc

import (
    "crypto/rand"
    "fmt"    
)

// randBytes gera um array de bytes aleatório usando a 
// biblioteca nativa do go para isso
func RandBytes(size int) (bytes []byte, err error) {
    bytes = make([]byte, size)
    _, err = rand.Read(bytes)
    return
}

// RotateLeft rotaciona os bytes de um array
// para a esquerda uma vez
func RotateLeft(array []byte){
    tmp := array[0]
    for i := 0; i < len(array) - 1; i++ {
        array[i] = array[i + 1]
    }
    array[len(array)-1] = tmp
}

// xorBlock recebe dois arrays, dest e src e faz
// um xor entre esses dois arrays, armazendando o
// resultado em dest
func xorBlock(dest []byte, src []byte) {
    for i := 0; i < len(dest); i++ {
        dest[i] ^= src[i]
    }
}

// transposeMatrix faz a transposição de uma matrix, para
// o algoritmo de aes, que considera a matrix como lida coluna
// por coluna, em vez de linha por linha, logo, uma transposição
// é necessaria
func transposeMatrix(matrix []byte, lines, cols int) {
    for l := 0; l < lines; l++ {
        for c := l+1; c < cols; c++ {
            i := l * cols + c
            j := c * cols + l
            matrix[i], matrix[j] = matrix[j], matrix[i]
        }
    }
}

// prependSlice recebe dois arrays e prefixa o array `prefix`
// em `slice`, retornando o array resultante
func prependSlice(slice, prefix []byte) []byte {
	result := make([]byte, 0, len(prefix)+len(slice))
	result = append(result, prefix...)
	return append(result, slice...)
}

// printVecHex é uma função para debug. Percorre o vetor
// e printa todos os valores em hexadecimal
func printVecHex(bytes []byte) {
    for _, b := range bytes {
        fmt.Printf("%02x ", b)
    }
    fmt.Println()
}
