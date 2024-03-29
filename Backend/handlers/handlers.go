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
	"github.com/Bernardo46-2/AEDS-III/data/sorts"
	"github.com/Bernardo46-2/AEDS-III/logger"
	"github.com/Bernardo46-2/AEDS-III/models"
	"github.com/Bernardo46-2/AEDS-III/service"
	"github.com/Bernardo46-2/AEDS-III/utils"
)

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

// GetIdList faz a chamada do respectivo metodo que recupera
// todos os ids contidos na database
func GetIdList(w http.ResponseWriter, r *http.Request) {
	// Recuperar IDs
	idList, err := service.GetIdList()

	// Resposta
	if err != nil {
		writeError(w, http.StatusInternalServerError, 2)
		return
	}

	writeJson(w, idList)
}

// GetList recupera a lista de pokemons de acordo com a lista de ids fornecidos
// e também passa o metodo de pesquisa necessario.
//
// Por fim faz parse do objeto contendo a lista e o tempo de pesquisa para JSON
func GetList(w http.ResponseWriter, r *http.Request) {
	// struct de retorno para conversao em JSON
	type retorno struct {
		Pokemons []models.Pokemon `json:"pokemons"`
		Time     int64            `json:"time"`
	}

	method, _ := strconv.Atoi(r.URL.Query().Get("method"))

	// Recuperando lista de argumentos
	var list []int64
	if json.NewDecoder(r.Body).Decode(&list) != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	pokeList, time, err := service.GetList(list, method)

	// Resposta
	if err != nil {
		writeError(w, http.StatusInternalServerError, 2)
		return
	}

	writeJson(w, retorno{
		Pokemons: pokeList,
		Time:     time,
	})
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
	service.ReconstruirIndices()

	// Resposta
	writeSuccess(w, 6)
	logger.Println("INFO", "Database Carregada")
	logger.Println("INFO", "Hash Dinamica Criada")
	logger.Println("INFO", "B Tree Criada")
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

// Ordenacao faz a chamada do devido metodo de ordenacao indexado
// atraves de uma hash de funcoes
//
// TODO: Ordenacao por substituicao nao esta funcionando corretamente
func Ordenacao(w http.ResponseWriter, r *http.Request) {
	// Recuperar metodo
	metodo, _ := strconv.Atoi(r.URL.Query().Get("metodo"))

	sorts.SortingFunctions[metodo]()

	// Reconstruir Indices
	service.ReconstruirIndices()

	// Resposta
	writeSuccess(w, 7)
	logger.Println("INFO", "Database Ordenada com sucesso!")
}

// MergeSearch faz a chamada do metodo de pesquisa com ordenacao por
// incidencia e retorna a lista de ids ordenados e o respectivo tempo de
// execucao dos algoritmos
func MergeSearch(w http.ResponseWriter, r *http.Request) {
	// struct para conversao dos dados em json
	type retornoIndexacao struct {
		Pokemons []int64 `json:"ids"`
		Time     int64   `json:"time"`
	}

	var req service.SearchRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		println("err: ", err.Error())
	}

	// Pesquisa os valores no indice
	idList, duration, err := service.MergeSearch(req)

	// Resposta
	if err != nil {
		writeError(w, http.StatusInternalServerError, 2)
		return
	}
	if len(idList) == 0 {
		writeError(w, 2, 2)
		return
	}

	writeJson(w, retornoIndexacao{idList, duration})
}

// Encrypt faz o desempacotamento da requisicao para a chamada
// da criptografia
func Encrypt(w http.ResponseWriter, r *http.Request) {
	// Recuperar metodo
	metodo, _ := strconv.Atoi(r.URL.Query().Get("metodo"))

	k := service.Encrypt(metodo)

	// Resposta
	writeJson(w, k)
	logger.Println("INFO", "Database encriptada!")
}

// Encrypt faz o desempacotamento da requisicao para a chamada
// da descriptografia
func Decrypt(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Key string `json:"key"`
	}
	var requestBody RequestBody
	json.NewDecoder(r.Body).Decode(&requestBody)

	// Recuperar metodo
	metodo, _ := strconv.Atoi(r.URL.Query().Get("metodo"))

	ok := service.Decrypt(metodo, requestBody.Key)

	// Resposta
	if ok {
		writeSuccess(w, 10)
		logger.Println("INFO", "Database decriptada!")
	} else {
		writeSuccess(w, 9)
		logger.Println("INFO", "Chave de criptografia invalida!")
	}
}

// Zip faz o desempacotamento da requisicao para a chamada
// da compressao
func Zip(w http.ResponseWriter, r *http.Request) {
	// Recuperar metodo
	metodo, _ := strconv.Atoi(r.URL.Query().Get("metodo"))

	service.Zip(metodo)

	// Resposta
	writeSuccess(w, 11)
	logger.Println("INFO", "Database comprimida!")
}

// Unzip faz o desempacotamento da requisicao para a chamada
// da descompressao
func Unzip(w http.ResponseWriter, r *http.Request) {
	// Recuperar metodo
	metodo, _ := strconv.Atoi(r.URL.Query().Get("metodo"))

	service.Unzip(metodo)

	// Resposta
	writeSuccess(w, 12)
	logger.Println("INFO", "Database comprimida!")
}
