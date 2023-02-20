package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const FILE string = "csv/new_pokedex.csv"

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

var GenReleaseDates = map[int]string{
	1: "1996/02/27, ",
	2: "1999/11/21, ",
	3: "2002/11/21, ",
	4: "2006/09/28, ",
	5: "2010/09/18, ",
	6: "2013/10/12, ",
	7: "2016/11/18, ",
	8: "2019/11/15, ",
	9: "2022/11/18, ",
}

type Pokemon struct {
	Numero     int
	Nome       string
	NomeJap    int
	Geracao    int
	Lancamento time.Time
	Especie    string
	Lendario   bool
	Mitico     bool
	Tipo       []string
	Atk        int
	Def        int
	Hp         int
	Altura     float32
	Peso       float32
}

func importCSV() {
	// Abrir o arquivo CSV
	file, err := os.Open(FILE)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()

	// Lendo o conteúdo do arquivo CSV
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return
	}

	pokemons := []Pokemon{}

	/* 	for indice, valor := range lines {
		novoPokemon := Pokemon{Nome: "João", Idade: 30}
		pokemons = append(pokemons, novoPokemon)
	} */
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
