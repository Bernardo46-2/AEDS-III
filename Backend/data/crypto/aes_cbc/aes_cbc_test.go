// Testes de encriptação e descriptação com exemplos do artigo oficial
// da NIST sobre AES
//
// Obs: Testes feitos abaixo desconsideram o iv e o padding quando comparando
// os bytes, apenas porque o artigo também os ignora
//
// Artigo (exemplos do AES CBC começam na página 27):
// https://nvlpubs.nist.gov/nistpubs/Legacy/SP/nistspecialpublication800-38a.pdf
package aescbc

import "testing"

var IV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
var PLAIN_TEXT = []byte{0x6b, 0xc1, 0xbe, 0xe2, 0x2e, 0x40, 0x9f, 0x96, 0xe9, 0x3d, 0x7e, 0x11, 0x73, 0x93, 0x17, 0x2a,
	0xae, 0x2d, 0x8a, 0x57, 0x1e, 0x03, 0xac, 0x9c, 0x9e, 0xb7, 0x6f, 0xac, 0x45, 0xaf, 0x8e, 0x51,
	0x30, 0xc8, 0x1c, 0x46, 0xa3, 0x5c, 0xe4, 0x11, 0xe5, 0xfb, 0xc1, 0x19, 0x1a, 0x0a, 0x52, 0xef,
	0xf6, 0x9f, 0x24, 0x45, 0xdf, 0x4f, 0x9b, 0x17, 0xad, 0x2b, 0x41, 0x7b, 0xe6, 0x6c, 0x37, 0x10}

func TestKeyExpansion(t *testing.T) {
	key, _ := NewKeyFrom([]byte{0x2b, 0x7e, 0x15, 0x16, 0x28, 0xae, 0xd2, 0xa6, 0xab, 0xf7, 0x15, 0x88, 0x09, 0xcf, 0x4f, 0x3c})

	expectedOutput := []byte{0x2b, 0x7e, 0x15, 0x16, 0x28, 0xae, 0xd2, 0xa6, 0xab, 0xf7, 0x15, 0x88, 0x09, 0xcf, 0x4f, 0x3c,
		0xa0, 0xfa, 0xfe, 0x17, 0x88, 0x54, 0x2c, 0xb1, 0x23, 0xa3, 0x39, 0x39, 0x2a, 0x6c, 0x76, 0x05,
		0xf2, 0xc2, 0x95, 0xf2, 0x7a, 0x96, 0xb9, 0x43, 0x59, 0x35, 0x80, 0x7a, 0x73, 0x59, 0xf6, 0x7f,
		0x3d, 0x80, 0x47, 0x7d, 0x47, 0x16, 0xfe, 0x3e, 0x1e, 0x23, 0x7e, 0x44, 0x6d, 0x7a, 0x88, 0x3b,
		0xef, 0x44, 0xa5, 0x41, 0xa8, 0x52, 0x5b, 0x7f, 0xb6, 0x71, 0x25, 0x3b, 0xdb, 0x0b, 0xad, 0x00,
		0xd4, 0xd1, 0xc6, 0xf8, 0x7c, 0x83, 0x9d, 0x87, 0xca, 0xf2, 0xb8, 0xbc, 0x11, 0xf9, 0x15, 0xbc,
		0x6d, 0x88, 0xa3, 0x7a, 0x11, 0x0b, 0x3e, 0xfd, 0xdb, 0xf9, 0x86, 0x41, 0xca, 0x00, 0x93, 0xfd,
		0x4e, 0x54, 0xf7, 0x0e, 0x5f, 0x5f, 0xc9, 0xf3, 0x84, 0xa6, 0x4f, 0xb2, 0x4e, 0xa6, 0xdc, 0x4f,
		0xea, 0xd2, 0x73, 0x21, 0xb5, 0x8d, 0xba, 0xd2, 0x31, 0x2b, 0xf5, 0x60, 0x7f, 0x8d, 0x29, 0x2f,
		0xac, 0x77, 0x66, 0xf3, 0x19, 0xfa, 0xdc, 0x21, 0x28, 0xd1, 0x29, 0x41, 0x57, 0x5c, 0x00, 0x6e,
		0xd0, 0x14, 0xf9, 0xa8, 0xc9, 0xee, 0x25, 0x89, 0xe1, 0x3f, 0x0c, 0xc8, 0xb6, 0x63, 0x0c, 0xa6}

	if len(expectedOutput) != len(key.Expanded) {
		t.Errorf("Key Expansion error:\nkey_len = %d\nexpected_len = %d\n\n", len(key.Expanded), len(expectedOutput))
	}

	for i := 0; i < len(expectedOutput); i++ {
		if expectedOutput[i] != key.Expanded[i] {
			t.Errorf("Key Expansion error:\nkey[%d] = %02x\nexpected[%d] = %02x\n\n", i, key.Expanded[i], i, expectedOutput[i])
		}
	}
}

func TestAes128Cbc(t *testing.T) {
	key, _ := NewKeyFrom([]byte{0x2b, 0x7e, 0x15, 0x16, 0x28, 0xae, 0xd2, 0xa6, 0xab, 0xf7, 0x15, 0x88, 0x09, 0xcf, 0x4f, 0x3c})
	expectedOutput := []byte{0x76, 0x49, 0xab, 0xac, 0x81, 0x19, 0xb2, 0x46, 0xce, 0xe9, 0x8e, 0x9b, 0x12, 0xe9, 0x19, 0x7d,
		0x50, 0x86, 0xcb, 0x9b, 0x50, 0x72, 0x19, 0xee, 0x95, 0xdb, 0x11, 0x3a, 0x91, 0x76, 0x78, 0xb2,
		0x73, 0xbe, 0xd6, 0xb8, 0xe3, 0xc1, 0x74, 0x3b, 0x71, 0x16, 0xe6, 0x9e, 0x22, 0x22, 0x95, 0x16,
		0x3f, 0xf1, 0xca, 0xa1, 0x68, 0x1f, 0xac, 0x09, 0x12, 0x0e, 0xca, 0x30, 0x75, 0x86, 0xe1, 0xa7}

	cipherText := Encrypt(key, IV, PLAIN_TEXT)

	for i := BLOCK_SIZE; i < len(cipherText)-BLOCK_SIZE; i++ {
		if cipherText[i] != expectedOutput[i-BLOCK_SIZE] {
			t.Errorf("Encryption error:\nexpectedOutput[%d] = %02x\ncipherText[%d] = %02x\n\n", i-BLOCK_SIZE, expectedOutput[i-BLOCK_SIZE], i-BLOCK_SIZE, cipherText[i])
		}
	}

	decryptedText := Decrypt(key, cipherText)

	for i := 0; i < len(PLAIN_TEXT); i++ {
		if PLAIN_TEXT[i] != decryptedText[i] {
			t.Errorf("Encryption error:\nplainText[%d] = %02x\ndecryptedText[%d] = %02x\n\n", i, PLAIN_TEXT[i], i, decryptedText[i])
		}
	}
}

func TestAes192Cbc(t *testing.T) {
	key, _ := NewKeyFrom([]byte{0x8e, 0x73, 0xb0, 0xf7, 0xda, 0x0e, 0x64, 0x52, 0xc8, 0x10, 0xf3, 0x2b, 0x80, 0x90, 0x79, 0xe5,
		0x62, 0xf8, 0xea, 0xd2, 0x52, 0x2c, 0x6b, 0x7b})
	expectedOutput := []byte{0x4f, 0x02, 0x1d, 0xb2, 0x43, 0xbc, 0x63, 0x3d, 0x71, 0x78, 0x18, 0x3a, 0x9f, 0xa0, 0x71, 0xe8,
		0xb4, 0xd9, 0xad, 0xa9, 0xad, 0x7d, 0xed, 0xf4, 0xe5, 0xe7, 0x38, 0x76, 0x3f, 0x69, 0x14, 0x5a,
		0x57, 0x1b, 0x24, 0x20, 0x12, 0xfb, 0x7a, 0xe0, 0x7f, 0xa9, 0xba, 0xac, 0x3d, 0xf1, 0x02, 0xe0,
		0x08, 0xb0, 0xe2, 0x79, 0x88, 0x59, 0x88, 0x81, 0xd9, 0x20, 0xa9, 0xe6, 0x4f, 0x56, 0x15, 0xcd}

	cipherText := Encrypt(key, IV, PLAIN_TEXT)

	for i := BLOCK_SIZE; i < len(cipherText)-BLOCK_SIZE; i++ {
		if cipherText[i] != expectedOutput[i-BLOCK_SIZE] {
			t.Errorf("Encryption error:\nexpectedOutput[%d] = %02x\ncipherText[%d] = %02x\n\n", i-BLOCK_SIZE, expectedOutput[i-BLOCK_SIZE], i-BLOCK_SIZE, cipherText[i])
		}
	}

	decryptedText := Decrypt(key, cipherText)

	for i := 0; i < len(PLAIN_TEXT); i++ {
		if PLAIN_TEXT[i] != decryptedText[i] {
			t.Errorf("Encryption error:\nplainText[%d] = %02x\ndecryptedText[%d] = %02x\n\n", i, PLAIN_TEXT[i], i, decryptedText[i])
		}
	}
}

func TestAes256Cbc(t *testing.T) {
	key, _ := NewKeyFrom([]byte{0x60, 0x3d, 0xeb, 0x10, 0x15, 0xca, 0x71, 0xbe, 0x2b, 0x73, 0xae, 0xf0, 0x85, 0x7d, 0x77, 0x81,
		0x1f, 0x35, 0x2c, 0x07, 0x3b, 0x61, 0x08, 0xd7, 0x2d, 0x98, 0x10, 0xa3, 0x09, 0x14, 0xdf, 0xf4})
	expectedOutput := []byte{0xf5, 0x8c, 0x4c, 0x04, 0xd6, 0xe5, 0xf1, 0xba, 0x77, 0x9e, 0xab, 0xfb, 0x5f, 0x7b, 0xfb, 0xd6,
		0x9c, 0xfc, 0x4e, 0x96, 0x7e, 0xdb, 0x80, 0x8d, 0x67, 0x9f, 0x77, 0x7b, 0xc6, 0x70, 0x2c, 0x7d,
		0x39, 0xf2, 0x33, 0x69, 0xa9, 0xd9, 0xba, 0xcf, 0xa5, 0x30, 0xe2, 0x63, 0x04, 0x23, 0x14, 0x61,
		0xb2, 0xeb, 0x05, 0xe2, 0xc3, 0x9b, 0xe9, 0xfc, 0xda, 0x6c, 0x19, 0x07, 0x8c, 0x6a, 0x9d, 0x1b}

	cipherText := Encrypt(key, IV, PLAIN_TEXT)

	for i := BLOCK_SIZE; i < len(cipherText)-BLOCK_SIZE; i++ {
		if cipherText[i] != expectedOutput[i-BLOCK_SIZE] {
			t.Errorf("Encryption error:\nexpectedOutput[%d] = %02x\ncipherText[%d] = %02x\n\n", i-BLOCK_SIZE, expectedOutput[i-BLOCK_SIZE], i-BLOCK_SIZE, cipherText[i])
		}
	}

	decryptedText := Decrypt(key, cipherText)

	for i := 0; i < len(PLAIN_TEXT); i++ {
		if PLAIN_TEXT[i] != decryptedText[i] {
			t.Errorf("Encryption error:\nplainText[%d] = %02x\ndecryptedText[%d] = %02x\n\n", i, PLAIN_TEXT[i], i, decryptedText[i])
		}
	}
}
