package rabinKarp

import (
	"math"
	"strings"

	"github.com/Bernardo46-2/AEDS-III/data/binManager"
	"github.com/Bernardo46-2/AEDS-III/data/indexes/invertedIndex"
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
    pattern = strings.ToLower(pattern)
    text = strings.ToLower(text)

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

func SearchPokemon(search string, field string) (scoredDocuments []invertedIndex.ScoredDocument) {
	controller, _ := binManager.InicializarControleLeitura(binManager.BIN_FILE)
	scoredDocuments = make([]invertedIndex.ScoredDocument, 0)

	for err := controller.ReadNext(); err == nil; err = controller.ReadNext() {
		if !controller.RegistroAtual.IsDead() {
			needle := RabinKarp(search, controller.RegistroAtual.Pokemon.GetField(field))
			if len(needle) > 0 {
				scoredDocuments = append(scoredDocuments, invertedIndex.ScoredDocument{DocumentID: int64(controller.RegistroAtual.Pokemon.Numero), Score: len(needle)})
			}
		}
	}

	return scoredDocuments
}
