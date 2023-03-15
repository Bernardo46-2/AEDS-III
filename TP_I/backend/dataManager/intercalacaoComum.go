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

// IntercalacaoBalanceadaComum executa a ordenação externa do banco de dados binario.
// A função cria arquivos temporários de tamanho especificado, realiza a ordenação externa
// em cada um deles, e finalmente intercala os arquivos até obter um arquivo ordenado final.
func IntercalacaoBalanceadaComum() {
	// Divide o arquivo de entrada em blocos de tamanho especificado e cria os arquivos temporários
	arquivosTemp, _ := divideArquivoEmBlocos(BIN_FILE, 8192, TMP_DIR_PATH)

	// Realiza a intercalação dos arquivos temporários até obter um arquivo ordenado final
	arquivoOrdenado := intercalaDoisEmDois(arquivosTemp)

	// Copia o arquivo ordenado para o arquivo original e paga os paths
	CopyFile(BIN_FILE, arquivoOrdenado)
	RemoveFile(arquivoOrdenado)
}

// divideArquivoEmBlocos realiza uma ordenação externa do arquivo de entrada contendo dados da Pokedex utilizando o algoritmo de merge sort.
// O arquivo é dividido em vários arquivos menores, cada um contendo um bloco de dados de tamanho especificado, que são ordenados
// individualmente e posteriormente combinados de forma ordenada em um único arquivo de saída.
// A função retorna um slice com os caminhos dos arquivos temporários criados.
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

	// Itera enquanto houver registro
	for i, j := 0, 0; j < numRegistros; i++ {
		// Cria o path com o caminho especificado e salva em variavel
		caminhoTemp := filepath.Join(dirTemp, fmt.Sprintf("temp_%d.bin", i))
		arquivoTemp, _ := os.Create(caminhoTemp)
		arquivosTemp = append(arquivosTemp, caminhoTemp)
		defer arquivoTemp.Close()

		// Inicializa variaveis
		pokeSlice := []models.Pokemon{}
		tamBlocoAtual := int64(0)

		// Repete enquanto existir espaço no bloco
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
				// Caso tenha lapide faz uma leitura fazia para pular o registro
				readRegistro(file, inicioRegistro)
				j++
			}
		}

		// Ordena os elementos do bloco
		sort.Slice(pokeSlice, func(i, j int) bool {
			return pokeSlice[i].Numero < pokeSlice[j].Numero
		})

		// Guarda no inicio do arquivo a quantidade de elementos que ele ira possuir
		binary.Write(arquivoTemp, binary.LittleEndian, utils.IntToBytes(int32(len(pokeSlice))))

		// Serializa e grava os registros
		for i := 0; i < len(pokeSlice); i++ {
			tmp := pokeSlice[i].ToBytes()
			binary.Write(arquivoTemp, binary.LittleEndian, tmp)
		}
	}

	return arquivosTemp, err
}

// intercalaDoisEmDois é uma função recursiva que recebe um slice contendo os caminhos
// dos arquivos a serem intercalados e realiza a intercalação dois a dois atraves do
// metodo de mergeSort.
// A função retorna o caminho do arquivo resultante da intercalação.
//
// arquivos é um slice contendo os caminhos dos arquivos a serem intercalados.
//
// Retorna o caminho do arquivo resultante da intercalação.
func intercalaDoisEmDois(arquivos []string) string {
	// Para a recursao quando existir apenas um unico arquivo
	if len(arquivos) == 1 {
		return arquivos[0]
	}

	// Intercala os subconjuntos e remove arquivos de lixo
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

// Intercala recebe dois paths de arquivos binarios e realiza a intercalação externa dos mesmos
//
// # É utilizado uma versão do algoritmo merge sort
//
// Retorna uma string contendo o nome do novo arquivo gerado com os dados intercalados
func intercala(arquivo1, arquivo2 string) (string, error) {

	// Inicializa variaveis
	var err error
	i, j := 0, 0
	pokemon1 := models.Pokemon{}
	pokemon2 := models.Pokemon{}

	// Abre os dois arquivos para leitura
	file1, err := os.Open(arquivo1)
	file2, err := os.Open(arquivo2)
	defer file1.Close()
	defer file2.Close()

	// Cria um novo arquivo temporário para escrita
	novoArquivo, err := os.Create("data/tmp/tmp.bin")
	defer novoArquivo.Close()

	// Lê a primeira linha contendo o tamanho de cada arquivo
	var tamFile1 int32
	var tamFile2 int32
	binary.Read(file1, binary.LittleEndian, &tamFile1)
	binary.Read(file2, binary.LittleEndian, &tamFile2)

	// Seta o inicio correto de leitura dos arquivos
	file1.Seek(4, io.SeekStart)
	file2.Seek(4, io.SeekStart)

	// Separa o espaço de contador de registros do novo arquivo
	binary.Write(novoArquivo, binary.LittleEndian, int32(tamFile1+tamFile2))

	// Enquanto houver linhas em ambos os arquivos, compara e escreve no novo arquivo
	for i < int(tamFile1) && j < int(tamFile2) {
		// Guarda a posição a ser manipulada
		ponteiro1, _ := file1.Seek(0, io.SeekCurrent)
		ponteiro2, _ := file2.Seek(0, io.SeekCurrent)

		// Lê os registros
		pokemon1, _, err = readRegistro(file1, ponteiro1)
		pokemon2, _, err = readRegistro(file2, ponteiro2)

		// Insere o menor ID primeiro
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
		}
		pokemon2.CalculateSize()
		binary.Write(novoArquivo, binary.LittleEndian, pokemon2.ToBytes())
		j++
	}

	// Retorna o nome do novo arquivo criado
	return novoArquivo.Name(), err
}

// CopyFile copia um arquivo do caminho de origem para o caminho de destino.
//
// destPath é o caminho do arquivo de destino.
// srcPath é o caminho do arquivo de origem.
//
// Retorna um erro, se ocorrer algum problema durante a cópia do arquivo.
func CopyFile(destPath, srcPath string) (err error) {
	err = nil

	// Abre o arquivo de origem para leitura
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Cria o arquivo de destino para escrita
	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copia o conteúdo do arquivo de origem para o arquivo de destino
	if _, err := io.Copy(destFile, srcFile); err != nil {
		return err
	}

	// Garante que todos os dados sejam gravados no arquivo de destino
	if err := destFile.Sync(); err != nil {
		return err
	}

	// Obtém as informações do arquivo de origem e aplica as mesmas permissões ao arquivo de destino
	srcFileInfo, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	return os.Chmod(destPath, srcFileInfo.Mode())
}

// PrintBin é uma função criada para depuração do codigo.
//
// Ela abre um arquivo binario e printa o ID e Nome dos pokemons existentes
func PrintBin(path string) {
	// Abre o arquivo binário
	file, _ := os.Open(path)
	defer file.Close()

	// Lê o número de entradas no arquivo
	var numEntradas int32
	binary.Read(file, binary.LittleEndian, &numEntradas)

	// Percorre as entradas do arquivo
	pokeArray := []models.Pokemon{}
	for i := 0; i < int(numEntradas); i++ {
		// Grava a localização do inicio do registro
		inicioRegistro, _ := file.Seek(0, io.SeekCurrent)
		pokemonAtual, _, _ := readRegistro(file, inicioRegistro)
		pokeArray = append(pokeArray, pokemonAtual)
	}

	// Printa o conteudo do arquivo
	for i := 0; i < len(pokeArray); i++ {
		fmt.Printf("Id = %d | Nome = %s\n", pokeArray[i].Numero, pokeArray[i].Nome)
	}
}

// RemoveFile remove o arquivo no caminho especificado.
//
// Retorna um erro, se ocorrer algum problema durante a remoção do arquivo.
func RemoveFile(filePath string) error {
	// Tenta remover o arquivo no caminho especificado
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("erro ao remover arquivo: %v", err)
	}
	return nil
}
