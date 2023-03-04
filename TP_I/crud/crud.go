package crud

import (
	"fmt"

	"github.com/Bernardo46-2/AEDS-III/dataManager"
	"github.com/Bernardo46-2/AEDS-III/models"
	"github.com/Bernardo46-2/AEDS-III/utils"
)

func Create() {
	pokemon := models.ReadPokemon()
	fmt.Printf("%s", pokemon.ToString())
	pokeBytes := pokemon.ToBytes()

	if err := dataManager.AppendPokemon(pokeBytes); err != nil {
		fmt.Println("\n", err)
	} else {
		fmt.Println("Pokemon inserido com sucesso")
	}
}

func Read() {
	pokemon, _, err := dataManager.ReadBinToPoke(utils.LerInt("Digite o numero da pokedex a pesquisar:\n"))
	if err != nil {
		fmt.Printf("\n%s\n", err)
	} else {
		fmt.Printf("\n" + pokemon.ToString())
	}
}

func Update() {
	pokemon, pos, err := dataManager.ReadBinToPoke(utils.LerInt("Digite o numero da pokedex a pesquisar:\n"))

	if err != nil {
		if err.Error() == "Pokemon não encontrado" {
			fmt.Println("Criar Pokemon")
		} else {
			fmt.Printf("\n%s\n", err)
			return
		}
	}

	pokemon.AlterarCampo()
	pokeBytes := pokemon.ToBytes()
	if err := dataManager.DeletarPokemon(pos); err != nil {
		fmt.Printf("Erro ao alterar\n%s\n", err)
		return
	}

	if err := dataManager.AppendPokemon(pokeBytes); err != nil {
		fmt.Println("\n", err)
		return
	}

	fmt.Println("Pokemon alterado com sucesso")

}

func Delete() {
	pokemon, pos, err := dataManager.ReadBinToPoke(utils.LerInt("Digite o numero da pokedex a deletar:\n"))
	if err != nil {
		fmt.Printf("Erro ao excluir\n%s\n", err)
	} else {
		if err = dataManager.DeletarPokemon(pos); err != nil {
			fmt.Printf("Erro ao excluir\n%s\n", err)
		}
		// AlterarNumRegistros(-1)
		fmt.Printf("Pokemon %s excluido com exito", pokemon.Nome)
	}
}
