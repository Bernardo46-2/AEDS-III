package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Bernardo46-2/AEDS-III/crud"
	"github.com/Bernardo46-2/AEDS-III/dataManager"
	"github.com/Bernardo46-2/AEDS-III/models"
	"github.com/Bernardo46-2/AEDS-III/utils"
)

func writeError(w http.ResponseWriter, codes ...int) {
	w.Header().Set("Content-Type", "application/json")
	code := codes[0]
	w.WriteHeader(code)
	if len(codes) > 1 {
		code = codes[1]
	}
	json.NewEncoder(w).Encode(models.ErrorResponse(code))
}

func writeSuccess(w http.ResponseWriter, codes ...int) {
	w.Header().Set("Content-Type", "application/json")
	code := codes[0]
	w.WriteHeader(code)
	if len(codes) > 1 {
		code = codes[1]
	}
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

func GetPokemon(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	pokemon, err := crud.Read(id)

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

	writeSuccess(w, http.StatusOK, 2)
}

func DeletePokemon(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	_, err := crud.Delete(id)

	if err != nil {
		writeError(w, http.StatusInternalServerError, 5)
		return
	}

	writeSuccess(w, 3)
}

func LoadDatabase(w http.ResponseWriter, r *http.Request) {
	dataManager.ImportCSV().CsvToBin()

	writeSuccess(w, 4)
}

func ToKatakana(w http.ResponseWriter, r *http.Request) {
	stringToConvert := r.URL.Query().Get("stringToConvert")

	convertedString := utils.ToKatakana(stringToConvert)

	writeJson(w, convertedString)
}
