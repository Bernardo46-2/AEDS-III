package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func main() {
	// Abrir o arquivo CSV
	file, err := os.Open("csv/pokedex.csv")
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()

	// Criar um leitor CSV
	reader := csv.NewReader(file)

	// Ler as linhas do arquivo
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return
	}

	// Imprimir o conteúdo do arquivo
	for _, record := range records {
		fmt.Println(record)
	}
}
