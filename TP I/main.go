package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const FILE string = "csv/pokedex2.csv"

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
	1: "1996/02/27",
	2: "1999/11/21",
	3: "2002/11/21",
	4: "2006/09/28",
	5: "2010/09/18",
	6: "2013/10/12",
	7: "2016/11/18",
	8: "2019/11/15",
	9: "2022/11/18",
}

func removeAfterSpace(str string) string {
	parts := strings.Split(str, " ")
	return parts[0]
}

type Pokemon struct {
	Numero     int
	Nome       string
	NomeJap    string
	Geracao    int
	Lancamento time.Time
	Especie    string
	Lendario   bool
	Mitico     bool
	Tipo       string
	Atk        int
	Def        int
	Hp         int
	Altura     float64
	Peso       float64
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

	for _, valor := range lines {
		var novoPokemon Pokemon
		novoPokemon.Numero, _ = strconv.Atoi(valor[1])
		novoPokemon.Nome = valor[2]
		novoPokemon.NomeJap = removeAfterSpace(valor[4])
		geracao, _ := strconv.Atoi(valor[5])
		novoPokemon.Geracao = geracao
		novoPokemon.Lancamento, _ = time.Parse("2006/01/02", GenReleaseDates[geracao])
		novoPokemon.Especie = valor[9]
		novoPokemon.Lendario, _ = strconv.ParseBool(valor[7])
		novoPokemon.Mitico, _ = strconv.ParseBool(valor[8])
		novoPokemon.Tipo = valor[11] + valor[12]
		novoPokemon.Atk, _ = strconv.Atoi(valor[21])
		novoPokemon.Def, _ = strconv.Atoi(valor[22])
		novoPokemon.Hp, _ = strconv.Atoi(valor[20])
		novoPokemon.Altura, _ = strconv.ParseFloat(valor[13], 64)
		novoPokemon.Peso, _ = strconv.ParseFloat(valor[14], 64)
		pokemons = append(pokemons, novoPokemon)
	}

	fmt.Printf("\n")
	fmt.Printf("Numero     = %d\n", pokemons[5].Numero)
	fmt.Printf("Nome       = %s\n", pokemons[5].Nome)
	fmt.Printf("NomeJap    = %s\n", pokemons[5].NomeJap)
	fmt.Printf("Geracao    = %d\n", pokemons[5].Geracao)
	fmt.Printf("Lancamento = %s\n", pokemons[5].Lancamento.Format("02/01/2006"))
	fmt.Printf("Especie    = %s\n", pokemons[5].Especie)
	fmt.Printf("Lendario   = %t\n", pokemons[5].Lendario)
	fmt.Printf("Mitico     = %t\n", pokemons[5].Mitico)
	fmt.Printf("Tipo       = %s\n", pokemons[5].Tipo)
	fmt.Printf("Atk        = %d\n", pokemons[5].Atk)
	fmt.Printf("Def        = %d\n", pokemons[5].Def)
	fmt.Printf("Hp         = %d\n", pokemons[5].Hp)
	fmt.Printf("Altura     = %f\n", pokemons[5].Altura)
	fmt.Printf("Peso       = %f\n", pokemons[5].Peso)
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
