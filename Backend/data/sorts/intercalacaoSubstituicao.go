// intercalacaoSubstituicao implementa o método de intercalação por substituição (heap)
// para a ordenação de arquivos grandes.
//
// Este pacote manipula arquivos grandes que não podem ser carregados na memória principal,
// utilizando uma estrutura de dados de heap para manter a ordenação dos registros durante
// as operações de intercalação. Os registros são lidos na memória, organizados em um heap
// mínimo para garantir a ordenação e escritos de volta em arquivos temporários menores,
// que são então mesclados para produzir a versão ordenada do arquivo original.
package sorts

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Bernardo46-2/AEDS-III/data/binManager"
	"github.com/Bernardo46-2/AEDS-III/models"
	"github.com/Bernardo46-2/AEDS-III/utils"
)

// heapNode faz a adaptação de pokemons para ser utilizado junto de uma variavel de peso
// para ser utilizado em um heap
type heapNode struct {
	Peso    int
	Pokemon models.Pokemon
}

// IntercalacaoPorSubstituicao lê um arquivo binário e o ordena utilizando o algoritmo HeapSort. O processo de ordenação
// é realizado em blocos de tamanho fixo 7, armazenados em um heap. Quando não for possível adicionar mais valores ao heap,
// um novo arquivo temporário é criado e os valores armazenados no heap são escritos nele, já ordenados. Ao final da
// leitura do arquivo, é realizado um merge entre os arquivos temporários para obter um único arquivo ordenado.
//
// INSTAVEL: algoritmo nao esta funcionando corretamente
// TODO: corrigir e descobrir onde esta o erro
func IntercalacaoPorSubstituicao() {
	// Abrir arquivo de entrada
	file, err := os.OpenFile(BIN_FILE, os.O_RDWR, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	// Ler o número total de registros
	numRegistros, _, _ := binManager.NumRegistros()
	file.Seek(4, io.SeekStart)

	// Criar os arquivos temporários
	arquivosTemp := []string{}
	pokeHeap := make([]heapNode, 7)
	os.Mkdir(TMP_DIR_PATH, 0755)

	var lidos int
	for i := 0; i < 7 && i < numRegistros; lidos++ {
		inicioRegistro, _ := file.Seek(0, io.SeekCurrent)
		pokemonAtual, _, _ := binManager.ReadRegistro(file, inicioRegistro)
		if pokemonAtual.Numero != -1 {
			pokeHeap[0] = heapNode{0, pokemonAtual}
			balanceHeap(pokeHeap, 0)
			i++
		}
	}

	// Cria as variaveis de path
	caminhoTemp := filepath.Join(TMP_DIR_PATH, "temp_0.bin")
	arquivosTemp = append(arquivosTemp, caminhoTemp)
	arquivoTemp, _ := os.Create(caminhoTemp)

	// Reserva o espaço de contagem de registros
	binary.Write(arquivoTemp, binary.LittleEndian, utils.IntToBytes(int32(0)))

	peso := 0
	for i := 0; i < (numRegistros - lidos); i++ {
		// Guarda a posicao de inicio do registro e verifica sua lapide
		inicioRegistro, _ := file.Seek(0, io.SeekCurrent)
		_, lapide, _ := binManager.TamanhoProxRegistro(file, inicioRegistro)
		file.Seek(-8, io.SeekCurrent)

		// Se nao possuir lapide lê o registro e adiciona ao heap,
		// depois retira do tipo do heap e adiciona ao arquivo
		if lapide == 0 {
			// pega a cabeca do heap
			pokeTmp := pokeHeap[0].Pokemon
			pokeTmp.CalculateSize()

			// escreve no arquivo
			arquivoTemp.Seek(0, io.SeekEnd)
			binary.Write(arquivoTemp, binary.LittleEndian, pokeTmp.ToBytes())
			aumentaNumRegistros(arquivoTemp)

			// adiciona novo valor ao heap, se for menor cria novo arquivo e aumenta o peso
			pokemonAtual, _, _ := binManager.ReadRegistro(file, inicioRegistro)
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
			// Realiza uma leitura vazia para descartar o valor
			binManager.ReadRegistro(file, inicioRegistro)
		}
	}

	// Esvaziar valores restantes do heap
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

// balanceHeap recebe um heap e um index e o retorna balanceado.
// Realiza o balanceamento de maneira recursiva e por isso é necessario
// fazer a chamada utilizando index = 0
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

// biggerThan Testa a hierarquia de valores do struct heapNode.
// Onde Peso vem primeiro, Numero depois
func biggerThan(node1, node2 heapNode) bool {
	if node1.Peso > node2.Peso {
		return true
	} else if node1.Peso == node2.Peso && node1.Pokemon.Numero > node2.Pokemon.Numero {
		return true
	}
	return false
}

// aumentaNumRegistros recebe um arquivo e aumenta o numero de registros.
// A função é um complemento da função ja existente no binManager
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

// removeLastItem recebe o heap e o retorna sem o ultimo elemento
func removeLastItem(s []heapNode) []heapNode {
	return append([]heapNode(nil), s[:len(s)-1]...)
}

// testaEDeletaArquivo testa a existencia de um arquivo,
// caso a mesma seja valida o arquivo sera deletado
func testaEDeletaArquivo(path string) bool {
	result := false
	// Obtém informações sobre o arquivo
	info, _ := os.Stat(path)

	// Verifica se o tamanho do arquivo é zero
	if info.Size() == 0 {
		// Deleta o arquivo
		os.Remove(path)
		result = true
	}

	return result
}

// apagaVazios recebe um array de strings e retorna sem os valores vazios
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
