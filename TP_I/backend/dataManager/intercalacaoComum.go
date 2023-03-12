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

		arquivoTemp.Close()

		pokeSlice = []models.Pokemon{}
	}

	return arquivosTemp, err
}

func IntercalacaoBalanceadaComum() {
	arquivosTemp, _ := divideArquivoEmBlocos(BIN_FILE, 8192, TMP_DIR_PATH)
	arquivoOrdenado := intercalaDoisEmDois(arquivosTemp)
	CopyFile(BIN_FILE, arquivoOrdenado)
	RemoveFile(arquivoOrdenado)
}

func intercalaDoisEmDois(arquivos []string) string {
	if len(arquivos) == 1 {
		return arquivos[0]
	}

	novosArquivos := []string{}
	for i := 0; i < len(arquivos); i += 2 {
		if i+1 < len(arquivos) {
			novoArquivo, _ := intercala(arquivos[i], arquivos[i+1])
			CopyFile(arquivos[i], novoArquivo)
			RemoveFile(arquivos[i+1])
			RemoveFile(novoArquivo)
			novosArquivos = append(novosArquivos, arquivos[i])
		} else {
			// Caso ímpar, só adiciona o arquivo na lista de novos arquivos
			novosArquivos = append(novosArquivos, arquivos[i])
		}
	}
	// Faz a chamada recursiva até restar apenas um arquivo
	return intercalaDoisEmDois(novosArquivos)
}

func intercala(arquivo1, arquivo2 string) (string, error) {
	var err error

	i, j := 0, 0
	pokemon1 := models.Pokemon{}
	pokemon2 := models.Pokemon{}

	// Abre os dois arquivos para leitura
	file1, err := os.Open(arquivo1)
	if err != nil {
		fmt.Println(err.Error())
		panic(1)
	}

	defer file1.Close()

	file2, err := os.Open(arquivo2)
	if err != nil {
		fmt.Println(err.Error())
		panic(2)
	}

	defer file2.Close()

	// Cria um novo arquivo temporário para escrita
	novoArquivo, err := os.Create("data/tmp/tmp.bin")
	if err != nil {
		fmt.Println(err.Error())
		panic(3)
	}

	defer novoArquivo.Close()

	// Lê a primeira linha de cada arquivo
	var tamFile1 int32
	var tamFile2 int32
	err = binary.Read(file1, binary.LittleEndian, &tamFile1)
	if err != nil {
		fmt.Println(err.Error() + " erro no arquivo: " + arquivo1)
		panic(33)
	}
	err = binary.Read(file2, binary.LittleEndian, &tamFile2)
	if err != nil {
		fmt.Println(err.Error() + " erro no arquivo: " + arquivo2)
		panic(44)
	}

	_, err = file1.Seek(4, io.SeekStart)
	if err != nil {
		fmt.Println(err.Error())
		panic(4)
	}

	_, err = file2.Seek(4, io.SeekStart)
	if err != nil {
		fmt.Println(err.Error())
		panic(5)
	}

	binary.Write(novoArquivo, binary.LittleEndian, int32(tamFile1+tamFile2))

	// Enquanto houver linhas em ambos os arquivos, compara e escreve no novo arquivo
	for i < int(tamFile1) && j < int(tamFile2) {
		ponteiro1, _ := file1.Seek(0, io.SeekCurrent)
		ponteiro2, _ := file2.Seek(0, io.SeekCurrent)

		pokemon1, _, err = readRegistro(file1, ponteiro1)
		if err != nil {
			fmt.Println(err.Error())
			panic(6)
		}

		pokemon2, _, err = readRegistro(file2, ponteiro2)
		if err != nil {
			fmt.Println(err.Error())
			panic(7)
		}

		if pokemon1.Numero < pokemon2.Numero {
			pokemon1.CalculateSize()
			binary.Write(novoArquivo, binary.LittleEndian, pokemon1.ToBytes())
			file2.Seek(ponteiro2, io.SeekStart)
			i++
		} else {
			pokemon2.CalculateSize()
			binary.Write(novoArquivo, binary.LittleEndian, pokemon2.ToBytes())
			file1.Seek(ponteiro1, io.SeekStart)
			j++
		}
	}

	// Se houver linhas restantes em um dos arquivos, escreve no novo arquivo
	for i < int(tamFile1) {
		ponteiro1, _ := file1.Seek(0, io.SeekCurrent)
		pokemon1, _, err = readRegistro(file1, ponteiro1)
		if err != nil {
			fmt.Println(err.Error())
			panic(8)
		}
		pokemon1.CalculateSize()
		binary.Write(novoArquivo, binary.LittleEndian, pokemon1.ToBytes())
		i++
	}

	for j < int(tamFile2) {
		ponteiro2, _ := file2.Seek(0, io.SeekCurrent)
		pokemon2, _, err = readRegistro(file2, ponteiro2)
		if err != nil {
			fmt.Println(err.Error())
			panic(9)
		}
		pokemon2.CalculateSize()
		binary.Write(novoArquivo, binary.LittleEndian, pokemon2.ToBytes())
		j++
	}

	// Retorna o nome do novo arquivo criado
	return novoArquivo.Name(), err
}

func CopyFile(destPath, srcPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return err
	}

	if err := destFile.Sync(); err != nil {
		return err
	}

	srcFileInfo, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	return os.Chmod(destPath, srcFileInfo.Mode())
}

func PrintBin(path string) {
	// Abre o arquivo binário
	file, _ := os.Open(path)
	defer file.Close()

	// Lê o número de entradas no arquivo
	var numEntradas int32
	binary.Read(file, binary.LittleEndian, &numEntradas)

	pokeArray := []models.Pokemon{}
	// Percorre as entradas do arquivo
	for i := 0; i < int(numEntradas); i++ {
		// Grava a localização do inicio do registro
		inicioRegistro, _ := file.Seek(0, io.SeekCurrent)
		pokemonAtual, _, _ := readRegistro(file, inicioRegistro)
		pokeArray = append(pokeArray, pokemonAtual)
	}

	for i := 0; i < len(pokeArray); i++ {
		fmt.Printf("Id = %d | Nome = %s\n", pokeArray[i].Numero, pokeArray[i].Nome)
	}
}

func RemoveFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("erro ao remover arquivo: %v", err)
	}
	return nil
}
