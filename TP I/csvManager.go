package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

const FILE string = "csv/pokedex.csv"
const BIN_FILE string = "csv/pokedex.dat"

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

func (self* CSV) getByteArray(pokemon Pokemon, size PokemonSize) ([]byte, []rune, []byte) {
    var bytes1 []byte
    var bytes2 []byte
    runes := []rune(pokemon.NomeJap)

    // bytes1 = append(bytes1, []byte(size.Total)...)
    // bytes1 = append(bytes1, []byte(pokemon.Nome)...)

    return bytes1, runes, bytes2
}

func (self* CSV) CsvToBin() {
    file, err := os.OpenFile(BIN_FILE, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
    if err != nil {
        panic("Opora")
    }
    defer file.Close()

    for i := 1; i < len(self.file); i++ {
        row := self.file[i]
        pokemon, size := ParsePokemon(row)

        b1, _, _ := self.getByteArray(pokemon, size)

        file.Write(b1)
    }
}
