package huffman

import (
	"fmt"
	"os"
)

// getCharMap separa uma Mapa com todos os caracteres existentes
// preparando para a criacao do dicionario
func getCharMap(arr []byte) map[byte]int {
	charMap := make(map[byte]int)

	for _, b := range arr {
		charMap[b]++
	}

	return charMap
}

func Zip(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("erro do tipo: %s", err.Error())
	}

	charMap := getCharMap(content)
	fmt.Printf("%+v", charMap)

	return nil
}
