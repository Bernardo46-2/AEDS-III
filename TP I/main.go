package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func pause() {
	var input string
	fmt.Printf("\nPressione Enter para continuar...\n")
	fmt.Scanf("%s\n", &input)
}

func importCSV() {
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

func main() {
	for {
		clearScreen()
		fmt.Printf("1 - Create\n")
		fmt.Printf("2 - Read\n")
		fmt.Printf("3 - Update\n")
		fmt.Printf("4 - Delete\n")
		fmt.Printf("9 - Import CSV\n")
		fmt.Printf("0 - Exit\n")
		fmt.Printf("\n> ")

		var tmp string
		if _, err := fmt.Scanln(&tmp); err != nil {
			fmt.Println("Erro ao ler opção:", err)
		}
		if opcao, err := strconv.Atoi(tmp); err == nil {
			switch opcao {
			case 0:
				fmt.Printf("Saindo do programa...\n\n")
				return
			case 1:
				fmt.Printf("Create\n")
			case 2:
				fmt.Printf("Read\n")
			case 3:
				fmt.Printf("Update\n")
			case 4:
				fmt.Printf("Delete\n")
			case 9:
				importCSV()
			default:
				fmt.Println("Opção inválida")
			}
		} else {
			fmt.Println("Erro ao ler opção:", err)
		}
		pause()
	}
}
