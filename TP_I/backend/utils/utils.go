package utils

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
	"strconv"
	"strings"
	"time"
)

func Atoi32(s string) (int32, error) {
	i, err := strconv.Atoi(s)
	return int32(i), err
}

func IntToBytes(n int32) []byte {
	var buf []byte
	return binary.LittleEndian.AppendUint32(buf, uint32(n))
}

func FloatToBytes(f float32) []byte {
	b := make([]byte, 4)
	bits := math.Float32bits(f)
	binary.LittleEndian.PutUint32(b, bits)

	return b
}

func RemoveAfterSpace(str string) string {
	parts := strings.Split(str, " ")
	return parts[0]
}

func BytesToVarSize(registro []byte, ptr int) (int, int) {
	return int(binary.LittleEndian.Uint32(registro[ptr : ptr+4])), ptr + 4
}

func BytesToInt32(registro []byte, ptr int) (int32, int) {
	return int32(binary.LittleEndian.Uint32(registro[ptr : ptr+4])), ptr + 4
}

func BytesToString(registro []byte, ptr int) (string, int) {
	size, ptr := BytesToVarSize(registro, ptr)
	nomeBytes := make([]byte, size)
	io.ReadFull(bytes.NewReader(registro[ptr:ptr+size]), nomeBytes)
	return strings.TrimSpace(string(nomeBytes)), ptr + size
}

func BytesToFixedSizeString(registro []byte, ptr int, maxSize int) (string, int) {
	nome := make([]byte, maxSize)
	io.ReadFull(bytes.NewReader(registro[ptr:ptr+maxSize]), nome)
	return strings.TrimRight(string(nome), "\x00"), ptr + maxSize
}

func BytesToArrayString(registro []byte, ptr int) ([]string, int) {
	size, ptr := BytesToVarSize(registro, ptr)
	stringBytes := make([]byte, size)
	io.ReadFull(bytes.NewReader(registro[ptr:ptr+size]), stringBytes)
	s := strings.TrimRight(string(stringBytes), ",")
	return strings.Split(s, ","), ptr + size
}

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

func BytesToTime(registro []byte, ptr int) (time.Time, int) {
	size, ptr := BytesToVarSize(registro, ptr)
	b := make([]byte, size)
	io.ReadFull(bytes.NewReader(registro[ptr:ptr+size]), b)
	var t time.Time
	t.UnmarshalBinary(b)
	return t, ptr + size
}

func BytesToBool(registro []byte, ptr int) (bool, int) {
	if registro[ptr] != 0 {
		return true, ptr + 1
	} else {
		return false, ptr + 1
	}
}

func BytesToFloat32(registro []byte, ptr int) (float32, int) {
	size := 4
	bits := binary.LittleEndian.Uint32(registro[ptr : ptr+size])
	float := math.Float32frombits(bits)
	return float, ptr + size
}
