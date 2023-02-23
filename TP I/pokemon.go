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

// bytesToVarSize é responsável por extrair um valor de tamanho variável
// de um registro binário e retornar seu valor e a próxima posição no registro.
//
// registro é um slice de bytes representando um registro binário contendo a
// informação a ser extraída.
//
// ptr é um inteiro representando a posição atual do ponteiro de leitura no
// registro.
//
// A função retorna um inteiro representando o valor extraído e um inteiro
// representando a próxima posição do ponteiro de leitura no registro.
func bytesToVarSize(registro []byte, ptr int) (int, int) {
	return int(binary.LittleEndian.Uint32(registro[ptr : ptr+4])), ptr + 4
}

// bytesToInt32 é responsável por extrair um valor inteiro de 32 bits
// de um registro binário e retornar seu valor e a próxima posição no registro.
//
// registro é um slice de bytes representando um registro binário contendo a
// informação a ser extraída.
//
// ptr é um inteiro representando a posição atual do ponteiro de leitura no
// registro.
//
// A função retorna um inteiro de 32 bits representando o valor extraído e um
// inteiro representando a próxima posição do ponteiro de leitura no registro.
func bytesToInt32(registro []byte, ptr int) (int32, int) {
	size, ptr := bytesToVarSize(registro, ptr)
	return int32(binary.LittleEndian.Uint32(registro[ptr : ptr+size])), ptr + size
}

// bytesToString é responsável por extrair uma string de um registro binário
// e retornar seu valor e a próxima posição no registro.
//
// registro é um slice de bytes representando um registro binário contendo a
// informação a ser extraída.
//
// ptr é um inteiro representando a posição atual do ponteiro de leitura no
// registro.
//
// A função retorna uma string representando o valor extraído e um
// inteiro representando a próxima posição do ponteiro de leitura no registro.
func bytesToString(registro []byte, ptr int) (string, int) {
	size, ptr := bytesToVarSize(registro, ptr)
	nomeBytes := make([]byte, size)
	io.ReadFull(bytes.NewReader(registro[ptr:ptr+size]), nomeBytes)
	return strings.TrimSpace(string(nomeBytes)), ptr + size
}

// bytesToArrayString é responsável por extrair um array de strings com
// tabulação em ',' de um registro binário e retornar seu valor e a próxima
// posição no registro.
//
// registro é um slice de bytes representando um registro binário contendo a
// informação a ser extraída.
//
// ptr é um inteiro representando a posição atual do ponteiro de leitura no
// registro.
//
// A função retorna um slice de strings representando o valor extraído e um
// inteiro representando a próxima posição do ponteiro de leitura no registro.
func bytesToArrayString(registro []byte, ptr int) ([]string, int) {
	size, ptr := bytesToVarSize(registro, ptr)
	stringBytes := make([]byte, size)
	io.ReadFull(bytes.NewReader(registro[ptr:ptr+size]), stringBytes)
	s := strings.TrimRight(string(stringBytes), ",")
	return strings.Split(s, ","), ptr + size
}

// bytesToJapName é responsável por converter um slice de bytes
// que representa um nome em japonês para uma string.
//
// A função recebe como argumentos um slice de bytes e um ponteiro para
// uma posição no slice. A partir do ponteiro, a função lê os bytes
// necessários para obter o tamanho do nome japonês. Em seguida, a função
// lê os bytes correspondentes aos runes do nome, converte-os em runes
// e adiciona à slice de runes japNameRunes. Finalmente, a função retorna
// a string criada a partir da slice de runes japNameRunes e o ponteiro
// atualizado.
//
// A função retorna uma string e um inteiro representando a posição
// atual no slice de bytes. Se ocorrer algum erro durante a leitura do
// arquivo binário, a função retornará uma string vazia e o último
// valor do ponteiro.
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

// bytesToTime decodifica um valor de tipo time.Time a partir de um registro de bytes e um ponteiro
// para a posição do próximo dado no registro.
//
// registro: um registro de bytes a ser decodificado.
// ptr: um inteiro que representa a posição do próximo dado no registro.
//
// Retorna um valor do tipo time.Time decodificado e o novo valor de ptr que representa a
// posição do próximo dado no registro.
func bytesToTime(registro []byte, ptr int) (time.Time, int) {
	size, ptr := bytesToVarSize(registro, ptr)
	b := make([]byte, size)
	io.ReadFull(bytes.NewReader(registro[ptr:ptr+size]), b)
	var t time.Time
	t.UnmarshalBinary(b)
	return t, ptr + size
}

// bytesToBool converte um booleano representado em bytes para um valor bool.
// Recebe um slice de bytes 'registro' que contém os dados do booleano e um inteiro
// 'ptr' que aponta para a posição atual no registro. Retorna um valor bool e um
// inteiro representando o novo ponteiro para o registro após a conversão.
func bytesToBool(registro []byte, ptr int) (bool, int) {
	_, ptr = bytesToVarSize(registro, ptr)
	if registro[ptr] != 0 {
		return true, ptr + 1
	} else {
		return false, ptr + 1
	}
}

// bytesToFloat32 converte um slice de bytes em um valor float32 e retorna o valor
// convertido e o novo ponteiro.
//
// A função espera que o slice de bytes fornecido comece com uma sequência de bytes
// que representam o tamanho do valor a ser convertido. Em seguida, converte os
// bytes restantes do slice para um uint32, e usa a função math.Float32frombits
// para converter o uint32 em um valor float32. O novo ponteiro retornado aponta para
// o byte seguinte ao final da sequência de bytes que representam o valor float32.
//
// registro é o slice de bytes a ser convertido.
// ptr é o ponteiro inicial do registro onde a conversão deve começar.
//
// A função retorna o valor float32 convertido e o novo ponteiro após a conversão.
func bytesToFloat32(registro []byte, ptr int) (float32, int) {
	size, ptr := bytesToVarSize(registro, ptr)
	bits := binary.LittleEndian.Uint32(registro[ptr : ptr+size])
	float := math.Float32frombits(bits)
	return float, ptr + size
}

func parsePokemon(line []string) Pokemon {
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

func readPokemon() Pokemon {
	var p Pokemon
	prompt := ""
	p.Numero, prompt = lerInt32("Numero da pokedex", prompt)
	p.Nome, prompt = lerString("Nome", prompt)
	p.NomeJap, prompt = lerString("Nome Japones", prompt)
	p.Geracao, prompt = lerInt32("Geraçao", prompt, len(GenReleaseDates))
	p.Lancamento, _ = time.Parse("2006/01/02", GenReleaseDates[int(p.Geracao)])
	fmt.Printf("Data da geracao = %s\n", GenReleaseDates[int(p.Geracao)])
	prompt += "Data da geracao = " + GenReleaseDates[int(p.Geracao)] + "\n"
	p.Especie, prompt = lerString("Especie", prompt)
	p.Lendario, prompt = lerBool("É Lendario", prompt)
	p.Mitico, prompt = lerBool("É Mitico", prompt)
	p.Tipo, prompt = lerStringSlice("Tipo do pokemon", prompt, 2)
	p.Atk, prompt = lerInt32("Atk", prompt)
	p.Def, prompt = lerInt32("Def", prompt)
	p.Hp, prompt = lerInt32("Hp", prompt)
	p.Altura, prompt = lerFloat32("Altura", prompt)
	p.Peso, prompt = lerFloat32("Peso", prompt)

	return p
}
