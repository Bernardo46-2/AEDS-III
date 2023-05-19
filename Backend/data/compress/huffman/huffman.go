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
	defer file.Close()

	content := make([]byte, 0)
	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			break
		}
		content = append(content, buffer[:n]...)
	}

	fmt.Println(string(content))

	return nil
}
