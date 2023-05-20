package lzw

import (
	// "fmt"
	"math"
	"os"
    "encoding/binary"

	"github.com/Bernardo46-2/AEDS-III/utils"
)

// ================================================ consts ================================================ //

// Maior tamanho possivel do dicionario definido
// pelo metodo de compressão LZW
const DICT_MAX_SIZE uint16 = 4096

// ================================================ ZipDict ================================================ //

// ZipDict é o dicionario usado para compactar arquivos usando a tecnica
// de compressão LZW. Possui um mapa de string e uint16 e uma variavel
// contando o numero de elementos inseridos no mapa
type ZipDict struct {
    dict map[string]uint16
    size uint16
}

// initZipDict inicializa um dicionario de compactação com seus
// valores padrões
func initZipDict() ZipDict {
    return ZipDict{
        dict: make(map[string]uint16, DICT_MAX_SIZE),
        size: 0,
    }
}

// isFull testa se o dicionario está cheio
func (d *ZipDict) isFull() bool {
    return d.size == DICT_MAX_SIZE - 1
}

// push insere um elemento no dicionario
func (d *ZipDict) push(s string) {
    if !d.isFull() {
        d.size++
        d.dict[s] = d.size
    }
}

// get pesquisa um elemento no dicionario, retornando seu possivel
// valor e um booleano indicando se o valor foi encontrado ou nao
// Caso o valor nao for encontrado, o uint16 retornado sera o valor
// padrao: 0
func (d *ZipDict) get(s string) (uint16, bool) {
    bytes := []byte(s)

	if len(bytes) == 1 {
		value, _ := utils.BytesToUint16(bytes, 0)
		return value, true
	} else {
		value, contains := d.dict[s]
		return value + math.MaxUint8, contains
	}
}

// ================================================ Bit Compress ================================================ //

// compress12bitArray percorre um []uint16 que possui numeros de exclusivamente
// 12 bits e os concatena com os proximos bits
// 
// in:  0000 aaaa aaaa aaaa  0000 bbbb bbbb bbbb  0000 cccc cccc cccc  0000 dddd dddd dddd
// 
// out: aaaa aaaa aaaa bbbb  bbbb bbbb cccc cccc  cccc dddd dddd dddd
func compress12bitArray(bytes []uint16) []uint16 {
    newLen := len(bytes) * 8 / 12 + 1
    if len(bytes) * 8 % 12 != 0 {
        newLen++
    }
    result := make([]uint16, newLen)
    index := 0

    for i := 0; i < len(bytes) - 1 && index < newLen; i++ {
        if i % 4 == 0 || i % 4 == 3 {
            result[index] = bytes[i] << 4
            result[index] |= bytes[i+1] >> 8
        } else if i % 4 == 1 {
            result[index] = bytes[i] << 8
            result[index] |= bytes[i+1] >> 4
        } else if i % 4 == 2 {
            result[index] = bytes[i] << 12
            result[index] |= bytes[i+1]
            i++
        }

        index++
    }

    return result
}

// ================================================ LZW ================================================ //

// writeFile escreve todos os bytes de um array de uint16 em
// um arquivo
func writeFile(f *os.File, content []uint16) error {
    for _, b := range content {
        err := binary.Write(f, binary.LittleEndian, b)
        if err != nil {
            return err
        }
    }
    return nil
}

// parseValue busca o maior valor possivel no dicionario enviado
// por parametro, eventualmente inserindo um novo valor no dicionario
// caso ele nao esteja cheio
func parseValue(dict *ZipDict, content []byte) (uint16, int) {
    offset := 0
    var value uint16

	for j := 1; j <= len(content); j++ {
		tmp, contains := dict.get(string(content[:j]))

		if contains {
			offset = len(content[:j])
			value = tmp
		} else {
			dict.push(string(content[:j]))
			break
		}
	}

	return value, offset
}

// zip recebe um dicionario e um conteudo em bytes e percorre
// todo o conteudo pesquisando cada valor encontrado e inserindo
// novos valores no caminho
func zip(dict ZipDict, content []byte) []uint16 {
    zipped := make([]uint16, 0, len(content))

	for i := 0; i < len(content); i++ {
		value, offset := parseValue(&dict, content[i:])
		i += offset - 1
		zipped = append(zipped, value)
	}

	return zipped
}

// Zip recebe uma string contendo o caminho do arquivo a ser compactado
// e o compacta usando o metodo de compressao LZW
func Zip(path string) error {
    path = "data/files/database/pokedex.bin"

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	dict := initZipDict()
    zipped := zip(dict, content)
    zipped = compress12bitArray(zipped)

    f, _ := os.Create("data/files/database/pokedex.lzw")
    defer f.Close()
    return writeFile(f, zipped)
}
