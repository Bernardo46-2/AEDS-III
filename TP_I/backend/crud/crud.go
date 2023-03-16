// O pacote Crud realiza a conversa entre as requisiçoes e o DataManager
// recebendo dados ja em formato struct e fazendo as devidas chamadas de ediçao
// no arquivo binario
package crud

import (
	"github.com/Bernardo46-2/AEDS-III/dataManager"
	"github.com/Bernardo46-2/AEDS-III/models"
)

// Create adiciona um novo pokemon ao banco de dados.
//
// Recebe um modelo pokemon e serializa para inserir
// Por fim retorna o ID do pokemon criado e erro se houver.
func Create(pokemon models.Pokemon) (int, error) {
	// Gera o id fazendo total de pokemons + 1
	id, _, _ := dataManager.NumRegistros()
	id++
	pokemon.Numero = int32(id)

	// Prepara, serializa e insere
	pokemon.CalculateSize()
	pokeBytes := pokemon.ToBytes()
	err := dataManager.AppendPokemon(pokeBytes)

	return id, err
}

// Read recebe o ID de um pokemon, procura no banco de dados e
// o retorna, se nao achar gera um erro
func Read(id int) (models.Pokemon, error) {
	pokemon, _, err := dataManager.ReadBinToPoke(id)
	return pokemon, err
}

// ReadAll retorna um slice de modelos Pokemon a partir de um ID especificado.
// Se houver a função lê até 60 registros a partir do ID fornecido e adiciona
// a um slice de Pokemon, retornando o slice e um erro, se houver.
func ReadAll(id int) (pokemon []models.Pokemon, err error) {
	// Recuperação do numero de registros totais
	numRegistros, _, _ := dataManager.NumRegistros()

	// Recupera 60 ids enquanto houverem
	for i, total := id+1, 0; total < 60 && i < numRegistros; i++ {
		tmp, _, _ := dataManager.ReadBinToPoke(i)
		if tmp.Numero > 0 {
			pokemon = append(pokemon, tmp)
			total++
		}
	}
	return
}

// Update atualiza um registro no arquivo binário de acordo com o número do pokemon informado.
// Recebe uma struct do tipo models.Pokemon a ser atualizada.
// Retorna um erro caso ocorra algum problema ao atualizar o registro.
func Update(pokemon models.Pokemon) (err error) {

	// Recupera a posição do id no arquivo
	_, pos, err := dataManager.ReadBinToPoke(int(pokemon.Numero))
	if err != nil {
		return
	}

	// Serializa os dados
	pokemon.CalculateSize()
	pokeBytes := pokemon.ToBytes()

	// Deleta o antigo e insere o novo registro
	if err = dataManager.DeletarPokemon(pos); err != nil {
		return err
	}

	if err = dataManager.AppendPokemon(pokeBytes); err != nil {
		return err
	}

	return
}

// Delete recebe um ID, procura no arquivo e gera a remoçao logica do mesmo
func Delete(id int) (pokemon models.Pokemon, err error) {
	// Tenta encontrar a posiçao do pokemon no arquivo binario
	pokemon, pos, err := dataManager.ReadBinToPoke(id)
	if err != nil {
		return
	}

	// Efetiva a remoção logica
	if err = dataManager.DeletarPokemon(pos); err != nil {
		return
	}

	return
}
