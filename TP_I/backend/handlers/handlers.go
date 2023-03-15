package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Bernardo46-2/AEDS-III/crud"
	"github.com/Bernardo46-2/AEDS-III/dataManager"
	"github.com/Bernardo46-2/AEDS-III/logger"
	"github.com/Bernardo46-2/AEDS-III/models"
	"github.com/Bernardo46-2/AEDS-III/utils"
)

func GetPokemon(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	pokemon, err := crud.Read(id)

	if err != nil {
		writeError(w, http.StatusInternalServerError, 2)
		return
	}

	logger.Println("GET", "Id de numero "+strconv.Itoa(id))
	writeJson(w, pokemon)
}

func GetAllPokemon(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	pokemon, err := crud.ReadAll(id)

	if err != nil {
		writeError(w, http.StatusInternalServerError, 2)
		return
	}

	writeJson(w, pokemon)
}

func PostPokemon(w http.ResponseWriter, r *http.Request) {
	var pokemon models.Pokemon

	err := json.NewDecoder(r.Body).Decode(&pokemon)
	defer r.Body.Close()
	if err != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	id, err := crud.Create(pokemon)

	if err != nil {
		writeError(w, http.StatusInternalServerError, 3)
		return
	}

	pokemonID := models.PokemonID{ID: id}
	writeJson(w, pokemonID)
}

func PutPokemon(w http.ResponseWriter, r *http.Request) {
	var pokemon models.Pokemon

	err := json.NewDecoder(r.Body).Decode(&pokemon)
	defer r.Body.Close()

	if err != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	err = crud.Update(pokemon)

	if err != nil {
		writeError(w, http.StatusInternalServerError)
		return
	}

	writeSuccess(w, 4)
}

func DeletePokemon(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	_, err := crud.Delete(id)

	if err != nil {
		writeError(w, http.StatusInternalServerError, 5)
		return
	}

	writeSuccess(w, 5)
}

func LoadDatabase(w http.ResponseWriter, r *http.Request) {
	dataManager.ImportCSV().CsvToBin()

	writeSuccess(w, 6)
	logger.Println("INFO", "Database Recarregada")
}

func ToKatakana(w http.ResponseWriter, r *http.Request) {
	stringToConvert := r.URL.Query().Get("stringToConvert")

	convertedString := utils.ToKatakana(stringToConvert)

	writeJson(w, convertedString)
}

func IntercalacaoComum(w http.ResponseWriter, r *http.Request) {
	dataManager.IntercalacaoBalanceadaComum()
	writeSuccess(w, 7)
}

func IntercalacaoVariavel(w http.ResponseWriter, r *http.Request) {
	dataManager.IntercalacaoBalanceadaVariavel()
	writeSuccess(w, 8)
}

func SelecaoPorSubstituicao(w http.ResponseWriter, r *http.Request) {
	dataManager.IntercalacaoPorSubstituicao()
	writeSuccess(w, 9)
}

func writeError(w http.ResponseWriter, codes ...int) {
	w.Header().Set("Content-Type", "application/json")
	code := codes[0]
	w.WriteHeader(code)
	if len(codes) > 1 {
		code = codes[1]
	}
	json.NewEncoder(w).Encode(models.ErrorResponse(code))
}

func writeSuccess(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.SuccessResponse(code))
}

func writeJson(w http.ResponseWriter, v any) {
	jsonData, err := json.Marshal(v)
	if err != nil {
		writeError(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
