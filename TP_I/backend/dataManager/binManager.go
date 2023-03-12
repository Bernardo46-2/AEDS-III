package dataManager

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/Bernardo46-2/AEDS-III/models"
)

// readBinToPoke lê um arquivo binário com informações de Pokémons e retorna
// um Pokémon com o número especificado. Caso o número não seja encontrado,
// o Pokémon retornado terá seu número igual a -1.
//
// O arquivo binário deve conter uma sequência de registros de tamanho variável.
// Cada registro deve conter um cabeçalho de 4 bytes representando o tamanho
// do registro em bytes, seguido dos dados do Pokémon.
//
// A função retorna um Pokémon com as informações encontradas no arquivo binário.
// Se o número do Pokémon não for encontrado, o número retornado será -1.
func ReadBinToPoke(id int) (models.Pokemon, int64, error) {
	// Abre o arquivo binário
	file, err := os.Open(BIN_FILE)
	if err != nil {
		return models.Pokemon{}, 0, fmt.Errorf("erro ao abrir o arquivo: %v", err)
	}
	defer file.Close()

	// Lê o número de entradas no arquivo
	var numEntradas int32
	if err = binary.Read(file, binary.LittleEndian, &numEntradas); err != nil {
		pos, _ := file.Seek(0, io.SeekCurrent)
		return models.Pokemon{}, pos, fmt.Errorf("erro ao ler número de entradas: %v Linha Corrompida: %d", err, pos)
	}

	// Percorre as entradas do arquivo
	for i := 0; i < int(numEntradas); i++ {
		// Grava a localização do inicio do registro
		inicioRegistro, _ := file.Seek(0, io.SeekCurrent)

		pokemonAtual, inicioRegistro, _ := readRegistro(file, inicioRegistro)

		// Verifica se o número do Pokémon atual é o procurado
		if pokemonAtual.Numero == int32(id) {
			return pokemonAtual, inicioRegistro, nil
		}
	}

	// Se não encontrou o Pokémon procurado, retorna um erro
	pos, _ := file.Seek(0, io.SeekCurrent)
	return models.Pokemon{}, pos, fmt.Errorf("pokemon não encontrado")
}

func readRegistro(file *os.File, inicioRegistro int64) (models.Pokemon, int64, error) {
	file.Seek(inicioRegistro, io.SeekStart)

	// Lê e confere a lapide do arquivo
	var lapide int32
	if err := binary.Read(file, binary.LittleEndian, &lapide); err != nil {
		return models.Pokemon{}, inicioRegistro, fmt.Errorf("erro ao ler lapide: %v Linha Corrompida: %d", err, inicioRegistro)
	}

	// Lê o tamanho do registro atual
	var tamReg int32
	if err := binary.Read(file, binary.LittleEndian, &tamReg); err != nil {
		return models.Pokemon{}, inicioRegistro, fmt.Errorf("erro ao ler tamanho do registro: %v Linha Corrompida: %d", err, inicioRegistro)
	}

	// Lê os bytes correspondentes ao registro atual
	pokeBytes := make([]byte, tamReg-4)
	if _, err := io.ReadFull(file, pokeBytes); err != nil {
		return models.Pokemon{}, inicioRegistro, fmt.Errorf("erro ao ler registro: %v Linha Corrompida: %d", err, inicioRegistro)
	}

	// Converte os bytes para uma struct models.Pokemon se nao houver lapide
	var pokemonAtual models.Pokemon
	if lapide != 0 {
		if err := pokemonAtual.ParseBinToPoke(pokeBytes); err != nil {
			return models.Pokemon{}, inicioRegistro, fmt.Errorf("erro ao converter registro para Pokemon: %v Linha Corrompida: %d", err, inicioRegistro)
		}
	}

	return pokemonAtual, inicioRegistro, nil
}

func tamanhoProxRegistro(file *os.File, ponteiroRegistro int64) (int64, int32, error) {
	// Lê e confere a lapide do arquivo
	var lapide int32
	if err := binary.Read(file, binary.LittleEndian, &lapide); err != nil {
		return ponteiroRegistro, lapide, fmt.Errorf("erro ao ler lapide: %v Linha Corrompida: %d", err, ponteiroRegistro)
	}

	// Lê o tamanho do registro atual
	var tamReg int32
	if err := binary.Read(file, binary.LittleEndian, &tamReg); err != nil {
		return ponteiroRegistro, lapide, fmt.Errorf("erro ao ler tamanho do registro: %v Linha Corrompida: %d", err, ponteiroRegistro)
	}

	return int64(tamReg), lapide, nil
}

/* func pularRegistro(inicioRegistro int64, tamanhoRegistro int64) int64 {
	return inicioRegistro + tamanhoRegistro
} */

func NumRegistros() (int, int64, error) {
	file, err := os.Open(BIN_FILE)
	inicioRegistros := int64(0)
	if err != nil {
		return 0, 0, fmt.Errorf("erro ao abrir o arquivo: %v", err)
	}
	defer file.Close()

	// Lê o número de entradas no arquivo
	var numEntradas int32
	if err = binary.Read(file, binary.LittleEndian, &numEntradas); err != nil {
		inicioRegistros, _ := file.Seek(0, io.SeekCurrent)
		return 0, inicioRegistros, fmt.Errorf("erro ao ler número de entradas: %v Linha Corrompida: 0", err)
	}

	inicioRegistros, _ = file.Seek(0, io.SeekCurrent)
	return int(numEntradas), inicioRegistros, nil
}

func DeletarPokemon(posicao int64) error {
	file, err := os.OpenFile(BIN_FILE, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo: %v", err)
	}
	defer file.Close()

	if _, err = file.Seek(posicao, 0); err != nil {
		return fmt.Errorf("erro ao posicionar ponteiro no arquivo: %v", err)
	}

	if err = binary.Write(file, binary.LittleEndian, int32(0)); err != nil {
		return fmt.Errorf("erro ao escrever valor no arquivo: %v", err)
	}

	return nil
}

func AlterarNumRegistros(n int32) error {
	// Abrir o arquivo no modo de leitura e escrita
	file, err := os.OpenFile(BIN_FILE, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Ler o valor atual
	var numRegistros int32
	err = binary.Read(file, binary.LittleEndian, &numRegistros)
	if err != nil {
		return err
	}

	// Alterar o valor
	numRegistros += n

	// Voltar para o início do arquivo para escrever o novo valor
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	// Escrever o novo valor
	err = binary.Write(file, binary.LittleEndian, &numRegistros)
	if err != nil {
		return err
	}

	return nil
}

func AppendPokemon(pokemon []byte) error {
	file, err := os.OpenFile(BIN_FILE, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	err = binary.Write(file, binary.LittleEndian, pokemon)
	if err != nil {
		return err
	}

	AlterarNumRegistros(1)

	return err
}
