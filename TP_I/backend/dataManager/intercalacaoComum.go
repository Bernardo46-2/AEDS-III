package dataManager

import (
	"fmt"
	"io"
	"os"

	"github.com/Bernardo46-2/AEDS-III/models"
)

const TMP_DIR_PATH string = "data/tmp/"

func IntercalacaoBalanceadaComum() {
	strings, err := divideArquivoEmBlocos(BIN_FILE, 4096, TMP_DIR_PATH)
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
	numRegistros, inicioRegistro, _ := NumRegistros()
	inicioRegistro, _ = file.Seek(inicioRegistro, io.SeekStart)
	fmt.Printf("inicioRegistro = %d\n", inicioRegistro)

	// Criar os arquivos temporários
	/* arquivosTemp := []string{} */

	slicePokemons := [][]models.Pokemon{}
	for i := 0; i < 30 && i < numRegistros; i++ {
		tamBlocoAtual := int64(0)
		slicePokemons = append(slicePokemons, []models.Pokemon{})
		for continuar := true; continuar; {
			inicioRegistro, _ = file.Seek(0, io.SeekCurrent)
			ponteiroAtual := inicioRegistro

			// Pega tamanho do registro e se possui lapide
			tamanhoRegistro, lapide, _ := tamanhoProxRegistro(file, ponteiroAtual)

			// Se nao tem lapide le o registro e salva, se nao pula
			if lapide != 0 {
				// Se nao couber no bloco finaliza e da append, se nao le e adiciona ao slice atual
				if tamBlocoAtual+tamanhoRegistro > tamanhoBloco {
					file.Seek(-8, io.SeekCurrent)
					continuar = false
				} else {
					tamBlocoAtual += tamanhoRegistro
					pokemonAtual, _, _ := readRegistro(file, inicioRegistro)
					slicePokemons[i] = append(slicePokemons[i], pokemonAtual)
				}
			} else {
				fmt.Printf("TA COM LAPIDE CACETE!\n")
			}
		}

	}

	fmt.Printf("Numero de blocos = %d\n", len(slicePokemons))
	for i := 0; i < len(slicePokemons); i++ {
		fmt.Printf("[%d] = %d\n", i, len(slicePokemons[i]))
	}

	return nil, nil
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
