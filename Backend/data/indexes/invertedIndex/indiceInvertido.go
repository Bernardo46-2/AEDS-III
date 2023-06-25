// O pacote invertedIndex fornece uma implementação de um Índice Invertido para
// uso eficiente na recuperação de informações. Um índice invertido é uma estrutura
// de dados projetada para melhorar a velocidade das operações de pesquisa, permitindo
// rápido acesso à lista de documentos que contêm uma determinada palavra.
//
// O pacote invertedIndex cria um mapa na memória principal e, em seguida, salva
// este mapa na memória secundária para uso duradouro em formato GOB
//
// Durante as operações, o pacote recupera os valores salvos da memória secundária
// para a memória principal em formato GOB e realiza a operação requisitada.
//
// Exemplo de uso:
//
//	invertedIndex.New(controler, "nome", binManager.FILES_PATH, 0) // Cria um novo Índice Invertido
//	invertedIndex.Create(pokemon, binManager.FILES_PATH, models.PokeStrings()...) // Adiciona o par documento-palavra ao índice
//
// O pacote tambem implementa um sistema de scoredDocument com operacao de Merge
// para pesquisas baseadas em repeticao de campos
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

// Documento e a sua frequencia
type Posting struct {
	DocumentID int64
	Frequency  int
}

// Dados necessarios para o indice invertido
type InvertedIndex struct {
	Index           map[string][]Posting
	removeThreshold float64
}

// Documento pontuado, utilizado para ordenar os valores
// de acordo com sua frequencia de aparição
type ScoredDocument struct {
	DocumentID int64
	Score      int
}

// Interface para leitura da database
type Reader interface {
	ReadNextGeneric() (any, bool, int64, error)
}

// Interface para recuperacao do campo do objeto indexavel
type IndexableObject interface {
	GetField(fieldName string) string
}

// NewInvertedIndex é um wrapper para a inicializacao de uma
// estrutura de indice invertido
func NewInvertedIndex() InvertedIndex {
	return InvertedIndex{
		Index:           make(map[string][]Posting),
		removeThreshold: 0,
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

// Print realiza a traducao para linguagem legivel do conteudo de um indice invertido.
func (ii *InvertedIndex) Print() {
	for word, postings := range ii.Index {
		fmt.Printf("%s: ", word)
		for _, posting := range postings {
			fmt.Printf("(%d,%d) ", posting.DocumentID, posting.Frequency)
		}
		fmt.Println()
	}
}

// RemoveHighFrequencyTerms realiza a limpeza de valores que se repetem muito,
// o que significa valores que não possuem significado de busca, como por exemplo
// "this, or, it, and"
//
// Deve ser fornecido um percentual de limite para a remocao.
// O valor ideal sugerido é 0.5
func (ii *InvertedIndex) RemoveHighFrequencyTerms(percentageThreshold float64) {
	if percentageThreshold == 0 {
		return
	}
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

// writeFile realiza a escrita do documento de indice invertido em um path fornecido
// e grava os arquivos em formato padrao GOB
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

// readFile realiza a leitura de um arquivo contendo um indice invertido em formato GOB
// e realiza o parse, retornando um ponteiro para o objeto
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

// NewScoredDocumentSlice é um wraper para criar um scoredDocument
func NewScoredDocumentSlice(id int64, score int) (slice []ScoredDocument) {
	return append(slice, ScoredDocument{
		DocumentID: id,
		Score:      score,
	})
}

// ======================================= Crud ======================================== //

// New inicializa a criacao de um indice invertido utilizando um controler por interface o
// campo a ser indexado, o endereco do arquivo e o threshold de limiar para palavras repetidas
//
// # Interfaces necessarias
//
// Controler:
//
//	ReadNextGeneric() (any, bool, int64, error)
//
// Objeto:
//
//	GetField(fieldName string) string
func New(controler Reader, fieldToIndex string, path string, removeFrequency float64) error {
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

	invIndex.removeThreshold = removeFrequency
	invIndex.RemoveHighFrequencyTerms(removeFrequency)

	return invIndex.writeFile(path, fieldToIndex)
}

// Merge realiza a uniao de diferentes listas de scoredDocuments pelos seus IDS
// e somando seus scores.
//
// Por fim os documentos finais sao ordenados de acordo com o score resultado e retornados
func Merge(scoredDocumentsLists ...[]ScoredDocument) []ScoredDocument {
	scoreMap := make(map[int64]int)

	// Soma os scores dos documentos em todas as listas
	for _, scoredDocuments := range scoredDocumentsLists {
		for _, scoredDocument := range scoredDocuments {
			scoreMap[scoredDocument.DocumentID] += scoredDocument.Score
		}
	}

	// Cria um slice com os documentos e seus scores acumulados
	mergedScoredDocuments := make([]ScoredDocument, 0, len(scoreMap))
	for documentID, score := range scoreMap {
		mergedScoredDocuments = append(mergedScoredDocuments, ScoredDocument{DocumentID: documentID, Score: score})
	}

	// Ordena o slice de acordo com a pontuação em ordem decrescente
	sort.Slice(mergedScoredDocuments, func(i, j int) bool {
		return mergedScoredDocuments[i].Score > mergedScoredDocuments[j].Score
	})

	return mergedScoredDocuments
}

// Create faz a abertura de um indice invertido e adiciona um valor a ele
func Create(myObj any, path string, fields ...string) error {
	obj, ok := myObj.(IndexableObject)
	if !ok {
		return fmt.Errorf("failed to convert object to IndexableObject")
	}

	id, _ := strconv.ParseInt(obj.GetField("id"), 10, 64)

	for _, field := range fields {
		// Ler o índice invertido atual do arquivo
		invIndex := readFile(field, path)

		// Adicionar o novo documento ao índice invertido
		content := obj.GetField(field)
		words := Tokenize(content)
		invIndex.AddDocument(id, words)

		// Escrever o índice invertido atualizado de volta ao arquivo
		err := invIndex.writeFile(path, field)
		if err != nil {
			fmt.Printf("error creating field '%s': %v", field, err)
			return fmt.Errorf("error creating field '%s': %v", field, err)
		}
	}

	return nil
}

// Read realiza a busca de um documento em um arquivo
func Read(path string, field string, keys ...string) (scoredDocuments []ScoredDocument) {
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

	// Criar um slice de ScoredDocument e preencher com DocumentID e soma das frequências
	scoredDocuments = make([]ScoredDocument, 0, len(frequencies))
	for id, freq := range frequencies {
		scoredDocuments = append(scoredDocuments, ScoredDocument{DocumentID: id, Score: freq})
	}

	// Ordenar o slice de ScoredDocument por score em ordem decrescente
	sort.Slice(scoredDocuments, func(i, j int) bool {
		return scoredDocuments[i].Score > scoredDocuments[j].Score
	})

	return scoredDocuments
}

// Update realiza a atualizacao de um objeto indexavel no indice invertido
func Update(obj IndexableObject, path string, fields ...string) error {
	id, _ := strconv.ParseInt(obj.GetField("id"), 10, 64)

	for _, field := range fields {
		// Ler o índice invertido atual do arquivo
		invIndex := readFile(field, path)

		// Remover o documento existente do índice invertido
		invIndex.RemoveDocument(id)

		// Tokenizar o novo conteúdo do campo e adicionar o novo documento ao índice invertido
		content := obj.GetField(field)
		words := Tokenize(content)
		invIndex.AddDocument(id, words)

		// Escrever o índice invertido atualizado de volta ao arquivo
		invIndex.RemoveHighFrequencyTerms(invIndex.removeThreshold)
		err := invIndex.writeFile(path, field)
		if err != nil {
			return fmt.Errorf("error updating field '%s': %v", field, err)
		}
	}

	return nil
}

// Delete remove um objeto de um indice invertido
func Delete(obj IndexableObject, path string, fields ...string) error {
	id, _ := strconv.ParseInt(obj.GetField("id"), 10, 64)

	for _, field := range fields {
		// Ler o índice invertido atual do arquivo
		invIndex := readFile(field, path)

		// Remover o documento existente do índice invertido
		invIndex.RemoveDocument(id)

		// Escrever o índice invertido atualizado de volta ao arquivo
		err := invIndex.writeFile(path, field)
		if err != nil {
			return fmt.Errorf("error deleting field '%s': %v", field, err)
		}
	}

	return nil
}
