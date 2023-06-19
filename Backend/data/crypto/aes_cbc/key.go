package aescbc

import (
	"errors"
)

// Key é um struct que contém informações úteis sobrea a chave de
// criptografia aes, assim como a chave em si junto com sua versão
// expandida
type Key struct {
	Key       []byte
	Expanded  []byte
	SizeBits  int
	SizeBytes int
}

// NewKey inicializa uma chave com um array de bytes preenchido
// aleatóriamente com o tamanho especificado em bits, podendo ser
// 128, 192 ou 256 bits (16, 24 e 32 bytes, respectivamente)
func NewKey(sizeBits int) (key Key, err error) {
	if sizeBits != 128 && sizeBits != 192 && sizeBits != 256 {
		err = errors.New("invalid key size")
		return
	}

	sizeBytes := sizeBits / 8
	k, err := RandBytes(sizeBytes)
	if err != nil {
		return
	}

	expanded := expand(k)

	key = Key{
		SizeBits:  sizeBits,
		SizeBytes: sizeBytes,
		Key:       k,
		Expanded:  expanded,
	}

	return
}

// NewKeyFrom inicializa uma chave a partir de um array de bytes
// recebido por parâmetro. O tamanho do array será usado como tamanho
// da chave, podendo ser 32, 24 ou 32 bytes (128, 192 e 256 bits, respectivamente)
func NewKeyFrom(k []byte) (key Key, err error) {
	sizeBytes := len(k)
	sizeBits := sizeBytes * 8

	if sizeBits != 128 && sizeBits != 192 && sizeBits != 256 {
		err = errors.New("invalid key size")
		return
	}

	expanded := expand(k)

	key = Key{
		SizeBits:  sizeBits,
		SizeBytes: sizeBytes,
		Key:       k,
		Expanded:  expanded,
	}

	return
}

// expand faz a expansão da chave do aes, gerando mais blocos,
// chamados de round keys, para serem usados na criptografia.
// O número de round keys depende do tamanho da chave provida,
// sendo que a chave resultante sempre vai começar com os
// mesmos bytes da chave original, apenas então sucedidos pelos
// novos bytes gerados pela função
func expand(key []byte) (keySchedule []byte) {
	keyLen := len(key)
	var nr, nk int

	switch keyLen {
	case 16:
		nr, nk = 10, 4
	case 24:
		nr, nk = 12, 6
	case 32:
		nr, nk = 14, 8
	}

	keyScheduleLen := 16 * (nr + 1)
	keySchedule = make([]byte, 0, keyScheduleLen)
	keySchedule = append(keySchedule, key...)

	for i := 0; i < 4*nr+4-nk; i++ {
		word := make([]byte, 4)
		copy(word, keySchedule[i*4+keyLen-4:i*4+keyLen])

		if i%nk == 0 {
			RotateLeft(word)
			subBytes(word)
			word[0] ^= RCON[i/nk]
		} else if nk > 6 && i%nk == 4 {
			subBytes(word)
		}

		prevWord := keySchedule[i*4 : i*4+4]
		xorBlock(word, prevWord)
		keySchedule = append(keySchedule, word...)
	}

	return
}
