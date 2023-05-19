package huffman

import (
	"fmt"
	"os"
)

func Zip(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Erro: O arquivo '%s' não existe.\n", path)
		} else {
			fmt.Printf("Erro ao abrir o arquivo: %v\n", err)
		}
		return err
	}
}
