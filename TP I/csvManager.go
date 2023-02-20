package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"unsafe"
	// "unsafe"
)

const FILE string = "csv/pokedex.csv"
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
    // size := unsafe.Sizeof(row)
    var c rune

    fmt.Println(unsafe.Sizeof(c))

    for i := 1; i < len(self.file); i++ {
        row := self.file[i]
        pokemon := ParsePokemon(row)
        
        fmt.Println(pokemon.ToString())
    }
}
