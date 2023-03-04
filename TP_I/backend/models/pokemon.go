package models

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"time"
	"unsafe"

	"github.com/Bernardo46-2/AEDS-III/utils"
)

const MAX_NAME_LEN = 40

type PokemonID struct {
	ID int `json:"id"`
}

type Pokemon struct {
	Numero     int32     `json:"numero"`
	Nome       string    `json:"nome,omitempty"`
	NomeJap    string    `json:"nome_jap,omitempty"`
	Geracao    int32     `json:"geracao"`
	Lancamento time.Time `json:"lancamento"`
	Especie    string    `json:"especie"`
	Lendario   bool      `json:"lendario"`
	Mitico     bool      `json:"mitico"`
	Tipo       []string  `json:"tipo"`
	Atk        int32     `json:"atk"`
	Def        int32     `json:"def"`
	Hp         int32     `json:"hp"`
	Altura     float32   `json:"altura"`
	Peso       float32   `json:"peso"`
	Size       PokeSize  `json:"-"`
}

type PokeSize struct {
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

func (p *Pokemon) ToString() string {
	str := ""

	str += fmt.Sprintf("Numero     = %d\n", p.Numero)
	str += fmt.Sprintf("Nome       = %s\n", p.Nome)
	str += fmt.Sprintf("NomeJap    = %s\n", p.NomeJap)
	str += fmt.Sprintf("Geracao    = %d\n", p.Geracao)
	str += fmt.Sprintf("Lancamento = %s\n", p.Lancamento.Format("02/01/2006"))
	str += fmt.Sprintf("Especie    = %s\n", p.Especie)
	str += fmt.Sprintf("Lendario   = %t\n", p.Lendario)
	str += fmt.Sprintf("Mitico     = %t\n", p.Mitico)
	str += fmt.Sprintf("Tipo       = %s\n", p.Tipo)
	str += fmt.Sprintf("Atk        = %d\n", p.Atk)
	str += fmt.Sprintf("Def        = %d\n", p.Def)
	str += fmt.Sprintf("Hp         = %d\n", p.Hp)
	str += fmt.Sprintf("Altura     = %f\n", p.Altura)
	str += fmt.Sprintf("Peso       = %f\n", p.Peso)

	return str
}

func copyBytes(dest []byte, src []byte, offset int) ([]byte, int) {
	copy(dest[offset:], src)
	return dest, offset + len(src)
}

func (p *Pokemon) ToBytes() []byte {
	pokeBytes := make([]byte, p.Size.Total+4)
	var lendario, mitico []byte
	offset := 0
	valid := 1

	if p.Lendario {
		lendario = []byte{1}
	} else {
		lendario = []byte{0}
	}

	if p.Mitico {
		mitico = []byte{1}
	} else {
		mitico = []byte{0}
	}

	releaseDate, _ := p.Lancamento.MarshalBinary()
	filler := make([]byte, p.Size.Nome-int32(len(p.Nome)))
	runes := []rune(p.NomeJap)
	japName := make([]byte, len(runes)*4)

	for i, v := range runes {
		binary.LittleEndian.PutUint32(japName[i*4:(i+1)*4], uint32(v))
	}

	pokeBytes, offset = copyBytes(pokeBytes, utils.IntToBytes(int32(valid)), offset)
	pokeBytes, offset = copyBytes(pokeBytes, utils.IntToBytes(p.Size.Total), offset)

	pokeBytes, offset = copyBytes(pokeBytes, utils.IntToBytes(p.Numero), offset)

	pokeBytes, offset = copyBytes(pokeBytes, []byte(p.Nome), offset)
	pokeBytes, offset = copyBytes(pokeBytes, filler, offset)

	pokeBytes, offset = copyBytes(pokeBytes, utils.IntToBytes(int32(len(runes)*4)), offset)
	pokeBytes, offset = copyBytes(pokeBytes, japName, offset)

	pokeBytes, offset = copyBytes(pokeBytes, utils.IntToBytes(p.Geracao), offset)

	pokeBytes, offset = copyBytes(pokeBytes, utils.IntToBytes(p.Size.Lancamento), offset)
	pokeBytes, offset = copyBytes(pokeBytes, releaseDate, offset)

	pokeBytes, offset = copyBytes(pokeBytes, utils.IntToBytes(p.Size.Especie), offset)
	pokeBytes, offset = copyBytes(pokeBytes, []byte(p.Especie), offset)

	pokeBytes, offset = copyBytes(pokeBytes, lendario, offset)
	pokeBytes, offset = copyBytes(pokeBytes, mitico, offset)

	pokeBytes, offset = copyBytes(pokeBytes, utils.IntToBytes(p.Size.Tipo), offset)
	pokeBytes, offset = copyBytes(pokeBytes, []byte(p.Tipo[0]+","), offset)

	if len(p.Tipo) > 1 {
		pokeBytes, offset = copyBytes(pokeBytes, []byte(p.Tipo[1]), offset)
	}

	pokeBytes, offset = copyBytes(pokeBytes, utils.IntToBytes(p.Atk), offset)
	pokeBytes, offset = copyBytes(pokeBytes, utils.IntToBytes(p.Def), offset)
	pokeBytes, offset = copyBytes(pokeBytes, utils.IntToBytes(p.Hp), offset)
	pokeBytes, offset = copyBytes(pokeBytes, utils.FloatToBytes(p.Altura), offset)
	pokeBytes, _ = copyBytes(pokeBytes, utils.FloatToBytes(p.Peso), offset)

	return pokeBytes
}

func (p *Pokemon) ParseBinToPoke(registro []byte) error {
	ptr := 0

	p.Numero, ptr = utils.BytesToInt32(registro, ptr)
	p.Nome, ptr = utils.BytesToFixedSizeString(registro, ptr, MAX_NAME_LEN)
	p.NomeJap, ptr = utils.BytesToJapName(registro, ptr)
	p.Geracao, ptr = utils.BytesToInt32(registro, ptr)
	p.Lancamento, ptr = utils.BytesToTime(registro, ptr)
	p.Especie, ptr = utils.BytesToString(registro, ptr)
	p.Lendario, ptr = utils.BytesToBool(registro, ptr)
	p.Mitico, ptr = utils.BytesToBool(registro, ptr)
	p.Tipo, ptr = utils.BytesToArrayString(registro, ptr)
	p.Atk, ptr = utils.BytesToInt32(registro, ptr)
	p.Def, ptr = utils.BytesToInt32(registro, ptr)
	p.Hp, ptr = utils.BytesToInt32(registro, ptr)
	p.Altura, ptr = utils.BytesToFloat32(registro, ptr)
	p.Peso, _ = utils.BytesToFloat32(registro, ptr)
	p.CalculateSize()

	return nil
}

func ParsePokemon(line []string) Pokemon {
	var pokemon Pokemon

	pokemon.Numero, _ = utils.Atoi32(line[1])
	pokemon.Nome = line[2]
	pokemon.NomeJap = utils.RemoveAfterSpace(line[4])
	geracao, _ := utils.Atoi32(line[5])
	pokemon.Geracao = geracao
	pokemon.Lancamento, _ = time.Parse("2006/01/02", GenReleaseDates[int(geracao)])
	pokemon.Especie = line[9]
	pokemon.Lendario, _ = strconv.ParseBool(line[7])
	pokemon.Mitico, _ = strconv.ParseBool(line[8])
	pokemon.Tipo = append(pokemon.Tipo, line[11])
	if len(line[12]) > 0 {
		pokemon.Tipo = append(pokemon.Tipo, line[12])
	}
	pokemon.Atk, _ = utils.Atoi32(line[21])
	pokemon.Def, _ = utils.Atoi32(line[22])
	pokemon.Hp, _ = utils.Atoi32(line[20])
	altura, _ := strconv.ParseFloat(line[13], 32)
	peso, _ := strconv.ParseFloat(line[14], 32)

	pokemon.Altura = float32(altura)
	pokemon.Peso = float32(peso)

	pokemon.CalculateSize()

	return pokemon
}

func (p *Pokemon) CalculateSize() {
	p.Size.Numero = int32(unsafe.Sizeof(p.Numero))
	p.Size.Nome = MAX_NAME_LEN
	p.Size.NomeJap = int32(len(p.NomeJap) / 3 * 4)
	p.Size.Geracao = int32(unsafe.Sizeof(p.Geracao))

	date_size, err := p.Lancamento.MarshalBinary()
	if err != nil {
		panic("Opora")
	}

	p.Size.Lancamento = int32(len(date_size))
	p.Size.Especie = int32(len(p.Especie))
	p.Size.Lendario = int32(unsafe.Sizeof(p.Lendario))
	p.Size.Mitico = int32(unsafe.Sizeof(p.Mitico))
	p.Size.Tipo = int32(len(p.Tipo[0]) + 1)
	if len(p.Tipo) > 1 {
		p.Size.Tipo += int32(len(p.Tipo[1]))
	}
	p.Size.Atk = int32(unsafe.Sizeof(p.Atk))
	p.Size.Def = int32(unsafe.Sizeof(p.Def))
	p.Size.Hp = int32(unsafe.Sizeof(p.Hp))
	p.Size.Altura = int32(unsafe.Sizeof(p.Altura))
	p.Size.Peso = int32(unsafe.Sizeof(p.Peso))

	p.Size.Total = p.Size.Numero + 4 +
		MAX_NAME_LEN +
		p.Size.NomeJap + 4 +
		p.Size.Geracao +
		p.Size.Lancamento + 4 +
		p.Size.Especie + 4 +
		p.Size.Lendario +
		p.Size.Mitico +
		p.Size.Tipo + 4 +
		p.Size.Atk +
		p.Size.Def +
		p.Size.Hp +
		p.Size.Altura +
		p.Size.Peso + 4 + 1
}

func ReadPokemon() Pokemon {
	var p Pokemon
	var tmpNomeJap string
	prompt := ""
	p.Numero, prompt = utils.LerInt32("Numero da pokedex", prompt)
	p.Nome, prompt = utils.LerString("Nome", prompt)
	tmpNomeJap, prompt = utils.LerString("Nome Japones", prompt)
	p.NomeJap = utils.ToKatakana(tmpNomeJap)
	fmt.Printf("Conversao para japones: %s\n", p.NomeJap)
	prompt += "Conversao para japones: " + p.NomeJap + "\n"
	p.Geracao, prompt = utils.LerInt32("Geraçao", prompt, len(GenReleaseDates))
	p.Lancamento, _ = time.Parse("2006/01/02", GenReleaseDates[int(p.Geracao)])
	fmt.Printf("Data da geracao = %s\n", GenReleaseDates[int(p.Geracao)])
	prompt += "Data da geracao = " + GenReleaseDates[int(p.Geracao)] + "\n"
	p.Especie, prompt = utils.LerString("Especie", prompt)
	p.Lendario, prompt = utils.LerBool("É Lendario", prompt)
	p.Mitico, prompt = utils.LerBool("É Mitico", prompt)
	p.Tipo, prompt = utils.LerStringSlice("Tipo do pokemon", prompt, 2)
	p.Atk, prompt = utils.LerInt32("Atk", prompt)
	p.Def, prompt = utils.LerInt32("Def", prompt)
	p.Hp, prompt = utils.LerInt32("Hp", prompt)
	p.Altura, prompt = utils.LerFloat32("Altura", prompt)
	p.Peso, _ = utils.LerFloat32("Peso", prompt)
	p.CalculateSize()

	return p
}

func (p *Pokemon) AlterarCampo() {
	continuar := true
	for continuar {
		prompt := ""
		campo, prompt := utils.LerString(p.ToString()+"\nQual campo quer alterar", prompt)
		switch campo {
		case "Numero":
			p.Numero, prompt = utils.LerInt32("Novo Numero", prompt)
		case "Nome":
			p.Nome, prompt = utils.LerString("Novo Nome", prompt)
		case "NomeJap":
			tmpNomeJap, prompt := utils.LerString("Novo Nome Japones", prompt)
			p.NomeJap = utils.ToKatakana(tmpNomeJap)
			fmt.Printf("Conversao para japones: %s\n", p.NomeJap)
			prompt += "Conversao para japones: " + p.NomeJap + "\n"
		case "Geracao":
			p.Geracao, prompt = utils.LerInt32("Nova Geracao", prompt)
		case "Lancamento":
			p.Lancamento, prompt = utils.LerTime("Novo Data de Lancamento", prompt)
		case "Especie":
			p.Especie, prompt = utils.LerString("Nova Especie", prompt)
		case "Lendario":
			p.Lendario, prompt = utils.LerBool("Novo Status Lendario", prompt)
		case "Mitico":
			p.Mitico, prompt = utils.LerBool("Novo Status Mitico", prompt)
		case "Tipo":
			p.Tipo, prompt = utils.LerStringSlice("Novo Tipo", prompt, 2)
		case "Atk":
			p.Atk, prompt = utils.LerInt32("Novo Atk", prompt)
		case "Def":
			p.Def, prompt = utils.LerInt32("Nova Def", prompt)
		case "Hp":
			p.Hp, prompt = utils.LerInt32("Novo Hp", prompt)
		case "Altura":
			p.Altura, prompt = utils.LerFloat32("Nova Altura", prompt)
		case "Peso":
			p.Peso, prompt = utils.LerFloat32("Novo Peso", prompt)
		default:
			fmt.Println("\nCampo invalido, digite novamente")
			utils.Pause()
		}
		continuar, _ = utils.LerBool("\nDeseja alterar mais algum campo? S/N", prompt)
	}
	p.CalculateSize()
}
