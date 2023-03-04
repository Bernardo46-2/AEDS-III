package main

import (
	"fmt"
)

func create() {
	pokemon := readPokemon()
	fmt.Printf("%s", pokemon.ToString())
	pokeBytes := pokemon.ToBytes()

	if err := AppendPokemon(pokeBytes); err != nil {
		fmt.Println("\n", err)
	} else {
		fmt.Println("Pokemon inserido com sucesso")
	}
}

func read() {
	pokemon, err, _ := readBinToPoke(lerInt("Digite o numero da pokedex a pesquisar:\n"))
	if err != nil {
		fmt.Printf("\n%s\n", err)
	} else {
		fmt.Printf("\n" + pokemon.ToString())
	}
}

func update() {
	pokemon, err, pos := readBinToPoke(lerInt("Digite o numero da pokedex a pesquisar:\n"))

	if err != nil {
		if err.Error() == "Pokemon não encontrado" {
			fmt.Println("Criar Pokemon")
		} else {
			fmt.Printf("\n%s\n", err)
			return
		}
	}

	pokemon.alterarCampo()
	pokeBytes := pokemon.ToBytes()
	if err := deletarPokemon(pos); err != nil {
		fmt.Printf("Erro ao alterar\n%s\n", err)
		return
	}

	if err := AppendPokemon(pokeBytes); err != nil {
		fmt.Println("\n", err)
		return
	}

	fmt.Println("Pokemon alterado com sucesso")

}

func delete() {
	pokemon, err, pos := readBinToPoke(lerInt("Digite o numero da pokedex a deletar:\n"))
	if err != nil {
		fmt.Printf("Erro ao excluir\n%s\n", err)
	} else {
		if err = deletarPokemon(pos); err != nil {
			fmt.Printf("Erro ao excluir\n%s\n", err)
		}
		// AlterarNumRegistros(-1)
		fmt.Printf("Pokemon %s excluido com exito", pokemon.Nome)
	}
}
