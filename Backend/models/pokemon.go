// O arquivo Pokemon do pacote Models permite a criação de uma estrutura pokemon
// tao como a manipulação de seus dados em diferentes niveis e maneiras.
// Seja para binario, json, string ou o contrario
package models

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/Bernardo46-2/AEDS-III/utils"
)

const MAX_NAME_LEN = 40

// PokemonID faz a formatação do ID para serialização em JSON
type PokemonID struct {
	ID int `json:"id"`
}

// Pokemon representa um Pokémon e seus atributos, como número, nome,
// espécie, habilidades e características físicas.
//
// Todos os dados sao serializaveis com exceção de Size
type Pokemon struct {
	Numero     int32     `json:"numero"`
	Nome       string    `json:"nome,omitempty"`
	NomeJap    string    `json:"nomeJap,omitempty"`
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
	Descricao  string    `json:"descricao"`
	Size       PokeSize  `json:"-"`
}

// PokeSize faz a intermediação para geração de um array de bytes ligado ao
// tamanho de cada variavel para armazenamento em arquivo binario
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
	Descricao  int32
}

// GenReleaseDates é um mapa para facil conversão de geração em data de lançamento
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

// ToString faz o parsing da struct em formato legivel para debug
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
	str += fmt.Sprintf("Descricao  = %s\n", p.Descricao)

	return str
}

// copyBytes é uma função que copia os bytes do slice de origem (src)
// para o slice de destino (dest) a partir do deslocamento especificado (offset).
// Retorna o slice de destino atualizado e o novo deslocamento atualizado.
func copyBytes(dest []byte, src []byte, offset int) ([]byte, int) {
	copy(dest[offset:], src)
	return dest, offset + len(src)
}

// ToBytes realiza a serialização em array binario da struct Pokemon
//
// O primeiro valor é o tamanho do registro
// O tamanho dos valores variaveis como strings sao armazenados antes do valor em si
// O padrão utilizado é int32 para otimizar espaço
func (p *Pokemon) ToBytes() []byte {
	// Inicializa dados
	pokeBytes := make([]byte, p.Size.Total+4)
	var lendario, mitico []byte
	offset := 0
	lapide := 0

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

	// serialização da data
	releaseDate, _ := p.Lancamento.MarshalBinary()
	filler := make([]byte, p.Size.Nome-int32(len(p.Nome)))
	runes := []rune(p.NomeJap)
	japName := make([]byte, len(runes)*4)

	for i, v := range runes {
		binary.LittleEndian.PutUint32(japName[i*4:(i+1)*4], uint32(v))
	}

	// Longo e chato processo de conversao de tamanho da variavel + variavel
	pokeBytes, offset = copyBytes(pokeBytes, utils.IntToBytes(int32(lapide)), offset)
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
	pokeBytes, offset = copyBytes(pokeBytes, utils.FloatToBytes(p.Peso), offset)
	pokeBytes, offset = copyBytes(pokeBytes, utils.IntToBytes(p.Size.Descricao), offset)
	pokeBytes, _ = copyBytes(pokeBytes, []byte(p.Descricao), offset)

	return pokeBytes
}

// ParseBinToPoke faz a desserialização do arquivo binario e retorna um
// struct do tipo Pokemon
func (p *Pokemon) ParseBinToPoke(registro []byte) error {
	// ptr serve para andar pelo registro de maneira incremental
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
	p.Peso, ptr = utils.BytesToFloat32(registro, ptr)
	p.Descricao, _ = utils.BytesToString(registro, ptr)
	p.CalculateSize()

	return nil
}

// ParsePokemon recebe um array de strings vindo do CSV e faz a conversao
// para a struct
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
	pokemon.Descricao = line[len(line)-1]

	pokemon.CalculateSize()

	return pokemon
}

// CalculateSize adiciona ao campo nao serializavel SIZE da struct Pokemon
// o somatorio do tamanho em byte de todos os campos + features necessarias
// para a serialização em binario
func (p *Pokemon) CalculateSize() {
	// Calcula os tamanhos
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
	p.Size.Descricao = int32(len(p.Descricao))

	// Soma e adiciona o espaço ocupado pelo bit de tamanho
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
		p.Size.Peso +
		p.Size.Descricao + 4 + 1
}

func (p Pokemon) GetField(fieldName string) string {
	field := strings.ToLower(fieldName)
	switch field {
	case "numero", "id":
		return fmt.Sprint(p.Numero)
	case "nome":
		return p.Nome
	case "nomejap":
		return p.NomeJap
	case "geracao":
		return fmt.Sprint(p.Geracao)
	case "lancamento":
		return p.Lancamento.Format(time.RFC3339)
	case "especie":
		return p.Especie
	case "lendario":
		return fmt.Sprint(p.Lendario)
	case "mitico":
		return fmt.Sprint(p.Mitico)
	case "tipo":
		return strings.Join(p.Tipo, ",")
	case "atk":
		return fmt.Sprint(p.Atk)
	case "def":
		return fmt.Sprint(p.Def)
	case "hp":
		return fmt.Sprint(p.Hp)
	case "altura":
		return fmt.Sprint(p.Altura)
	case "peso":
		return fmt.Sprint(p.Peso)
	case "descricao":
		return p.Descricao
	default:
		return ""
	}
}

func PokeStrings() []string {
	obj := Pokemon{}
	var fieldNames []string
	v := reflect.ValueOf(obj)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.String || (field.Kind() == reflect.Slice && field.Type().Elem().Kind() == reflect.String) {
			fieldNames = append(fieldNames, utils.Decaptalize(t.Field(i).Name))
		}
	}

	return fieldNames
}

func (p Pokemon) GetFieldF64(fieldName string) float64 {
	field := strings.ToLower(fieldName)
	switch field {
	case "numero", "id":
		return float64(p.Numero)
	case "geracao":
		return float64(p.Geracao)
	case "lancamento":
		return float64(p.Lancamento.Unix())
	case "atk":
		return float64(p.Atk)
	case "def":
		return float64(p.Def)
	case "hp":
		return float64(p.Hp)
	case "altura":
		return float64(p.Altura)
	case "peso":
		return float64(p.Peso)
	default:
		return -1
	}
}

func PokeNumbers() []string {
	fields := []string{}

	valueOf := reflect.ValueOf(Pokemon{})
	typeOf := valueOf.Type()

	for i := 0; i < valueOf.NumField(); i++ {
		// field := valueOf.Field(i)
		fieldType := typeOf.Field(i)

		if fieldType.Type.Kind() == reflect.Int || fieldType.Type.Kind() == reflect.Int8 ||
			fieldType.Type.Kind() == reflect.Int16 || fieldType.Type.Kind() == reflect.Int32 ||
			fieldType.Type.Kind() == reflect.Int64 {
			fields = append(fields, utils.Decaptalize(fieldType.Name))
		}

		if fieldType.Type.Kind() == reflect.Float32 || fieldType.Type.Kind() == reflect.Float64 {
			fields = append(fields, utils.Decaptalize(fieldType.Name))
		}

		if fieldType.Type.Kind() == reflect.Slice {
			sliceElemType := fieldType.Type.Elem()
			if sliceElemType.Kind() == reflect.Int || sliceElemType.Kind() == reflect.Int8 ||
				sliceElemType.Kind() == reflect.Int16 || sliceElemType.Kind() == reflect.Int32 ||
				sliceElemType.Kind() == reflect.Int64 {
				fields = append(fields, utils.Decaptalize(fieldType.Name))
			}

			if sliceElemType.Kind() == reflect.Float32 || sliceElemType.Kind() == reflect.Float64 {
				fields = append(fields, utils.Decaptalize(fieldType.Name))
			}
		}
	}

	return fields
}
