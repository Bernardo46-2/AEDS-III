// O pacote service realiza a conversa entre as requisiçoes e o DataManager
// recebendo dados ja em formato struct e fazendo as devidas chamadas de ediçao
// no arquivo binario
package service

import (
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/Bernardo46-2/AEDS-III/data/binManager"
	"github.com/Bernardo46-2/AEDS-III/data/indexes/hashing"
	"github.com/Bernardo46-2/AEDS-III/data/indexes/invertedIndex"
	"github.com/Bernardo46-2/AEDS-III/models"
	"github.com/Bernardo46-2/AEDS-III/utils"
)

// ReadPagesNumber retorna o numero de paginas disponiveis para a
// exibiçao dos pokemons na tela inicial do site, como um menu
// de navegação entre paginas
func ReadPagesNumber() (numeroPaginas int, err error) {
	// Recuperação do numero de registros totais
	numRegistros, _, _ := binManager.NumRegistros()

	// calcula e retorna o total
	numeroPaginas = int(math.Ceil((float64(numRegistros) / float64(60))))
	return
}

func GetIdList() (ids []int32, err error) {
	c, _ := binManager.InicializarControleLeitura(binManager.BIN_FILE)
	defer c.Close()

	for {
		err = c.ReadNext()
		if err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
				break
			}
		}
		if !c.RegistroAtual.IsDead() {
			ids = append(ids, c.RegistroAtual.Pokemon.Numero)
		}
	}

	utils.InsertionSort(ids)

	return
}

func GetList(idList []int64) (pokeList []models.Pokemon, err error) {
	c, _ := binManager.InicializarControleLeitura(binManager.BIN_FILE)
	defer c.Close()

	for _, id := range idList {
		pos, err := hashing.HashRead(id, binManager.FILES_PATH, "hashIndex")
		if err == nil {
			pokeList = append(pokeList, c.ReadTarget(pos))
		}
	}

	return
}

// Create adiciona um novo pokemon ao banco de dados.
//
// Recebe um modelo pokemon e serializa para inserir
// Por fim retorna o ID do pokemon criado e erro se houver.
//
// tambem realiza: HashCreate
func Create(pokemon models.Pokemon) (int, error) {
	// Recupera o ultimo ID para gerar o proximo
	ultimoID := binManager.GetLastPokemon()
	ultimoID++
	pokemon.Numero = ultimoID

	// Prepara, serializa e insere
	pokemon.CalculateSize()
	pokeBytes := pokemon.ToBytes()
	address, err := binManager.AppendPokemon(pokeBytes)

	// Indice invertido
	fmt.Printf("%+v", pokemon)
	invertedIndex.Create(pokemon, binManager.FILES_PATH, models.PokemonStringFields()...)

	// Tabela Hash
	hashing.HashCreate(int64(pokemon.Numero), address, binManager.FILES_PATH, "hashIndex")

	return int(ultimoID), err
}

// Read recebe o ID de um pokemon, procura no banco de dados atraves do
// indice hash e o retorna, se nao achar gera um erro
func Read(id int) (models.Pokemon, error) {
	pos, err := hashing.HashRead(int64(id), binManager.FILES_PATH, "hashIndex")
	pokemon := binManager.ReadTargetPokemon(pos)
	return pokemon, err
}

// Update atualiza um registro no arquivo binário de acordo com o número do pokemon informado.
// Recebe uma struct do tipo models.Pokemon a ser atualizada.
// Retorna um erro caso ocorra algum problema ao atualizar o registro.
//
// O update é feito deletando um valor e adicionando outro ao final do arquivo.
//
// tambem realiza: HashUpdate
func Update(pokemon models.Pokemon) (err error) {

	// Recupera a posição do id no arquivo
	pos, err := hashing.HashRead(int64(pokemon.Numero), binManager.FILES_PATH, "hashIndex")
	if err != nil {
		return
	}

	// Serializa os dados
	pokemon.CalculateSize()
	pokeBytes := pokemon.ToBytes()

	// Deleta o antigo e insere o novo registro
	err = binManager.DeletarPokemon(pos)
	if err != nil {
		return
	}

	newAddress, err := binManager.AppendPokemon(pokeBytes)
	if err != nil {
		return
	}

	// Indice invertido
	invertedIndex.Update(pokemon, binManager.FILES_PATH, models.PokemonStringFields()...)

	// Tabela Hash
	err = hashing.HashUpdate(int64(pokemon.Numero), newAddress, binManager.FILES_PATH, "hashIndex")

	return
}

// Delete recebe um ID, procura no arquivo e gera a remoçao logica do mesmo
//
// tambem realiza: HashDelete
func Delete(id int) (pokemon models.Pokemon, err error) {
	// Tenta encontrar a posiçao do pokemon no arquivo binario
	var pos int64
	pos, err = hashing.HashRead(int64(id), binManager.FILES_PATH, "hashIndex")
	pokemon = binManager.ReadTargetPokemon(pos)
	if err != nil {
		return
	}

	// Efetiva a remoção logica
	if err = binManager.DeletarPokemon(pos); err != nil {
		return
	}

	// Indice invertido
	invertedIndex.Delete(pokemon, binManager.FILES_PATH, models.PokemonStringFields()...)

	// Tabela Hash
	hashing.HashDelete(int64(pokemon.Numero), binManager.FILES_PATH, "hashIndex")

	return
}
