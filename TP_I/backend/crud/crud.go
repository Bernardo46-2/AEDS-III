package crud

import (
	"fmt"

	"github.com/Bernardo46-2/AEDS-III/dataManager"
	"github.com/Bernardo46-2/AEDS-III/models"
)

func Create(pokemon models.Pokemon) (id int, err error) {
	pokemon.CalculateSize()

	fmt.Printf("%s", pokemon.ToString())
	pokeBytes := pokemon.ToBytes()

	id, err = dataManager.AppendPokemon(pokeBytes)

	return
}

func Read(id int) (pokemon models.Pokemon, err error) {
	pokemon, _, err = dataManager.ReadBinToPoke(id)
	return
}

func Update(pokemon models.Pokemon) (err error) {
	pokemon, pos, err := dataManager.ReadBinToPoke(int(pokemon.Numero))

	if err != nil {
		return
	}

	pokeBytes := pokemon.ToBytes()
	if err = dataManager.DeletarPokemon(pos); err != nil {
		return
	}

	if _, err = dataManager.AppendPokemon(pokeBytes); err != nil {
		return
	}

	return
}

func Delete(id int) (pokemon models.Pokemon, err error) {
	pokemon, pos, err := dataManager.ReadBinToPoke(id)
	if err != nil {
		return
	}

	if err = dataManager.DeletarPokemon(pos); err != nil {
		return
	}

	return
}
