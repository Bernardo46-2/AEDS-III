package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

const MAX_NAME_LEN = 40

type Pokemon struct {
	Numero     int32
	Nome       string
	NomeJap    string
	Geracao    int32
	Lancamento time.Time
	Especie    string
	Lendario   bool
	Mitico     bool
	Tipo       []string
	Atk        int32
	Def        int32
	Hp         int32
	Altura     float32
	Peso       float32
	Size       PokemonSize
}

type PokemonSize struct {
	Total      int32
	Numero     int32
	Nome       int32
	NomeJap    int32
	Geracao    int32
	Lancamento int32
	Especie    int32
	Lendario   int32
	Mitico     int32
	Tipo       int32
	Atk        int32
	Def        int32
	Hp         int32
	Altura     int32
	Peso       int32
}

var GenReleaseDates = map[int]string{
	1: "1996/02/27",
	2: "1999/11/21",
	3: "2002/11/21",
	4: "2006/09/28",
	5: "2010/09/18",
	6: "2013/10/12",
	7: "2016/11/18",
	8: "2019/11/15",
	9: "2022/11/18",
}

func Atoi32(s string) (int32, error) {
	i, err := strconv.Atoi(s)
	return int32(i), err
}

func (self *Pokemon) ToString() string {
	str := ""

	str += fmt.Sprintf("\n")
	str += fmt.Sprintf("Numero     = %d\n", self.Numero)
	str += fmt.Sprintf("Nome       = %s\n", self.Nome)
	str += fmt.Sprintf("NomeJap    = %s\n", self.NomeJap)
	str += fmt.Sprintf("Geracao    = %d\n", self.Geracao)
	str += fmt.Sprintf("Lancamento = %s\n", self.Lancamento.Format("02/01/2006"))
	str += fmt.Sprintf("Especie    = %s\n", self.Especie)
	str += fmt.Sprintf("Lendario   = %t\n", self.Lendario)
	str += fmt.Sprintf("Mitico     = %t\n", self.Mitico)
	str += fmt.Sprintf("Tipo       = %s\n", self.Tipo)
	str += fmt.Sprintf("Atk        = %d\n", self.Atk)
	str += fmt.Sprintf("Def        = %d\n", self.Def)
	str += fmt.Sprintf("Hp         = %d\n", self.Hp)
	str += fmt.Sprintf("Altura     = %f\n", self.Altura)
	str += fmt.Sprintf("Peso       = %f\n", self.Peso)

	return str
}

// parseBinToPoke é responsável por interpretar os bytes de um registro binário
// contendo informações sobre um Pokémon e criar um objeto do tipo Pokemon
// a partir desses dados.
//
// A função retorna um objeto do tipo Pokemon preenchido com as informações do
// registro binário.
// Caso ocorra algum erro na leitura do arquivo binário, um erro será retornado.
func (p *Pokemon) parseBinToPoke(registro []byte) error {

	ptr := 0
	p.Numero, ptr = bytesToInt32(registro, ptr)
	p.Nome, ptr = bytesToString(registro, ptr)
	p.NomeJap, ptr = bytesToJapName(registro, ptr)
	p.Geracao, ptr = bytesToInt32(registro, ptr)
	p.Lancamento, ptr = bytesToTime(registro, ptr)
	p.Especie, ptr = bytesToString(registro, ptr)
	p.Lendario, ptr = bytesToBool(registro, ptr)
	p.Mitico, ptr = bytesToBool(registro, ptr)
	p.Tipo, ptr = bytesToArrayString(registro, ptr)
	p.Atk, ptr = bytesToInt32(registro, ptr)
	p.Def, ptr = bytesToInt32(registro, ptr)
	p.Hp, ptr = bytesToInt32(registro, ptr)
	p.Altura, ptr = bytesToFloat32(registro, ptr)
	p.Peso, ptr = bytesToFloat32(registro, ptr)

	return nil
}

func bytesToVarSize(registro []byte, ptr int) (int, int) {
	return int(binary.LittleEndian.Uint32(registro[ptr : ptr+4])), ptr + 4
}

func bytesToInt32(registro []byte, ptr int) (int32, int) {
	size, ptr := bytesToVarSize(registro, ptr)
	return int32(binary.LittleEndian.Uint32(registro[ptr : ptr+size])), ptr + size
}

func bytesToString(registro []byte, ptr int) (string, int) {
	size, ptr := bytesToVarSize(registro, ptr)
	nomeBytes := make([]byte, size)
	io.ReadFull(bytes.NewReader(registro[ptr:ptr+size]), nomeBytes)
	return strings.TrimSpace(string(nomeBytes)), ptr + size
}

func bytesToArrayString(registro []byte, ptr int) ([]string, int) {
	size, ptr := bytesToVarSize(registro, ptr)
	stringBytes := make([]byte, size)
	io.ReadFull(bytes.NewReader(registro[ptr:ptr+size]), stringBytes)
	s := strings.TrimRight(string(stringBytes), ",")
	return strings.Split(s, ","), ptr + size
}

func bytesToJapName(registro []byte, ptr int) (string, int) {
	size, ptr := bytesToVarSize(registro, ptr)

	japNameRunes := make([]rune, size/4)
	for i := 0; i < size/4; i++ {
		// Converte os 4 bytes em um uint32 correspondente ao rune.
		runeUint := binary.LittleEndian.Uint32(registro[ptr : ptr+4])
		// Converte o uint32 em um rune e adiciona à slice de runes.
		japNameRunes[i] = rune(runeUint)
		ptr += 4
	}

	return string(japNameRunes), ptr
}

func bytesToTime(registro []byte, ptr int) (time.Time, int) {
	size, ptr := bytesToVarSize(registro, ptr)
	b := make([]byte, size)
	io.ReadFull(bytes.NewReader(registro[ptr:ptr+size]), b)
	var t time.Time
	t.UnmarshalBinary(b)
	return t, ptr + size
}

func bytesToBool(registro []byte, ptr int) (bool, int) {
	_, ptr = bytesToVarSize(registro, ptr)
	if registro[ptr] != 0 {
		return true, ptr + 1
	} else {
		return false, ptr + 1
	}
}

func bytesToFloat32(registro []byte, ptr int) (float32, int) {
	size, ptr := bytesToVarSize(registro, ptr)
	bits := binary.LittleEndian.Uint32(registro[ptr : ptr+size])
	float := math.Float32frombits(bits)
	return float, ptr + size
}

func ParsePokemon(line []string) Pokemon {
	var pokemon Pokemon
	var size PokemonSize

	pokemon.Numero, _ = Atoi32(line[1])
	pokemon.Nome = line[2]
	pokemon.NomeJap = RemoveAfterSpace(line[4])
	geracao, _ := Atoi32(line[5])
	pokemon.Geracao = geracao
	pokemon.Lancamento, _ = time.Parse("2006/01/02", GenReleaseDates[int(geracao)])
	pokemon.Especie = line[9]
	pokemon.Lendario, _ = strconv.ParseBool(line[7])
	pokemon.Mitico, _ = strconv.ParseBool(line[8])
	pokemon.Tipo = append(pokemon.Tipo, line[11])
	if len(line[12]) > 0 {
		pokemon.Tipo = append(pokemon.Tipo, line[12])
	}
	pokemon.Atk, _ = Atoi32(line[21])
	pokemon.Def, _ = Atoi32(line[22])
	pokemon.Hp, _ = Atoi32(line[20])
	altura, _ := strconv.ParseFloat(line[13], 32)
	peso, _ := strconv.ParseFloat(line[14], 32)

	pokemon.Altura = float32(altura)
	pokemon.Peso = float32(peso)

	size.Numero = int32(unsafe.Sizeof(pokemon.Numero))
	size.Nome = MAX_NAME_LEN
	size.NomeJap = int32(len(pokemon.NomeJap) / 3 * 4)
	size.Geracao = int32(unsafe.Sizeof(pokemon.Geracao))

	date_size, err := pokemon.Lancamento.MarshalBinary()
	if err != nil {
		panic("Opora")
	}

	size.Lancamento = int32(len(date_size))
	size.Especie = int32(len(pokemon.Especie))
	size.Lendario = int32(unsafe.Sizeof(pokemon.Lendario))
	size.Mitico = int32(unsafe.Sizeof(pokemon.Mitico))
	size.Tipo = int32(len(pokemon.Tipo[0]) + 1)
	if len(pokemon.Tipo) > 1 {
		size.Tipo += int32(len(pokemon.Tipo[1]))
	}
	size.Atk = int32(unsafe.Sizeof(pokemon.Atk))
	size.Def = int32(unsafe.Sizeof(pokemon.Def))
	size.Hp = int32(unsafe.Sizeof(pokemon.Hp))
	size.Altura = int32(unsafe.Sizeof(pokemon.Altura))
	size.Peso = int32(unsafe.Sizeof(pokemon.Peso))

	size.Total = size.Numero + 4 +
		size.Nome + 4 +
		size.NomeJap + 4 +
		size.Geracao + 4 +
		size.Lancamento + 4 +
		size.Especie + 4 +
		size.Lendario + 4 +
		size.Mitico + 4 +
		size.Tipo + 4 +
		size.Atk + 4 +
		size.Def + 4 +
		size.Hp + 4 +
		size.Altura + 4 +
		size.Peso + 4 + 4

	pokemon.Size = size

	return pokemon
}
