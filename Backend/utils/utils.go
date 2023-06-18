// O pacote utils contém uma série de funções utilitárias para ajudar no
// desenvolvimento do software. As funções incluem conversão de tipos de
// dados, manipulação de bytes e strings, entre outras.
// Essas funções foram criadas para melhorar a eficiência e legibilidade
// do código, além de fornecer soluções para problemas comuns encontrados
// no desenvolvimento.
package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Atoi32 converte uma string em um int32 e retorna o resultado
// e um erro, se houver algum problema na conversão.
func Atoi32(s string) (int32, error) {
	i, err := strconv.Atoi(s)
	return int32(i), err
}

// UintToBytes converte um número uint16 em uma slice de bytes
// usando a ordem Little Endian e retorna a slice resultante.
func Uint16ToBytes(n uint16) []byte {
	var buf []byte
	return binary.LittleEndian.AppendUint16(buf, uint16(n))
}

// IntToBytes converte um número int32 em uma slice de bytes
// usando a ordem Little Endian e retorna a slice resultante.
func IntToBytes(n int32) []byte {
	var buf []byte
	return binary.LittleEndian.AppendUint32(buf, uint32(n))
}

// IntToBytes converte um número int64 em uma slice de bytes
// usando a ordem Little Endian e retorna a slice resultante.
func Int64ToBytes(n int64) []byte {
	var buf []byte
	return binary.LittleEndian.AppendUint64(buf, uint64(n))
}

// FloatToBytes converte um número float32 em uma slice de bytes
// usando a ordem Little Endian e retorna a slice resultante.
func FloatToBytes(f float32) []byte {
	b := make([]byte, 4)
	bits := math.Float32bits(f)
	binary.LittleEndian.PutUint32(b, bits)

	return b
}

// BoolToByte converte um booleano em um byte
// usando a ordem Little Endian e retorna um byte 0 ou 1.
func BoolToByte(b bool) byte {
	if b {
		return byte(1)
	} else {
		return byte(0)
	}
}

// RemoveAfterSpace remove tudo depois do primeiro espaço em branco
// encontrado na string e retorna o resultado.
func RemoveAfterSpace(str string) string {
	parts := strings.Split(str, " ")
	return parts[0]
}

// BytesToVarSize retorna o tamanho de um campo de tamanho variável e avança o ponteiro
func BytesToVarSize(registro []byte, ptr int) (int, int) {
	return int(binary.LittleEndian.Uint32(registro[ptr : ptr+4])), ptr + 4
}

// BytesToUint16 retorna um inteiro unsigned de 16 bits e avança o ponteiro ptr
func BytesToUint16(registro []byte, ptr int) (uint16, int) {
	return binary.LittleEndian.Uint16(registro[ptr : ptr+4]), ptr + 4
}

// BytesToInt32 retorna um inteiro de 32 bits e avança o ponteiro ptr
func BytesToInt32(registro []byte, ptr int) (int32, int) {
	return int32(binary.LittleEndian.Uint32(registro[ptr : ptr+4])), ptr + 4
}

// BytesToInt64 retorna um inteiro de 64 bits e avança o ponteiro ptr
func BytesToInt64(registro []byte, ptr int) (int64, int) {
	return int64(binary.LittleEndian.Uint64(registro[ptr : ptr+8])), ptr + 8
}

// BytesToString retorna uma string de tamanho variável e avança o ponteiro ptr
func BytesToString(registro []byte, ptr int) (string, int) {
	size, ptr := BytesToVarSize(registro, ptr)
	nomeBytes := make([]byte, size)
	io.ReadFull(bytes.NewReader(registro[ptr:ptr+size]), nomeBytes)
	return strings.TrimSpace(string(nomeBytes)), ptr + size
}

// BytesToFixedSizeString retorna uma string de tamanho fixo e avança o ponteiro ptr
func BytesToFixedSizeString(registro []byte, ptr int, maxSize int) (string, int) {
	nome := make([]byte, maxSize)
	io.ReadFull(bytes.NewReader(registro[ptr:ptr+maxSize]), nome)
	return strings.TrimRight(string(nome), "\x00"), ptr + maxSize
}

// BytesToArrayString recebe um registro contendo strings tabuladas em virgula
// e retorna um array de strings e avança o ponteiro ptr
func BytesToArrayString(registro []byte, ptr int) ([]string, int) {
	size, ptr := BytesToVarSize(registro, ptr)
	stringBytes := make([]byte, size)
	io.ReadFull(bytes.NewReader(registro[ptr:ptr+size]), stringBytes)
	s := strings.TrimRight(string(stringBytes), ",")
	return strings.Split(s, ","), ptr + size
}

// BytesToJapName retorna uma string de tamanho variavel escrita em caracteres japoneses
// (4 bytes) e avança o ponteiro ptr
func BytesToJapName(registro []byte, ptr int) (string, int) {
	size, ptr := BytesToVarSize(registro, ptr)

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

// BytesToTime converte um slice de bytes contendo uma representação binária de um valor de tempo
// para um valor do tipo time.Time e move o ponteiro ptr
func BytesToTime(registro []byte, ptr int) (time.Time, int) {
	size, ptr := BytesToVarSize(registro, ptr)
	b := make([]byte, size)
	io.ReadFull(bytes.NewReader(registro[ptr:ptr+size]), b)
	var t time.Time
	t.UnmarshalBinary(b)
	return t, ptr + size
}

// BytesToBool converte um slice de bytes contendo uma representação binária de um valor booleano
// para um valor do tipo bool e move o ponteiro ptr
func BytesToBool(registro []byte, ptr int) (bool, int) {
	if registro[ptr] != 0 {
		return true, ptr + 1
	} else {
		return false, ptr + 1
	}
}

// BytesToFloat32 converte um slice de bytes contendo uma representação binária de um valor float32
// para um valor do tipo float32 e move o ponteiro ptr
func BytesToFloat32(registro []byte, ptr int) (float32, int) {
	size := 4
	bits := binary.LittleEndian.Uint32(registro[ptr : ptr+size])
	float := math.Float32frombits(bits)
	return float, ptr + size
}

// BytesToFloat64 converte um slice de bytes contendo uma representação binária de um valor float64
// para um valor do tipo float64 e move o ponteiro ptr
func BytesToFloat64(registro []byte, ptr int) (float64, int) {
	size := 8
	bits := binary.LittleEndian.Uint64(registro[ptr : ptr+size])
	float := math.Float64frombits(bits)
	return float, ptr + size
}

func AbsInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func InsertionSort(arr []int32) {
	for i := 1; i < len(arr); i++ {
		key := arr[i]
		j := i - 1
		for j >= 0 && arr[j] > key {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
	}
}

func Decaptalize(str string) string {
	if len(str) < 1 {
		return str
	}
	return strings.ToLower(str[0:1]) + str[1:]
}

func BoolToFloat(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}

func FormatDate(dateStr string) string {
	date, _ := time.Parse("02/01/2006", dateStr)

	unixTime := date.Unix()
	return strconv.FormatInt(unixTime, 10)
}

func FormatByte(b byte, size int) string {
	// Converte o byte em uma string binária
	binaryStr := fmt.Sprintf("%08b", b)

	// Verifica se o tamanho da string binária é maior que o tamanho desejado
	if len(binaryStr) > size {
		// Remove os bits excedentes do início da string binária
		binaryStr = binaryStr[len(binaryStr)-size:]
	} else if len(binaryStr) < size {
		// Adiciona zeros à esquerda para atingir o tamanho desejado
		padding := strings.Repeat("0", size-len(binaryStr))
		binaryStr = padding + binaryStr
	}

	return binaryStr
}

func ByteSize(v interface{}) int {
	var size int
	value := reflect.ValueOf(v)

	switch value.Kind() {
	case reflect.Int8, reflect.Uint8:
		size = 1
	case reflect.Int16, reflect.Uint16:
		size = 2
	case reflect.Int32, reflect.Uint32, reflect.Float32:
		size = 4
	case reflect.Int64, reflect.Uint64, reflect.Float64, reflect.Complex64:
		size = 8
	case reflect.Complex128:
		size = 16
	case reflect.String:
		size = int(len(v.(string)))
	case reflect.Slice:
		size = int(value.Len()) * ByteSize(value.Index(0).Interface())
	case reflect.Int, reflect.Uint:
		size = strconv.IntSize / 8
	default:
		fmt.Printf("Unhandled kind: %v\n", value.Kind())
		size = 0
	}

	return size
}

func RotateLeft(value uint8, n uint) uint8 {
	return (value << n) | (value >> (8 - n))
}

func RotateRight(value byte, n uint) byte {
	return (value >> n) | (value << (8 - n))
}

func ChangeExtension(filePath, newExtension string) string {
	// Obter a extensão atual do arquivo
	ext := filepath.Ext(filePath)

	// Remover a extensão atual do arquivo
	fileName := filePath[:len(filePath)-len(ext)]

	// Adicionar a nova extensão ao nome do arquivo
	newFilePath := fmt.Sprintf("%s%s", fileName, newExtension)

	return newFilePath
}

func ByteArrayToAscii(b [10]byte) string {
	str := fmt.Sprint(b)
	str = strings.Trim(str, "[]")
	return str
}

func StringToByteArray(s string) ([10]byte, error) {
	parts := strings.Split(s, " ")
	if len(parts) != 10 {
		return [10]byte{}, fmt.Errorf("string must be 10 bytes long")
	}

	var byteArray [10]byte
	for i, part := range parts {
		val, err := strconv.Atoi(part)
		if err != nil {
			return [10]byte{}, err
		}
		byteArray[i] = byte(val)
	}

	return byteArray, nil
}

const (
	VERIFIER string = "data/files/database/autenticity.txt"
	V_TEXT   string = "VERIFICADO!"
)

func Create_verifier() {

	// Cria os diretórios se eles não existirem
	os.MkdirAll("data/files/database", os.ModePerm)

	// Cria o arquivo
	file, _ := os.Create(VERIFIER)
	defer file.Close()

	// Escreve a mensagem no arquivo
	file.WriteString(V_TEXT)

}

func Verify(content string) bool {
	return content == V_TEXT
}
