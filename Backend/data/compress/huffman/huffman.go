package huffman

import (
	"fmt"
	"os"
)

func Zip(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("erro do tipo: %s", err.Error())
	}
}
