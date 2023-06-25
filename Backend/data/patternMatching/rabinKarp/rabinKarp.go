// Package rabinKarp implementa o algoritmo de busca de strings Rabin-Karp. Este pacote oferece
// uma forma eficiente e prática de encontrar a ocorrência de uma string de padrão dentro de um texto.
//
// O algoritmo Rabin-Karp é um algoritmo de busca de string que usa hash para encontrar
// uma string de padrão dentro de um texto. A vantagem do Rabin-Karp é que ele processa
// múltiplas ocorrências do padrão simultaneamente, tornando-o adequado para a busca de padrões
// em textos longos.
//
// A função principal deste pacote é a função RabinKarp, que recebe um padrão e um texto e
// retorna os índices do texto onde o padrão foi encontrado. Para cada correspondência de hash,
// a função RabinKarp verifica se os caracteres correspondem exatamente, eliminando assim a
// possibilidade de falsos positivos.
//
// Exemplo de uso:
//
//	pattern := "abc"
//	text := "abcdeabcdabcabc"
//	indexes := rabinKarp.RabinKarp(pattern, text)
//	// indexes: [0 5 9 12]
//
// A função RabinKarp retorna um slice vazio se o padrão não for encontrado no texto. Este pacote
// converte tanto o texto quanto o padrão para minúsculas antes de realizar a comparação, portanto,
// a busca é insensível a maiúsculas e minúsculas.
//
// Este pacote não oferece suporte à busca de padrões que contenham caracteres especiais.
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

// RabinKarp implementa o algoritmo de busca de strings Rabin-Karp, que busca
// uma string 'pattern' em um texto 'text' usando cálculos de hash. O algoritmo
// desliza uma janela do tamanho do padrão sobre o texto, recalculando o hash a
// cada passo e comparando com o hash do padrão. Isso faz do Rabin-Karp um algoritmo
// eficiente para a busca de padrões fixos em um texto.
//
// O 'pattern' e o 'text' são convertidos para letras minúsculas antes da busca para
// assegurar que a busca seja insensível à caixa.
//
// A função retorna um slice de inteiros contendo todos os índices no 'text' onde o
// 'pattern' foi encontrado. Se o 'pattern' não for encontrado no 'text', um slice
// vazio é retornado.
//
// Note que a função de hash usada pelo algoritmo Rabin-Karp pode levar a colisões de
// hash, o que significa que em alguns raros casos a função pode retornar um falso
// positivo, onde parece que o padrão foi encontrado, mas na verdade foi apenas uma
// coincidência de hash. Nestes casos de colisao os caracteres são testados manualmente
// um a um.
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

// SearchPokemon realiza uma busca por um termo específico (search) em um campo específico
// (field) dos registros de Pokemon, retornando os documentos que contêm o Id do pokemon
// que possui aquele termo e a frequencia de aparição.
//
// Param:
//
//	search: O termo de busca que será procurado no campo especificado dos registros de Pokemon.
//	field: O campo dos registros de Pokemon onde o termo de busca será procurado.
//
// Retorno:
//
//	Um slice de ScoredDocument, onde cada ScoredDocument contém o ID do documento e a pontuação,
//	que é o número de ocorrências do termo de busca no campo especificado. Retorna um slice vazio
//	se o termo de busca não for encontrado.
//
// A função inicia o controlador de leitura e lê todos os registros de Pokemon um por um. Para
// cada registro de Pokemon, a função realiza a busca pelo termo de busca no campo especificado
// usando a função RabinKarp. Se o termo de busca for encontrado, um ScoredDocument é criado
// com o ID do documento sendo o número do Pokemon e a pontuação sendo o número de ocorrências
// do termo de busca. O ScoredDocument é então adicionado ao slice de ScoredDocument.
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
