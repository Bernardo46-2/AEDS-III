package main

import (
	"fmt"
)

func main() {
	var csvFile CSV
	menu := "1 - Create\n2 - Read\n3 - Update\n4 - Delete\n8 - Convert CSV to Bin\n0 - Exit\n"

	for quit := false; !quit; {
		switch lerInt(menu) {
		case 0:
			fmt.Printf("\nSaindo do programa...\n\n")
			quit = true
		case 1:
			pokemon := readPokemon()
			fmt.Printf("%s", pokemon.ToString())
			// incrementNumRegistros()
		case 2:
			pokemon, err, _ := readBinToPoke(lerInt("Digite o numero da pokedex a pesquisar:\n"))
			if err != nil {
				fmt.Printf("\n%s\n", err)
			} else {
				fmt.Printf(pokemon.ToString())
			}
		case 3:
			fmt.Printf("Update\n")
		case 4:
			pokemon, err, pos := readBinToPoke(lerInt("Digite o numero da pokedex a deletar:\n"))
			if err != nil {
				fmt.Printf("Erro ao excluir\n%s\n", err)
			} else {
				if err = deletarPokemon(pos); err != nil {
					fmt.Printf("Erro ao excluir\n%s\n", err)
				}
				fmt.Printf("Pokemon %s excluido com exito", pokemon.Nome)
			}
		case 8:
			csvFile = importCSV()
			csvFile.CsvToBin()
		default:
			fmt.Println("Opção inválida")
		}
		pause()
	}
}
