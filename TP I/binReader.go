package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
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
func readBinToPoke(id int) (Pokemon, error) {
	// Abre o arquivo binário
	file, err := os.Open(BIN_FILE)
	if err != nil {
		return Pokemon{}, fmt.Errorf("Erro ao abrir o arquivo: %v", err)
	}
	defer file.Close()

	// Lê o número de entradas no arquivo
	var numEntradas int32
	if err = binary.Read(file, binary.LittleEndian, &numEntradas); err != nil {
		return Pokemon{}, fmt.Errorf("Erro ao ler número de entradas: %v", err)
	}

	fmt.Printf("Numero de registros = %d\n", numEntradas)

	// Percorre as entradas do arquivo
	for i := 0; i < int(numEntradas); i++ {
		// Lê o tamanho do registro atual
		var tamReg int32

		if err = binary.Read(file, binary.LittleEndian, &tamReg); err != nil {
			return Pokemon{}, fmt.Errorf("Erro ao ler tamanho do registro: %v", err)
		}

		// Lê os bytes correspondentes ao registro atual
		pokeBytes := make([]byte, tamReg-4)
		if _, err := io.ReadFull(file, pokeBytes); err != nil {
			return Pokemon{}, fmt.Errorf("Erro ao ler registro: %v", err)
		}

		// Converte os bytes para uma struct Pokemon
		var pokemonAtual Pokemon
		if err = pokemonAtual.parseBinToPoke(pokeBytes); err != nil {
			return Pokemon{}, fmt.Errorf("Erro ao converter registro para Pokemon: %v", err)
		}

		// Verifica se o número do Pokémon atual é o procurado
		if pokemonAtual.Numero == int32(id) {
			return pokemonAtual, nil
		}
	}

	// Se não encontrou o Pokémon procurado, retorna um erro
	return Pokemon{}, fmt.Errorf("Pokemon não encontrado")
}

func incrementNumRegistros() error {
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

	// Incrementar o valor
	numRegistros++

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
