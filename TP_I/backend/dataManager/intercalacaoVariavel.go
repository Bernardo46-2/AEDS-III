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

func IntercalacaoBalanceadaVariavel() {
	arquivosTemp, _ := divideArquivoEmBlocosVariaveis(BIN_FILE, 8192, TMP_DIR_PATH)
	PrintBin(arquivosTemp[0])
	/* arquivoOrdenado := intercalaDoisEmDois(arquivosTemp)
	CopyFile(BIN_FILE, arquivoOrdenado)
	RemoveFile(arquivoOrdenado) */
}

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

	// Criar os arquivos temporários
	arquivosTemp := []string{}
	pokeSlice := []models.Pokemon{}
	lastPoke := models.Pokemon{}
	for i, j := 0, 0; j < numRegistros; i++ {
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

		if i > 0 {
			if lastPoke.Numero < pokeSlice[0].Numero {
				path := fmt.Sprintf(dirTemp+"temp_%d.bin", i-1)
				fileAppendFinal, _ := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

				for i := 0; i < len(pokeSlice); i++ {
					tmp := pokeSlice[i].ToBytes()
					binary.Write(fileAppendFinal, binary.LittleEndian, tmp)
				}

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
			}
		} else {
			lastPoke = pokeSlice[len(pokeSlice)-1]

			caminhoTemp := filepath.Join(dirTemp, fmt.Sprintf("temp_%d.bin", i))
			arquivoTemp, _ := os.Create(caminhoTemp)
			arquivosTemp = append(arquivosTemp, caminhoTemp)

			binary.Write(arquivoTemp, binary.LittleEndian, utils.IntToBytes(int32(len(pokeSlice))))

			for i := 0; i < len(pokeSlice); i++ {
				tmp := pokeSlice[i].ToBytes()
				binary.Write(arquivoTemp, binary.LittleEndian, tmp)
			}

			arquivoTemp.Close()
		}

		pokeSlice = []models.Pokemon{}
	}

	return arquivosTemp, err
}

/* func AppendPokemon(pokemon []byte) error {
/* func AppendPokemon(pokemon []byte) error {
	file, err := os.OpenFile(BIN_FILE, os.O_WRONLY|os.O_APPEND, 0644
	if err != nil.OpenFile(BIN_FILE, os.O_WRONLY|os.O_APPEND, 0644
	if err != il
		eturn er
	}
defer file.Close(

	err = binary.ite(file, binary.LittleEndian, pokemon
	if err != il
return er


AlterarNumRegistros(1

	retrn er
} */
