package rabinKarp

import (
	"fmt"
	"math"
)

const (
	// Inteiro usado na função de hash
	D int = 257

	// Número primo usado na hash para garantir que os valores
	// não ficarão grandes demais
	Q int = 997
)

// hash recebe uma string e calcula um valor que será
// usado durante o casamento de padrões
func hash(str string) (h int) {
	for _, c := range str {
		h = (h*D + int(c)) % Q
	}
	return
}

func RabinKarp(pattern string, text string) (indexes []int) {
	pl := len(pattern) // Tamanho do padrão
	tl := len(text)    // Tamanho do texto

	if tl < pl || tl == 0 || pl == 0 {
		return nil
	}

	ph := hash(pattern)   // Hash do padrão
	th := hash(text[:pl]) // Hash do primeiro slice do texto

	// Valor de hash usado para reajustar a hash do slice a cada iteração
	h := int(math.Pow(float64(D), float64(pl)-1)) % Q

	for i := 0; i <= tl-pl; i++ {
		// Se hash do texto e do slice são iguais, comparar
		// os caracteres
		if ph == th {
			contains := true

			// Iterando pelo slice testando se o padrão e slice
			// são iguais
			for j := 0; j < pl && contains; j++ {
				contains = text[i+j] == pattern[j]
			}

			// Se `contains` continuar como true até esse ponto,
			// significa que o padrão foi encontrado
			if contains {
				indexes = append(indexes, i)
			}
		}

		// Testando se não é a ultima iteração
		if i != tl-pl {
			// Recalculando Hash para próxima posição no texto
			th = (D*(th-int(text[i])*h) + int(text[i+pl])) % Q

			// Tratando caso o valor da hash ficar negativo
			if th < 0 {
				th += Q
			}
		}
	}

	return
}

func Test() {
	pattern := "aaba"
	text := "aabaabbaabaaba"
	is := RabinKarp(pattern, text)
	fmt.Println(is)
}
