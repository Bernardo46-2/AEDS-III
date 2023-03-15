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

// CSV struct contendo os dados csv
type CSV struct {
	file [][]string
}

// ImportCSV abre o arquivo csv em "data/pokedex.csv" e retorna
// um ponteiro para um struct CSV contendo as linhas lidas
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

// writeBytes foi feito para facilitar a escrita do arquivo em binario.
//
// Feita mais por conveniencia
func writeBytes(file *os.File, b []byte) {
	binary.Write(file, binary.LittleEndian, b)
}

// CsvToBin converte o .csv inteiro para binario e escrever
// em um arquivo .dat
func (csv *CSV) CsvToBin() {
	// Abre o arquivo com permissão completa
	file, err := os.OpenFile(BIN_FILE, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic("Opora")
	}
	defer file.Close()

	// Inicializa o array e realiza os parsings
	arrayPokemons := []models.Pokemon{}
	for i := 1; i < len(csv.file); i++ {
		row := csv.file[i]
		pokemon := models.ParsePokemon(row)
		arrayPokemons = append(arrayPokemons, pokemon)
	}

	// Remove possiveis valores com o mesmo ID
	arrayPokemons = removeDuplicates(arrayPokemons, "Numero")

	// Grava no inicio do arquivo a quantidade de registros encontrados
	writeBytes(file, utils.IntToBytes(int32(len(arrayPokemons))))

	// Serializa e grava os registros
	for i := 0; i < len(arrayPokemons); i++ {
		bytes := arrayPokemons[i].ToBytes()
		writeBytes(file, bytes)
	}
}

// removeDuplicates remove Pokémons duplicados de uma slice dada com base no campo especificado.
// A função retorna uma nova slice com os elementos não duplicados.
//
// s: slice de Pokémons a serem removidos os duplicados
// field: string representando o campo para comparar valores
//
// Retorna uma nova slice com os elementos não duplicados
func removeDuplicates(s []models.Pokemon, field string) []models.Pokemon {
	// Cria um mapa para armazenar os valores únicos do campo
	keys := make(map[interface{}]bool)
	// Cria uma nova slice para armazenar os elementos não duplicados
	var res []models.Pokemon

	// Itera sobre a slice de Pokémons e adiciona os não duplicados à nova slice
	for _, p := range s {
		value := reflect.ValueOf(p).FieldByName(field).Interface()
		if _, ok := keys[value]; !ok {
			keys[value] = true
			res = append(res, p)
		}
	}

	// Retorna a nova slice com os elementos não duplicados
	return res
}
