package main

import (
	"encoding/csv"
    "fmt"
	"os"
	"unsafe"
)

const FILE string = "csv/pokedex2.csv"
const BIN_FILE string = "bin/pokedex.dat"

type CSV struct {
	file [][]string
}

func importCSV() CSV {
	// Abrir o arquivo CSV
	file, err := os.Open(FILE)
	if err != nil {
		panic(fmt.Errorf("Erro ao abrir o arquivo: %v", err))
	}
	defer file.Close()

	// Lendo o conteúdo do arquivo CSV
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		panic(fmt.Errorf("Erro ao ler o arquivo: %v", err))
	}

	return CSV { file: lines }
}

func (self* CSV) CsvToBin() {
    row := self.file[0]
	size := unsafe.Sizeof(row)

	// for c := range row {
		fmt.Println(row)
		fmt.Println(size)
	// }
}
