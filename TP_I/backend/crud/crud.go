package crud

import (
	"github.com/Bernardo46-2/AEDS-III/dataManager"
	"github.com/Bernardo46-2/AEDS-III/models"
)

func Create(pokemon models.Pokemon) (id int, err error) {
	id, _, _ = dataManager.NumRegistros()
	id++
	pokemon.Numero = int32(id)

	pokemon.CalculateSize()

	pokeBytes := pokemon.ToBytes()

	err = dataManager.AppendPokemon(pokeBytes)

	return
}

func Read(id int) (pokemon models.Pokemon, err error) {
	pokemon, _, err = dataManager.ReadBinToPoke(id)
	return
}

func ReadAll(page int) (pokemon []models.Pokemon, err error) {
	numRegistros, _, _ := dataManager.NumRegistros()
	i := page * 60

	for total := 0; total < 60 && i < numRegistros; i++ {
		tmp, _, _ := dataManager.ReadBinToPoke(i)
		if tmp.Numero != 0 {
			pokemon = append(pokemon, tmp)
			total++
		}
	}
	return
}

func Update(pokemon models.Pokemon) (err error) {
	_, pos, err := dataManager.ReadBinToPoke(int(pokemon.Numero))

	if err != nil {
		return err
	}

	pokemon.CalculateSize()
	pokeBytes := pokemon.ToBytes()
	if err = dataManager.DeletarPokemon(pos); err != nil {
		return err
	}

	if err = dataManager.AppendPokemon(pokeBytes); err != nil {
		return err
	}

	return err
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
