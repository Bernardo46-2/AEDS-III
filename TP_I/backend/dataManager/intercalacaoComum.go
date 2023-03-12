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

const TMP_DIR_PATH string = "data/tmp/"

func IntercalacaoBalanceadaComum() {
	strings, err := divideArquivoEmBlocos(BIN_FILE, 8192, TMP_DIR_PATH)
	for i := 0; i < len(strings); i++ {
		fmt.Println(strings[i])
	}
	if err != nil {
		fmt.Println(err.Error())
	}
}

func divideArquivoEmBlocos(caminhoEntrada string, tamanhoBloco int64, dirTemp string) ([]string, error) {
	// Abrir arquivo de entrada
	file, err := os.Open(caminhoEntrada)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Ler o número total de registros
	numRegistros, _, _ := NumRegistros()
	file.Seek(4, io.SeekStart)

	// Criar os arquivos temporários
	arquivosTemp := []string{}
	pokeSlice := []models.Pokemon{}

	for i, j := 0, 0; j < numRegistros; i++ {
		caminhoTemp := filepath.Join(dirTemp, fmt.Sprintf("temp_%d.bin", i))
		arquivoTemp, _ := os.Create(caminhoTemp)
		arquivosTemp = append(arquivosTemp, caminhoTemp)

		tamBlocoAtual := int64(0)
		for j < numRegistros {
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
					tamBlocoAtual += tamanhoRegistro
					pokemonAtual, _, _ := readRegistro(file, inicioRegistro)
					pokemonAtual.CalculateSize()
					pokeSlice = append(pokeSlice, pokemonAtual)
					j++
				}
			} else {
				readRegistro(file, inicioRegistro)
				j++
			}
		}

		sort.Slice(pokeSlice, func(i, j int) bool {
			return pokeSlice[i].Numero < pokeSlice[j].Numero
		})

		binary.Write(arquivoTemp, binary.LittleEndian, utils.IntToBytes(int32(len(pokeSlice))))

		for i := 0; i < len(pokeSlice); i++ {
			tmp := pokeSlice[i].ToBytes()
			binary.Write(arquivoTemp, binary.LittleEndian, tmp)
		}

		pokeSlice = []models.Pokemon{}
	}

	return arquivosTemp, err
}

// Função auxiliar para escrever um bloco de tamanho fixo em um arquivo
/* func writeBloco(dest io.Writer, tamanho int64, src io.Reader) error {
	var buf [8192]byte
	for tamanho > 0 {
		n, err := io.CopyN(dest, src, int64(len(buf)))
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		tamanho -= n
	}
	return nil
} */
