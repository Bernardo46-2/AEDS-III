package main

import (
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"math"
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

	return CSV{file: lines}
}

func intToBytes(n int32) []byte {
	var buf []byte
	return binary.LittleEndian.AppendUint32(buf, uint32(n))
}

func floatToBytes(f float32) []byte {
	// Create a new byte slice with 4 bytes
	b := make([]byte, 4)

	// Use the math package's Float32bits function to get the binary representation of the float32
	bits := math.Float32bits(f)

	// Convert the 32-bit float to a 4-byte slice
	binary.LittleEndian.PutUint32(b, bits)

	return b
}

// TODO during ES class...
func writeBytes(file *os.File, b []byte) {

}

func (self *CSV) CsvToBin() {
	file, err := os.OpenFile(BIN_FILE, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic("Opora")
	}
	defer file.Close()

	binary.Write(file, binary.LittleEndian, intToBytes(int32(len(self.file))))

	for i := 1; i < len(self.file); i++ {
		row := self.file[i]
		pokemon := ParsePokemon(row)

        binary.Write(file, binary.LittleEndian, intToBytes(1))
		binary.Write(file, binary.LittleEndian, intToBytes(pokemon.Size.Total))

		binary.Write(file, binary.LittleEndian, intToBytes(pokemon.Size.Numero))
		binary.Write(file, binary.LittleEndian, intToBytes(pokemon.Numero))

		binary.Write(file, binary.LittleEndian, intToBytes(pokemon.Size.Nome))
		binary.Write(file, binary.LittleEndian, []byte(pokemon.Nome))

		filler := make([]byte, pokemon.Size.Nome-int32(len(pokemon.Nome)))
		binary.Write(file, binary.LittleEndian, filler)

		runes := []rune(pokemon.NomeJap)
		japName := make([]byte, len(runes)*4)

		for i, v := range runes {
			binary.LittleEndian.PutUint32(japName[i*4:(i+1)*4], uint32(v))
		}

		binary.Write(file, binary.LittleEndian, intToBytes(int32(len(runes)*4)))
		binary.Write(file, binary.LittleEndian, japName)

		binary.Write(file, binary.LittleEndian, intToBytes(pokemon.Size.Geracao))
		binary.Write(file, binary.LittleEndian, intToBytes(pokemon.Geracao))

		releaseDate, _ := pokemon.Lancamento.MarshalBinary()

		binary.Write(file, binary.LittleEndian, intToBytes(pokemon.Size.Lancamento))
		binary.Write(file, binary.LittleEndian, releaseDate)

		binary.Write(file, binary.LittleEndian, intToBytes(pokemon.Size.Especie))
		binary.Write(file, binary.LittleEndian, []byte(pokemon.Especie))

		binary.Write(file, binary.LittleEndian, intToBytes(pokemon.Size.Lendario))
		binary.Write(file, binary.LittleEndian, pokemon.Lendario)

		binary.Write(file, binary.LittleEndian, intToBytes(pokemon.Size.Mitico))
		binary.Write(file, binary.LittleEndian, pokemon.Mitico)

		binary.Write(file, binary.LittleEndian, intToBytes(pokemon.Size.Tipo))
		binary.Write(file, binary.LittleEndian, []byte(pokemon.Tipo[0]+","))

		if len(pokemon.Tipo) > 1 {
			binary.Write(file, binary.LittleEndian, []byte(pokemon.Tipo[1]))
		}

		binary.Write(file, binary.LittleEndian, intToBytes(int32(pokemon.Size.Atk)))
		binary.Write(file, binary.LittleEndian, intToBytes(int32(pokemon.Atk)))

		binary.Write(file, binary.LittleEndian, intToBytes(int32(pokemon.Size.Def)))
		binary.Write(file, binary.LittleEndian, intToBytes(int32(pokemon.Def)))

		binary.Write(file, binary.LittleEndian, intToBytes(int32(pokemon.Size.Hp)))
		binary.Write(file, binary.LittleEndian, intToBytes(int32(pokemon.Hp)))

		binary.Write(file, binary.LittleEndian, intToBytes(int32(pokemon.Size.Altura)))
		binary.Write(file, binary.LittleEndian, floatToBytes(pokemon.Altura))

		binary.Write(file, binary.LittleEndian, intToBytes(int32(pokemon.Size.Peso)))
		binary.Write(file, binary.LittleEndian, floatToBytes(pokemon.Peso))
	}
}
