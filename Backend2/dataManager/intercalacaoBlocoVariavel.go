package dataManager

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/Bernardo46-2/AEDS-III/models"
	"github.com/Bernardo46-2/AEDS-III/utils"
)

// IntercalacaoBalanceadaVariavel executa a ordenação externa do banco de dados binario.
// A função cria arquivos temporários de tamanho especificado, realiza a ordenação externa
// em cada um deles, e finalmente intercala os arquivos até obter um arquivo ordenado final.
//
// A diferença entre a intercalação Balanceada Comum é que a variavel tenta aproveitar ao maximo
// blocos que ja estejam coincidentemente ordenados no arquivo orignial
func IntercalacaoBalanceadaVariavel() {
	// Divide o arquivo de entrada em blocos de tamanho especificado e cria os arquivos temporários
	arquivosTemp, _ := divideArquivoEmBlocosVariaveis(BIN_FILE, 8192, TMP_DIR_PATH)

	// Realiza a intercalação dos arquivos temporários até obter um arquivo ordenado final
	arquivoOrdenado := intercalaDoisEmDois(arquivosTemp)

	// Copia o arquivo ordenado para o arquivo original e paga os paths
	CopyFile(BIN_FILE, arquivoOrdenado)
	RemoveFile(arquivoOrdenado)
}

// divideArquivoEmBlocosVariaveis realiza uma ordenação externa do arquivo de entrada contendo dados da Pokedex utilizando o algoritmo de merge sort.
// O arquivo é dividido em vários arquivos menores, cada um contendo um bloco de dados de tamanho especificado, que são ordenados
// individualmente e posteriormente combinados de forma ordenada em um único arquivo de saída.
// A função coloca dentro de um mesmo arquivo blocos que ja estejam coincidentemente ordenados
func divideArquivoEmBlocosVariaveis(caminhoEntrada string, tamanhoBloco int64, dirTemp string) ([]string, error) {
	// Abrir arquivo de entrada
	file, err := os.Open(caminhoEntrada)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Ler o número total de registros
	numRegistros, _, _ := NumRegistros()
	file.Seek(4, io.SeekStart)

	// Inicializa variaveis
	arquivosTemp := []string{}
	lastPoke := models.Pokemon{}
	os.Mkdir(dirTemp, 0755)

	// Recupera os registros e salva em blocos
	for i, j := 0, 0; j < numRegistros; i++ {
		pokeSlice := []models.Pokemon{}
		tamBlocoAtual := int64(0)

		for j < numRegistros {
			// Seta o ponteiro do arquivo
			inicioRegistro, _ := file.Seek(0, io.SeekCurrent)
			ponteiroAtual := inicioRegistro
			// Pega tamanho do registro e se possui lapide
			tamanhoRegistro, lapide, _ := tamanhoProxRegistro(file, ponteiroAtual)

			// Se nao tem lapide le o registro e salva, se nao pula
			if lapide != 0 {
				// Se nao couber no bloco finaliza e da append, se nao le e adiciona ao slice atual
				if tamBlocoAtual+tamanhoRegistro > tamanhoBloco {
					file.Seek(-8, io.SeekCurrent)
					break
				} else {
					// Se for valido realiza o parse
					tamBlocoAtual += tamanhoRegistro
					pokemonAtual, _, _ := readRegistro(file, inicioRegistro)
					pokemonAtual.CalculateSize()
					pokeSlice = append(pokeSlice, pokemonAtual)
					j++
				}
			} else {
				// Caso tenha lapide faz uma leitura vazia para pular o registro
				readRegistro(file, inicioRegistro)
				j++
			}
		}

		// Ordena os elementos do bloco
		sort.Slice(pokeSlice, func(i, j int) bool {
			return pokeSlice[i].Numero < pokeSlice[j].Numero
		})

		// Caso o valor salvo referente ao ultimo bloco lido seja menor do que o primeiro valor
		// do novo bloco ordenado, a escrita sera feita no mesmo arquivo,
		// caso contrario sera feita em um novo arquivo
		if i > 0 && lastPoke.Numero < pokeSlice[0].Numero {
			// Abre o arquivo novamente
			path := fmt.Sprintf(dirTemp+"temp_%d.bin", i-1)
			fileAppendFinal, _ := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

			// Escreve os pokemons no arquivo
			for i := 0; i < len(pokeSlice); i++ {
				tmp := pokeSlice[i].ToBytes()
				binary.Write(fileAppendFinal, binary.LittleEndian, tmp)
			}

			// Recupera a quantidade de elementos no arquivo e atualiza para o novo valor
			fileAppendFinal.Seek(0, io.SeekStart)
			var novoNumRegistros int32
			binary.Read(fileAppendFinal, binary.LittleEndian, &novoNumRegistros)
			novoNumRegistros += int32(len(pokeSlice))
			fileAppendFinal.Close()
			fileAppendStart, _ := os.OpenFile(path, os.O_RDWR, 0644)
			fileAppendStart.Seek(0, io.SeekStart)
			binary.Write(fileAppendStart, binary.LittleEndian, novoNumRegistros)
			fileAppendStart.Close()
			i--
		} else {

			// Cria um novo arquivo e guarda no slice
			lastPoke = pokeSlice[len(pokeSlice)-1]
			caminhoTemp := filepath.Join(dirTemp, fmt.Sprintf("temp_%d.bin", i))
			arquivoTemp, _ := os.Create(caminhoTemp)
			arquivosTemp = append(arquivosTemp, caminhoTemp)

			// Serializa e escreve os dados no arquivo
			binary.Write(arquivoTemp, binary.LittleEndian, utils.IntToBytes(int32(len(pokeSlice))))
			for i := 0; i < len(pokeSlice); i++ {
				tmp := pokeSlice[i].ToBytes()
				binary.Write(arquivoTemp, binary.LittleEndian, tmp)
			}

			arquivoTemp.Close()
		}
	}

	return arquivosTemp, err
}
