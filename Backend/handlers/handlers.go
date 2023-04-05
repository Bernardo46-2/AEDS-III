// O pacote handlers faz a ligação entre as requisições http e suas respectivas funções
// ligando o service para manipulação do banco de dados, ou chamando diretamente as funções
// de ordenação no binManager
// Handlers também realiza o parsing entre JSON e Objeto
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Bernardo46-2/AEDS-III/data/binManager"
	"github.com/Bernardo46-2/AEDS-III/data/indexes/hashing"
	"github.com/Bernardo46-2/AEDS-III/data/indexes/invertedIndex"
	"github.com/Bernardo46-2/AEDS-III/data/sorts"
	"github.com/Bernardo46-2/AEDS-III/logger"
	"github.com/Bernardo46-2/AEDS-III/models"
	"github.com/Bernardo46-2/AEDS-III/service"
	"github.com/Bernardo46-2/AEDS-III/utils"
)

// GetPagesNumber retorna a quantidade de paginas disponiveis
func GetPagesNumber(w http.ResponseWriter, r *http.Request) {
	// Recuperar ID e ler arquivo
	numeroPaginas, err := service.ReadPagesNumber()

	// Resposta
	if err != nil {
		writeError(w, http.StatusInternalServerError, 2)
		return
	}

	writeJson(w, numeroPaginas)
}

// GetAllPokemon recupera os 60 pokemons a partir do ID fornecido
func GetAllPokemon(w http.ResponseWriter, r *http.Request) {
	// Recuperar ID e ler arquivo
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pokemon, err := service.ReadAll(page)

	// Resposta
	if err != nil {
		writeError(w, http.StatusInternalServerError, 2)
		return
	}

	writeJson(w, pokemon)
}

// GetPokemon recupera o pokemon pelo ID fornecido
func GetPokemon(w http.ResponseWriter, r *http.Request) {
	// recuperar ID e ler do arquivo
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	pokemon, err := service.Read(id)

	// Gera resposta de acordo com o resultado
	if err != nil {
		writeError(w, http.StatusInternalServerError, 2)
		return
	}

	writeJson(w, pokemon)
}

// PostPokemon adiciona o pokemon ao banco de dados
func PostPokemon(w http.ResponseWriter, r *http.Request) {

	// Desserialização
	var pokemon models.Pokemon
	err := json.NewDecoder(r.Body).Decode(&pokemon)
	defer r.Body.Close()
	if err != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	// Create
	id, err := service.Create(pokemon)

	// Resposta
	if err != nil {
		writeError(w, http.StatusInternalServerError, 3)
		return
	}
	pokemonID := models.PokemonID{ID: id}
	writeJson(w, pokemonID)
}

// PutPokemon recebe um json e atualiza o valor no banco de dados
// de acordo com o dado recebido
func PutPokemon(w http.ResponseWriter, r *http.Request) {
	//  Desserialização
	var pokemon models.Pokemon
	err := json.NewDecoder(r.Body).Decode(&pokemon)
	defer r.Body.Close()

	if err != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	// Update
	err = service.Update(pokemon)

	// Resposta
	if err != nil {
		writeError(w, http.StatusInternalServerError)
		return
	}

	writeSuccess(w, 4)
}

// DeletePokemon recebe um ID, pesquisa no banco de dados
// e se existir efetiva sua remoção logica
func DeletePokemon(w http.ResponseWriter, r *http.Request) {
	// Recupera id
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	// Delete
	_, err := service.Delete(id)

	// Resposta
	if err != nil {
		writeError(w, http.StatusInternalServerError, 5)
		return
	}

	writeSuccess(w, 5)
}

// LoadDatabase faz o carregamento do arquivo CSV e o serializa em binario
//
// Tambem é criado indices para: Hash
func LoadDatabase(w http.ResponseWriter, r *http.Request) {
	// CSV
	binManager.ImportCSV().CsvToBin()

	// Reconstruir Indices
	reconstruirIndices()

	// Resposta
	writeSuccess(w, 6)
	logger.Println("INFO", "Database Carregada")
	logger.Println("INFO", "Hash Dinamica Criada")
}

// ToKatakana recebe uma string em alfabeto romato, converte para
// o padrão katakana da linguagem japonesa e retorna a string
// convertida
func ToKatakana(w http.ResponseWriter, r *http.Request) {
	// Intercepta
	stringToConvert := r.URL.Query().Get("stringToConvert")

	// Converte
	convertedString := utils.ToKatakana(stringToConvert)

	// Resposta
	writeJson(w, convertedString)
}

func Ordenacao(w http.ResponseWriter, r *http.Request) {
	// Recuperar metodo
	metodo, _ := strconv.Atoi(r.URL.Query().Get("metodo"))

	sorts.SortingFunctions[metodo]()

	// Reconstruir Indices
	reconstruirIndices()

	// Resposta
	writeSuccess(w, 7)
	logger.Println("INFO", "Database Ordenada com sucesso!")
}

// writeError recebe um erro de http responde e um id de erro interno,
// faz o parsing do modelo e gera uma resposta em formato json com o erro fornecido
func writeError(w http.ResponseWriter, codes ...int) {
	// Preparacao da resposta http
	w.Header().Set("Content-Type", "application/json")
	code := codes[0]
	w.WriteHeader(code)
	if len(codes) > 1 {
		code = codes[1]
	}

	// Gera uma resposta json personalizada
	err := models.ErrorResponse(code)
	json.NewEncoder(w).Encode(err)
	logger.Println("ERROR", fmt.Sprintf("code: %d, message: %s", err.Code, err.Message))
}

// writeSuccess gera uma resposta http de sucesso (200) e
// faz o parsing do modelo de sucesso para uma resposta json com a mensagem
// da ação realizada
func writeSuccess(w http.ResponseWriter, code int) {
	// Preparação da resposta http
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Gera uma resposta json personalizada
	json.NewEncoder(w).Encode(models.SuccessResponse(code))
}

// writeJson recebe qualquer tipo de dado ou struct e serializa o dado
// em formato json, gerando junto uma resposta de sucesso ou erro
func writeJson(w http.ResponseWriter, v any) {
	// Serialização
	jsonData, err := json.Marshal(v)

	// Resposta
	if err != nil {
		writeError(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func reconstruirIndices() {
	// Hashing
	controler, _ := binManager.InicializarControleLeitura(binManager.BIN_FILE)
	defer controler.Close()
	hashing.StartHashFile(controler, 8, binManager.FILES_PATH, "hashIndex")

	// Indice Invertido
	controler.Reset()
	invertedIndex.New(controler, "nome", binManager.FILES_PATH, 0)
	controler.Reset()
	invertedIndex.New(controler, "nomeJap", binManager.FILES_PATH, 0)
	controler.Reset()
	invertedIndex.New(controler, "especie", binManager.FILES_PATH, 0.8)
	controler.Reset()
	invertedIndex.New(controler, "tipo", binManager.FILES_PATH, 0)
	controler.Reset()
	invertedIndex.New(controler, "descricao", binManager.FILES_PATH, 0.8)
}
