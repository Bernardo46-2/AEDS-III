package invertedIndex

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Posting struct {
	DocumentID int64
	Frequency  int
}

type InvertedIndex struct {
	Index map[string][]Posting
}

type Reader interface {
	ReadNextGeneric() (any, bool, int64, error)
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

func (ii *InvertedIndex) Print() {
	for word, postings := range ii.Index {
		fmt.Printf("%s: ", word)
		for _, posting := range postings {
			fmt.Printf("(%d,%d) ", posting.DocumentID, posting.Frequency)
		}
		fmt.Println()
	}
}

func (ii *InvertedIndex) RemoveHighFrequencyTerms(percentageThreshold float64) {
	percentageThreshold /= 100
	// 1. Calcule a frequência total de todas as palavras no índice
	totalFrequency := 0
	for _, postings := range ii.Index {
		for _, posting := range postings {
			totalFrequency += posting.Frequency
		}
	}

	// 2. Para cada palavra, calcule sua frequência relativa e verifique se excede o limite
	for word, postings := range ii.Index {
		wordFrequency := 0
		for _, posting := range postings {
			wordFrequency += posting.Frequency
		}

		frequencyRatio := float64(wordFrequency) / float64(totalFrequency)

		if frequencyRatio > percentageThreshold {
			// fmt.Printf("Removendo termo '%s' com frequência total %d (frequência relativa: %f)\n", word, wordFrequency, frequencyRatio)
			delete(ii.Index, word)
		}
	}
}

func (ii *InvertedIndex) writeFile(path string, field string) error {
	filePath := filepath.Join(path, "invertedIndex")
	os.MkdirAll(filePath, 0755)
	fieldPath := filepath.Join(filePath, field+".bin")
	file, err := os.Create(fieldPath)
	if err != nil {
		return fmt.Errorf("error creating file: %s", err)
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	err = enc.Encode(ii)
	if err != nil {
		return fmt.Errorf("error encoding: %s", err)
	}

	return nil
}

func readFile(field string, path string) *InvertedIndex {
	fieldPath := filepath.Join(path, "invertedIndex", field+".bin")
	file, err := os.Open(fieldPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	var invIndex InvertedIndex
	dec := gob.NewDecoder(file)
	err = dec.Decode(&invIndex)
	if err != nil {
		fmt.Println("Error decoding:", err)
		return nil
	}

	return &invIndex
}

// ======================================= Crud ======================================== //

func New(controler Reader, fieldToIndex string, path string) error {
	invIndex := NewInvertedIndex()

	for {
		objInterface, isDead, _, err := controler.ReadNextGeneric()
		if err != nil {
			break
		}

		obj, ok := objInterface.(IndexableObject)
		if !ok {
			return fmt.Errorf("failed to convert object to IndexableObject\n%+v", objInterface)
		}

		if !isDead {
			content := obj.GetField(fieldToIndex)
			words := Tokenize(content)
			id, _ := strconv.ParseInt(obj.GetField("id"), 10, 64)
			invIndex.AddDocument(id, words)
		}
	}

	invIndex.RemoveHighFrequencyTerms(0.3)

	return invIndex.writeFile(path, fieldToIndex)
}

func Read(path string, field string, keys ...string) (documentIDs []int64) {
	invIndex := readFile(field, path)

	rawKeys := strings.Join(keys, " ")
	token := Tokenize(rawKeys)

	// Armazenar a soma das frequências para cada DocumentID
	frequencies := make(map[int64]int)

	for _, key := range token {
		postings, found := invIndex.Index[key]
		if found {
			for _, posting := range postings {
				frequencies[posting.DocumentID] += posting.Frequency
			}
		}
	}

	// Criar um slice de DocumentFrequency e preencher com DocumentID e soma das frequências
	type DocumentFrequency struct {
		DocumentID int64
		Frequency  int
	}

	docFrequencies := make([]DocumentFrequency, 0, len(frequencies))
	for id, freq := range frequencies {
		docFrequencies = append(docFrequencies, DocumentFrequency{DocumentID: id, Frequency: freq})
	}

	// Ordenar o slice de DocumentFrequency por frequência em ordem decrescente
	sort.Slice(docFrequencies, func(i, j int) bool {
		return docFrequencies[i].Frequency > docFrequencies[j].Frequency
	})

	// Extrair apenas os IDs dos documentos e retornar
	documentIDs = make([]int64, len(docFrequencies))
	for i, docFreq := range docFrequencies {
		documentIDs[i] = docFreq.DocumentID
	}

	return documentIDs
}
