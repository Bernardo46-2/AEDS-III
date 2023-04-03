package invertedIndex

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Bernardo46-2/AEDS-III/data/binManager"
)

type Posting struct {
	DocumentID int64
	Frequency  int
}

type InvertedIndex struct {
	Index map[string][]Posting
}

type IndexableObject interface {
	GetField(fieldName string) string
}

func NewInvertedIndex() InvertedIndex {
	return InvertedIndex{
		Index: make(map[string][]Posting),
	}
}

// Esta função aceita um documento como entrada e retorna uma lista de
// tokens (palavras) no documento. Você pode usar essa função para
// pré-processar o texto antes de criar o índice invertido.
func Tokenize(document string) []string {
	// Converte todas as letras para minúsculas
	lowerCaseDocument := strings.ToLower(document)

	// Expressão regular para dividir o texto em palavras, removendo pontuações
	regEx := regexp.MustCompile(`\p{L}+`)
	tokens := regEx.FindAllString(lowerCaseDocument, -1)

	return tokens
}

// Esta função aceita um ID de documento e uma lista de tokens e adiciona os
// tokens ao índice invertido, associando-os ao ID do documento.
func (ii *InvertedIndex) AddDocument(documentID int64, tokens []string) {
	tokenFrequency := make(map[string]int)

	for _, token := range tokens {
		tokenFrequency[token]++
	}

	for token, frequency := range tokenFrequency {
		ii.Index[token] = append(ii.Index[token], Posting{
			DocumentID: documentID,
			Frequency:  frequency,
		})
	}
}

// Esta função remove todas as ocorrências de um documento do índice invertido.
func (ii *InvertedIndex) RemoveDocument(documentID int64) {
	for token, postings := range ii.Index {
		newPostings := []Posting{}
		for _, posting := range postings {
			if posting.DocumentID != documentID {
				newPostings = append(newPostings, posting)
			}
		}

		if len(newPostings) == 0 {
			delete(ii.Index, token)
		} else {
			ii.Index[token] = newPostings
		}
	}
}

func CreateInvertedIndex(controler *binManager.ControleLeitura, fieldToIndex string) (InvertedIndex, error) {
	invIndex := NewInvertedIndex()

	for {
		objInterface, isDead, _, err := controler.ReadNextGeneric()
		if err != nil {
			break
		}

		obj, ok := objInterface.(IndexableObject)
		if !ok {
			fmt.Printf("%+v", objInterface)
			return NewInvertedIndex(), fmt.Errorf("failed to convert object to IndexableObject")
		}

		if !isDead {
			content := obj.GetField(fieldToIndex)
			words := Tokenize(content)
			id, _ := strconv.ParseInt(obj.GetField("id"), 10, 64)
			invIndex.AddDocument(id, words)
		}
	}

	return invIndex, nil
}

func (idx *InvertedIndex) Print() {
	for word, postings := range idx.Index {
		fmt.Printf("%s: ", word)
		for _, posting := range postings {
			fmt.Printf("(%d,%d) ", posting.DocumentID, posting.Frequency)
		}
		fmt.Println()
	}
}
