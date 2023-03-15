package dataManager

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Bernardo46-2/AEDS-III/models"
	"github.com/Bernardo46-2/AEDS-III/utils"
)

type heapNode struct {
	Peso    int
	Pokemon models.Pokemon
}

func IntercalacaoPorSubstituicao() {
	// Abrir arquivo de entrada
	file, err := os.OpenFile(BIN_FILE, os.O_RDWR, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	// Ler o número total de registros
	numRegistros, _, _ := NumRegistros()
	file.Seek(4, io.SeekStart)

	// Criar os arquivos temporários
	arquivosTemp := []string{}
	pokeHeap := make([]heapNode, 7)

	for i := 0; i < 7 && i < numRegistros; {
		inicioRegistro, _ := file.Seek(0, io.SeekCurrent)
		pokemonAtual, _, _ := readRegistro(file, inicioRegistro)
		if pokemonAtual.Numero != -1 {
			pokeHeap[0] = heapNode{0, pokemonAtual}
			balanceHeap(pokeHeap, 0)
			i++
		}
	}

	caminhoTemp := filepath.Join(TMP_DIR_PATH, "temp_0.bin")
	arquivosTemp = append(arquivosTemp, caminhoTemp)
	arquivoTemp, _ := os.Create(caminhoTemp)
	inicioRegistro, _ := file.Seek(0, io.SeekCurrent)
	/* 	fmt.Println(inicioRegistro) */
	binary.Write(arquivoTemp, binary.LittleEndian, utils.IntToBytes(int32(0)))

	// last := pokeHeap[0]
	peso := 0

	for i := 0; i < (numRegistros - 7); i++ {
		// fmt.Printf("i = %d\n", i)
		inicioRegistro, _ = file.Seek(0, io.SeekCurrent)
		_, lapide, _ := tamanhoProxRegistro(file, inicioRegistro)
		file.Seek(-8, io.SeekCurrent)

		if lapide != 0 {
			// pega a cabeca do heap
			pokeTmp := pokeHeap[0].Pokemon
			pokeTmp.CalculateSize()

			// escreve no arquivo
			arquivoTemp.Seek(0, io.SeekEnd)
			binary.Write(arquivoTemp, binary.LittleEndian, pokeTmp.ToBytes())
			aumentaNumRegistros(arquivoTemp)

			// adiciona novo valor ao heap, se for menor cria novo arquivo e aumenta o peso
			pokemonAtual, _, _ := readRegistro(file, inicioRegistro)
			pokemonAtual.CalculateSize()
			if pokemonAtual.Numero < pokeHeap[0].Pokemon.Numero {
				peso++
				arquivoTemp.Close()
				caminhoTemp = filepath.Join(TMP_DIR_PATH, fmt.Sprintf("temp_%d.bin", peso))
				arquivosTemp = append(arquivosTemp, caminhoTemp)
				arquivoTemp, _ = os.Create(caminhoTemp)
				binary.Write(arquivoTemp, binary.LittleEndian, utils.IntToBytes(int32(0)))
			}
			pokeHeap[0] = heapNode{peso, pokemonAtual}

			// balanceia o novo heap
			balanceHeap(pokeHeap, 0)
		} else {
			readRegistro(file, inicioRegistro)
		}
	}

	for i := 0; i < 7; i++ {
		// pega a cabeca do heap
		pokeTmp := pokeHeap[0].Pokemon
		pokeTmp.CalculateSize()

		// escreve no arquivo
		arquivoTemp.Seek(0, io.SeekEnd)
		binary.Write(arquivoTemp, binary.LittleEndian, pokeTmp.ToBytes())
		aumentaNumRegistros(arquivoTemp)

		// adiciona novo valor ao heap, se for menor cria novo arquivo e aumenta o peso
		pokeHeap[0] = pokeHeap[6-i]
		if pokeHeap[0].Pokemon.Numero < pokeHeap[6-i].Pokemon.Numero {
			peso++
			arquivoTemp.Close()
			caminhoTemp = filepath.Join(TMP_DIR_PATH, fmt.Sprintf("temp_%d.bin", peso))
			arquivosTemp = append(arquivosTemp, caminhoTemp)
			arquivoTemp, _ = os.Create(caminhoTemp)
			binary.Write(arquivoTemp, binary.LittleEndian, utils.IntToBytes(int32(0)))
		}
		pokeHeap[0] = heapNode{peso, pokeHeap[0].Pokemon}

		// balanceia o novo heap
		pokeHeap = removeLastItem(pokeHeap)
		balanceHeap(pokeHeap, 0)
	}

	arquivoTemp.Close()

	arquivoOrdenado := intercalaDoisEmDois(apagaVazios(arquivosTemp))
	CopyFile(BIN_FILE, arquivoOrdenado)
	RemoveFile(arquivoOrdenado)
}

func balanceHeap(heap []heapNode, index int) []heapNode {
	leftIndex := 2*index + 1
	rightIndex := 2*index + 2
	smallest := index

	// Encontra o menor elemento entre o índice atual e seus filhos
	if leftIndex < len(heap) && biggerThan(heap[smallest], heap[leftIndex]) {
		smallest = leftIndex
	}
	if rightIndex < len(heap) && biggerThan(heap[smallest], heap[rightIndex]) {
		smallest = rightIndex
	}

	// Se o menor elemento não for o índice atual, troca eles de lugar e chama a função recursivamente
	if smallest != index {
		heap[index], heap[smallest] = heap[smallest], heap[index]
		heap = balanceHeap(heap, smallest)
	}

	return heap
}

func biggerThan(node1, node2 heapNode) bool {
	if node1.Peso > node2.Peso {
		return true
	} else if node1.Peso == node2.Peso && node1.Pokemon.Numero > node2.Pokemon.Numero {
		return true
	}
	return false
}

func aumentaNumRegistros(file *os.File) {
	var err error
	// Ler o valor atual
	var numRegistros int32
	file.Seek(0, 0)
	err = binary.Read(file, binary.LittleEndian, &numRegistros)
	if err != nil {
		fmt.Println(err.Error() + "1")
	}

	// Alterar o valor
	numRegistros++

	// Voltar para o início do arquivo para escrever o novo valor
	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Println(err.Error() + "2")
	}

	// Escrever o novo valor
	err = binary.Write(file, binary.LittleEndian, &numRegistros)
	if err != nil {
		fmt.Println(err.Error() + "3")
	}
}

func removeLastItem(s []heapNode) []heapNode {
	return append([]heapNode(nil), s[:len(s)-1]...)
}

func testaEDeletaArquivo(path string) bool {
	result := false
	// Obtém informações sobre o arquivo
	info, err := os.Stat(path)
	if err != nil {
		panic(1)
	}

	// Verifica se o tamanho do arquivo é zero
	if info.Size() == 0 {
		// Deleta o arquivo
		err := os.Remove(path)
		result = true
		if err != nil {
			panic(2)

		}
	}

	return result
}

func apagaVazios(paths []string) []string {
	var result []string
	for _, path := range paths {
		if testaEDeletaArquivo(path) {
			continue
		}
		result = append(result, path)
	}
	return result
}
