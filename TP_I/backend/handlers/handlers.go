package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Bernardo46-2/AEDS-III/crud"
	"github.com/Bernardo46-2/AEDS-III/dataManager"
	"github.com/Bernardo46-2/AEDS-III/models"
)

func GetPokemon(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	// Procura pelo id passado
	pokemon, err := crud.Read(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retorna as informações relevantes em um formato JSON
	json.NewEncoder(w).Encode(pokemon)
}

func PostPokemon(w http.ResponseWriter, r *http.Request) {
	var pokemon models.Pokemon

	err := json.NewDecoder(r.Body).Decode(&pokemon)

	defer r.Body.Close()
	if err != nil {
		// Trata o erro
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Erro ao converter o JSON para Pokemon"))
		return
	}

	err = crud.Create(pokemon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Requisição recebida com sucesso!"))
}

func PutPokemon(w http.ResponseWriter, r *http.Request) {
	var pokemon models.Pokemon

	err := json.NewDecoder(r.Body).Decode(&pokemon)

	defer r.Body.Close()
	if err != nil {
		// Trata o erro
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Erro ao converter o JSON para Pokemon"))
		return
	}

	err = crud.Update(pokemon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Requisição recebida com sucesso!"))
}

func DeletePokemon(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	// Deleta o registro pelo id passado
	_, err := crud.Delete(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Registro deletado com sucesso!"))
}

func LoadDatabase(w http.ResponseWriter, r *http.Request) {
	// Leitura do arquivo CSV e popule o banco de dados binário
	dataManager.ImportCSV().CsvToBin()

	// Retorne uma resposta indicando que a operação foi concluída
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Banco de dados carregado com sucesso!"))
}
