package dataManager

import (
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"os"
	"reflect"

	"github.com/Bernardo46-2/AEDS-III/models"
	"github.com/Bernardo46-2/AEDS-III/utils"
)

const FILE string = "data/pokedex.csv"
const BIN_FILE string = "data/pokedex.dat"

type CSV struct {
	file [][]string
}

// Funcao para importar o arquivo .csv para uma matriz
// de strings
func ImportCSV() *CSV {
	// Abrir o arquivo CSV
	file, err := os.Open(FILE)
	if err != nil {
		panic(fmt.Errorf("erro ao abrir o arquivo: %v", err))
	}
	defer file.Close()

	// Lendo o conteúdo do arquivo CSV
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		panic(fmt.Errorf("erro ao ler o arquivo: %v", err))
	}

	return &CSV{file: lines}
}

// Funcao para facilitar a escrita do arquivo em binario
// Feita mais por conveniencia
func writeBytes(file *os.File, b []byte) {
	binary.Write(file, binary.LittleEndian, b)
}

// Funcao para converter o .csv inteiro para binario e escrever
// em um arquivo .dat
func (csv *CSV) CsvToBin() {
	file, err := os.OpenFile(BIN_FILE, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic("Opora")
	}
	defer file.Close()

	arrayPokemons := []models.Pokemon{}
	for i := 1; i < len(csv.file); i++ {
		row := csv.file[i]
		pokemon := models.ParsePokemon(row)
		arrayPokemons = append(arrayPokemons, pokemon)
	}

	arrayPokemons = removeDuplicates(arrayPokemons, "Numero")

	writeBytes(file, utils.IntToBytes(int32(len(arrayPokemons))))

	for i := 0; i < len(arrayPokemons); i++ {
		bytes := arrayPokemons[i].ToBytes()
		writeBytes(file, bytes)
	}
}

func removeDuplicates(s []models.Pokemon, field string) []models.Pokemon {
	keys := make(map[interface{}]bool)
	var res []models.Pokemon

	for _, p := range s {
		value := reflect.ValueOf(p).FieldByName(field).Interface()
		if _, ok := keys[value]; !ok {
			keys[value] = true
			res = append(res, p)
		}
	}

	return res
}
