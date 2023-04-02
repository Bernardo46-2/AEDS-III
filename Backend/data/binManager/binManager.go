// O arquivo binManager do pacote dataManager realiza o tratamento do arquivo binario
// Recebe as requisiçoes a partir do pacote service e termina recuperando ou editando
// os registros binarios necessarios.
package binManager

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/Bernardo46-2/AEDS-III/models"
)

// const FILE string = "data/files/pokedex2.csv"
// const BIN_FILE string = "data/files/pokedex.dat"

// ReadBinToPoke lê um arquivo binário com informações de Pokémons e retorna
// um Pokémon com o número especificado. Caso o número não seja encontrado,
// o Pokémon retornado terá seu número igual a -1.
//
// O arquivo binario esta estruturado em:
// Lapide (int32)
// Tamanho (int32)
// Registro
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

		pokemonAtual, inicioRegistro, _ := ReadRegistro(file, inicioRegistro)

		// Verifica se o número do Pokémon atual é o procurado
		if pokemonAtual.Numero == int32(id) {
			return pokemonAtual, inicioRegistro, nil
		}
	}

	// Se não encontrou o Pokémon procurado, retorna um erro
	pos, _ := file.Seek(0, io.SeekCurrent)
	return models.Pokemon{}, pos, fmt.Errorf("pokemon não encontrado")
}

// ReadRegistro recebe um arquivo e o ponto de onde a leitura deve ser iniciada
//
// Em caso de lapide retorna um objeto pokemon com id -1
// Em caso de erro gera uma mensagem formatada com o tipo e a linha corrompida
func ReadRegistro(file *os.File, inicioRegistro int64) (pokemonAtual models.Pokemon, pos int64, err error) {
	// Seta a leitura para a posição determinada
	pos, err = file.Seek(inicioRegistro, io.SeekStart)
	pokemonAtual = models.Pokemon{Numero: -1}

	// Lê e confere a lapide do arquivo
	var lapide int32
	if err := binary.Read(file, binary.LittleEndian, &lapide); err != nil {
		return pokemonAtual, pos, fmt.Errorf("erro ao ler lapide: %v Linha Corrompida: %d", err, pos)
	}

	// Lê o tamanho do registro atual
	var tamReg int32
	if err := binary.Read(file, binary.LittleEndian, &tamReg); err != nil {
		return pokemonAtual, pos, fmt.Errorf("erro ao ler tamanho do registro: %v Linha Corrompida: %d", err, pos)
	}

	// Lê os bytes correspondentes ao registro atual
	pokeBytes := make([]byte, tamReg-4)
	if _, err := io.ReadFull(file, pokeBytes); err != nil {
		return pokemonAtual, pos, fmt.Errorf("erro ao ler registro: %v Linha Corrompida: %d", err, pos)
	}

	// Converte os bytes para uma struct models.Pokemon se nao houver lapide
	if lapide != 1 {
		if err := pokemonAtual.ParseBinToPoke(pokeBytes); err != nil {
			return pokemonAtual, pos, fmt.Errorf("erro ao converter registro para Pokemon: %v Linha Corrompida: %d", err, pos)
		}
	}

	return
}

// TamanhoProxRegistro recebe um arquivo e uma posição de leitura e retorna
// o tamanho do registro a ser lido e se possivelmente possui lapide
func TamanhoProxRegistro(file *os.File, ponteiroRegistro int64) (int64, int32, error) {
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

// NumRegistros abre o arquivo binario e analisa o marcador de quantidade de registros
// que esta na posição 0 do arquivo
func NumRegistros() (numEntradas int, inicioRegistros int64, err error) {
	// Abre o arquivo para leitura
	file, err := os.Open(BIN_FILE)
	if err != nil {
		return 0, 0, fmt.Errorf("erro ao abrir o arquivo: %v", err)
	}
	defer file.Close()

	// Lê o número de entradas no arquivo
	var numEntradas32 int32
	if err = binary.Read(file, binary.LittleEndian, &numEntradas32); err != nil {
		inicioRegistros, _ := file.Seek(0, io.SeekCurrent)
		return 0, inicioRegistros, fmt.Errorf("erro ao ler número de entradas: %v Linha Corrompida: 0", err)
	}

	// Recupera a posição inicial do arquivo a partir do registrador
	inicioRegistros, _ = file.Seek(0, io.SeekCurrent)
	numEntradas = int(numEntradas32)
	return
}

func GetLastPokemon() (lastID int32) {
	atualPokemon := models.Pokemon{Numero: -1}
	ultimoPokemon := models.Pokemon{Numero: -1}

	file, err := os.Open(BIN_FILE)
	if err != nil {
		return -1
	}
	defer file.Close()

	// Lê o número de entradas no arquivo
	var numEntradas32 int32
	binary.Read(file, binary.LittleEndian, &numEntradas32)
	numEntradas := int(numEntradas32)

	// pula registros ate chegar ao ultimo
	for i := 0; i < numEntradas; i++ {
		var lapide int32
		binary.Read(file, binary.LittleEndian, &lapide)
		var tamReg int32
		binary.Read(file, binary.LittleEndian, &tamReg)
		pokeBytes := make([]byte, tamReg-4)
		io.ReadFull(file, pokeBytes)
		if lapide != 1 {
			atualPokemon.ParseBinToPoke(pokeBytes)
			if atualPokemon.Numero > ultimoPokemon.Numero {
				ultimoPokemon = atualPokemon
			}
		}
	}

	return ultimoPokemon.Numero
}

// DeletarPokemon recebe a posição da lapide a ser alterada no arquivo
//
// A lapide se localiza como primeira variavel do registro
// Lapide / tamanho registro / registro
func DeletarPokemon(posicao int64) error {
	// Abre o arquivo para leitura e edição
	file, err := os.OpenFile(BIN_FILE, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo: %v", err)
	}
	defer file.Close()

	// Seta a posição da lapide do registro
	if _, err = file.Seek(posicao, io.SeekStart); err != nil {
		return fmt.Errorf("erro ao posicionar ponteiro no arquivo: %v", err)
	}

	// Escreve o valor da lapide
	if err = binary.Write(file, binary.LittleEndian, int32(1)); err != nil {
		return fmt.Errorf("erro ao escrever valor no arquivo: %v", err)
	}

	return nil
}

// AlterarNumRegistros recebe uma marcação de atualização no numero de registros
//
// Adiciona ou subtrai do registro do arquivo o valor do parametro
func AlterarNumRegistros(n int32) (err error) {
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

	return
}

// AppendPokemon recebe um pokemon serializado em array de bytes e faz o append
// deste pokemon no final do arquivo
//
// Por fim atualiza o numero de registros em +1
func AppendPokemon(pokemon []byte) (address int64, err error) {

	// Abre o arquivo para leitura e append
	file, err := os.OpenFile(BIN_FILE, os.O_WRONLY|os.O_APPEND, 0644)
	address, _ = file.Seek(0, io.SeekEnd)
	if err != nil {
		return
	}
	defer file.Close()

	// Tenta fazer a escrita
	err = binary.Write(file, binary.LittleEndian, pokemon)
	if err != nil {
		return
	}

	// Atualiza a quantidade de registros
	AlterarNumRegistros(1)

	return
}

type ControleLeitura struct {
	Arquivo        *os.File  // ponteiro para o arquivo de registros
	TotalRegistros int32     // número total de registros no arquivo
	RegistrosLidos int32     // número de registros já lidos
	RegistroAtual  *Registro // ponteiro para o registro atual sendo lido
}

type Registro struct {
	Lapide   int32
	Tamanho  int32
	Pokemon  models.Pokemon
	Endereco int64
}

func (c *ControleLeitura) Close() error {
	return c.Arquivo.Close()
}

func InicializarControleLeitura(nomeArquivo string) (*ControleLeitura, error) {
	arquivo, err := os.Open(nomeArquivo)
	if err != nil {
		return nil, err
	}

	// ler o número total de registros do arquivo
	var totalRegistros int32
	err = binary.Read(arquivo, binary.LittleEndian, &totalRegistros)
	if err != nil {
		return nil, err
	}

	// criar uma instância do ControleLeitura
	controle := &ControleLeitura{
		Arquivo:        arquivo,
		RegistrosLidos: 0,
		TotalRegistros: totalRegistros,
		RegistroAtual:  nil,
	}

	return controle, nil
}

func (c *ControleLeitura) ReadTarget(targetPos int64) models.Pokemon {
	targetPokemon := models.Pokemon{Numero: -1}
	limiteArquivo, _ := c.Arquivo.Seek(0, io.SeekEnd)
	if targetPos < 0 || targetPos >= limiteArquivo {
		return targetPokemon
	}

	c.Arquivo.Seek(targetPos, io.SeekStart)
	var lapide int32
	var tamanho int32
	var conteudo models.Pokemon

	binary.Read(c.Arquivo, binary.LittleEndian, &lapide)
	binary.Read(c.Arquivo, binary.LittleEndian, &tamanho)
	conteudoBytes := make([]byte, tamanho-4)
	binary.Read(c.Arquivo, binary.LittleEndian, &conteudoBytes)

	// Converte os bytes para uma struct models.Pokemon se nao houver lapide
	if lapide != 1 {
		conteudo.ParseBinToPoke(conteudoBytes)
	} else {
		conteudo.Numero = -1
	}

	return conteudo
}

func (c *ControleLeitura) ReadNext() error {
	// verificar se todos os registros já foram lidos
	if c.RegistrosLidos >= c.TotalRegistros {
		return io.EOF // fim do arquivo
	}

	// ler os dados do registro do arquivo
	endereco, _ := c.Arquivo.Seek(0, io.SeekCurrent)
	var lapide int32
	var tamanho int32
	var conteudo models.Pokemon

	binary.Read(c.Arquivo, binary.LittleEndian, &lapide)
	binary.Read(c.Arquivo, binary.LittleEndian, &tamanho)
	conteudoBytes := make([]byte, tamanho-4)
	binary.Read(c.Arquivo, binary.LittleEndian, &conteudoBytes)

	// Converte os bytes para uma struct models.Pokemon se nao houver lapide
	if lapide != 1 {
		conteudo.ParseBinToPoke(conteudoBytes)
	} else {
		conteudo.Numero = -1
	}

	// atualizar o registro atual e o número de registros lidos
	registro := &Registro{
		Lapide:   lapide,
		Tamanho:  tamanho,
		Pokemon:  conteudo,
		Endereco: endereco,
	}
	c.RegistroAtual = registro
	c.RegistrosLidos++

	return nil
}
