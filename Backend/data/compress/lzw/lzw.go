package lzw

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/Bernardo46-2/AEDS-III/utils"
)

// ================================================ consts ================================================ //

// Maior tamanho possivel do dicionario definido
// pelo metodo de compressão LZW
const DICT_MAX_SIZE uint16 = 4096

// ================================================ Dict ================================================ //

// Dict é o dicionario usado para compactar arquivos usando a tecnica
// de compressão LZW. Possui um mapa de string e uint16 e uma variavel
// contando o numero de elementos inseridos no mapa
type Dict struct {
	dict      map[string]uint16
	unzipDict map[uint16]string
	sizeZip   uint16
	sizeUnzip uint16
}

// initZipDict inicializa um dicionario de compactação com seus
// valores padrões
func initZipDict() Dict {
	return Dict{
		dict:      make(map[string]uint16, DICT_MAX_SIZE),
		unzipDict: make(map[uint16]string, DICT_MAX_SIZE),
		sizeZip:   0,
		sizeUnzip: 0,
	}
}

// isFullZip testa se o dicionario está cheio
func (d *Dict) isFullZip() bool {
	return d.sizeZip == DICT_MAX_SIZE-1
}

// isFullUnzip testa se o dicionario está cheio
func (d *Dict) isFullUnzip() bool {
	return d.sizeZip == DICT_MAX_SIZE-1
}

// pushZip insere um elemento no dicionario
func (d *Dict) pushZip(s string) {
	if !d.isFullZip() {
		d.sizeZip++
		d.dict[s] = d.sizeZip
	}
}

// pushUnzip insere um elemento no dicionario
func (d *Dict) pushUnzip(s []byte) {
	if !d.isFullUnzip() {
		d.sizeUnzip++
		d.unzipDict[d.sizeUnzip-1] = string(s)
	}
}

// get pesquisa um elemento no dicionario, retornando seu possivel
// valor e um booleano indicando se o valor foi encontrado ou nao
// Caso o valor nao for encontrado, o uint16 retornado sera o valor
// padrao: 0
func (d *Dict) get(s string) (uint16, bool) {
	bytes := []byte(s)

	if len(bytes) == 1 {
		value, _ := utils.BytesToUint16(bytes, 0)
		return value, true
	} else {
		value, contains := d.dict[s]
		return value + math.MaxUint8, contains
	}
}

// getUnzip pesquisa um elemento no dicionario, retornando seu possivel
// valor e um booleano indicando se o valor foi encontrado ou nao
// Se o valor nao for encontrado, o retorno é nil
func (d *Dict) getUnzip(n uint16) []byte {
	if n < math.MaxUint8 {
		return []byte{byte(n)}
	} else {
		value, contains := d.unzipDict[n-math.MaxUint8-1]

		if contains {
			return []byte(value)
		}
	}

	return nil
}

// ================================================ Bit Compress ================================================ //

// compress12bitArray percorre um []uint16 que possui numeros de exclusivamente
// 12 bits e os concatena com os proximos bits
//
// in:  0000 aaaa aaaa aaaa | 0000 bbbb bbbb bbbb | 0000 cccc cccc cccc | 0000 dddd dddd dddd
//
// out: aaaa aaaa aaaa bbbb | bbbb bbbb cccc cccc | cccc dddd dddd dddd
func compress12bitArray(array []uint16) []uint16 {
	newLen := int(math.Ceil(float64(len(array)) * 3 / 4))
	result := make([]uint16, newLen)
	index := 0
	i := 0
	c := 0

	for i = 0; i < len(array)-1 && index < newLen; i++ {
		if c%3 == 0 {
			result[index] = array[i] << 4
			result[index] |= (array[i+1] >> 8) & 0x000F
			c++
		} else if c%3 == 1 {
			result[index] = array[i] << 8
			result[index] |= (array[i+1] >> 4) & 0x00FF
			c++
		} else {
			result[index] = array[i] << 12
			result[index] |= array[i+1] & 0x0FFF
			c = 0
			i++
		}

		index++
	}

	if i == len(array)-1 && index == newLen-1 {
		if c%3 == 0 {
			result[index] = array[i] << 4
		} else if c%3 == 1 {
			result[index] = array[i] << 8
		} else if c%3 == 2 {
			result[index] = array[i] << 12
		}
	}

	return result
}

// decompress12bitArray pega um array de uin16, consistindo de numeros de
// 12 bits compactados e os descompacta retornando um novo array de uint16
//
// in:  aaaa aaaa aaaa bbbb | bbbb bbbb cccc cccc | cccc dddd dddd dddd
//
// out: 0000 aaaa aaaa aaaa | 0000 bbbb bbbb bbbb | 0000 cccc cccc cccc | 0000 dddd dddd dddd
func decompress12bitArray(array []uint16) []uint16 {
	newLen := int(math.Ceil(float64(len(array))/3*4)) - 1
	result := make([]uint16, newLen)
	index := 0
	i := 0

	for i = 0; i < newLen-1 && index < len(array)-1; i += 2 {
		if i%4 == 0 {
			result[i] = array[index] >> 4
			result[i+1] = array[index]<<8&0x0FFF | array[index+1]>>8
		} else if i%4 == 2 {
			result[i] = array[index]<<4&0x0FFF | array[index+1]>>12
			result[i+1] = array[index+1] & 0x0FFF
			index++
		}

		index++
	}

	fmt.Println(index, len(array))

	if index == len(array)-1 {
		result[i] = array[index] >> 4
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

// unzip recebe um dicionario e um conteudo em bytes e percorre
// todo o conteudo pesquisando cada valor encontrado e inserindo
// novos valores no caminho
func unzip(dict Dict, content []uint16) []byte {
	unzipped := make([]byte, 0, len(content)*2)

	for i := 0; i < len(content)-1; i++ {
		value := dict.getUnzip(content[i])
		next := dict.getUnzip(content[i+1])
		unzipped = append(unzipped, value...)
		if len(next) > 0 {
			dict.pushUnzip(append(value, next[0]))
		}
	}

	value := dict.getUnzip(content[len(content)-1])
	unzipped = append(unzipped, value...)

	return unzipped
}

// parseValueZip busca o maior valor possivel no dicionario enviado
// por parametro, eventualmente inserindo um novo valor no dicionario
// caso ele nao esteja cheio
func parseValueZip(dict *Dict, content []byte) (uint16, int) {
	offset := 1
	var value uint16

	for i := 1; i <= len(content); i++ {
		tmp, contains := dict.get(string(content[:i]))

		if contains {
			offset = len(content[:i])
			value = tmp
		} else {
			dict.pushZip(string(content[:i]))
			break
		}
	}

	return value, offset
}

// zip recebe um dicionario e um conteudo em bytes e percorre
// todo o conteudo pesquisando cada valor encontrado e inserindo
// novos valores no caminho
func zip(dict Dict, content []byte) []uint16 {
	zipped := make([]uint16, 0, int(math.Ceil(float64(len(content))/2)))
	i := 0

	for i < len(content) {
		value, offset := parseValueZip(&dict, content[i:])
		zipped = append(zipped, value)
		i += offset
	}

	return zipped
}

// bytesToUint16s converte um array de bytes para um array de
// uint16
func bytesToUint16s(bytes []byte) []uint16 {
	length := len(bytes) / 2
	result := make([]uint16, length)
	j := 0

	for i := 0; i < length; i++ {
		result[i] = uint16(bytes[j]) | uint16(bytes[j+1])<<8
		j += 2
	}

	return result
}

// Zip recebe uma string contendo o caminho do arquivo a ser compactado
// e o compacta usando o metodo de compressao LZW
func Zip(path string) error {
	extension := filepath.Ext(path)

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	zipped := zip(initZipDict(), content)
	// zipped = compress12bitArray(zipped)

	f, _ := os.Create(utils.ChangeExtension(path, ".lzw"))
	writeFile(f, zipped)
	f.Close()

	content, _ = os.ReadFile(utils.ChangeExtension(path, ".lzw"))
	unzipped := unzip(initZipDict(), bytesToUint16s(content))
	// unzipped := decompress12bitArray(zipped)

	f, _ = os.Create(utils.ChangeExtension(path, "2"+extension))
	f.Write(unzipped)
	f.Close()

	return nil
}
