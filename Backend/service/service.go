// O pacote service realiza a conversa entre as requisiçoes e o DataManager
// recebendo dados ja em formato struct e fazendo as devidas chamadas de ediçao
// no arquivo binario
package service

import (
	"errors"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Bernardo46-2/AEDS-III/data/binManager"
	"github.com/Bernardo46-2/AEDS-III/data/compress/huffman"
	"github.com/Bernardo46-2/AEDS-III/data/compress/lzw"
	aescbc "github.com/Bernardo46-2/AEDS-III/data/crypto/aes_cbc"
	"github.com/Bernardo46-2/AEDS-III/data/crypto/trivium"
	"github.com/Bernardo46-2/AEDS-III/data/indexes/bplustree"
	"github.com/Bernardo46-2/AEDS-III/data/indexes/btree"
	"github.com/Bernardo46-2/AEDS-III/data/indexes/hashing"
	"github.com/Bernardo46-2/AEDS-III/data/indexes/invertedIndex"
	"github.com/Bernardo46-2/AEDS-III/data/patternMatching/kmp"
	"github.com/Bernardo46-2/AEDS-III/data/patternMatching/rabinKarp"
	"github.com/Bernardo46-2/AEDS-III/models"
	"github.com/Bernardo46-2/AEDS-III/utils"
)

type SearchRequest struct {
	Nome         string `json:"nome"`
	JapName      string `json:"japName"`
	Especie      string `json:"especie"`
	Tipo         string `json:"tipo"`
	Descricao    string `json:"descricao"`
	IDI          string `json:"idI"`
	IDF          string `json:"idF"`
	GeracaoI     string `json:"geracaoI"`
	GeracaoF     string `json:"geracaoF"`
	LancamentoI  string `json:"lancamentoI"`
	LancamentoF  string `json:"lancamentoF"`
	AtkI         string `json:"atkI"`
	AtkF         string `json:"atkF"`
	DefI         string `json:"defI"`
	DefF         string `json:"defF"`
	HpI          string `json:"hpI"`
	HpF          string `json:"hpF"`
	AlturaI      string `json:"alturaI"`
	AlturaF      string `json:"alturaF"`
	PesoI        string `json:"pesoI"`
	PesoF        string `json:"pesoF"`
	Lendario     string `json:"lendario"`
	Mitico       string `json:"mitico"`
	PatternMatch string `json:"patternMatch"`
}

// ReadPagesNumber retorna o numero de paginas disponiveis para a
// exibiçao dos pokemons na tela inicial do site, como um menu
// de navegação entre paginas
func ReadPagesNumber() (numeroPaginas int, err error) {
	// Recuperação do numero de registros totais
	numRegistros, _, _ := binManager.NumRegistros()

	// calcula e retorna o total
	numeroPaginas = int(math.Ceil((float64(numRegistros) / float64(60))))
	return
}

func GetIdList() (ids []int32, err error) {
	c, _ := binManager.InicializarControleLeitura(binManager.BIN_FILE)
	defer c.Close()

	for {
		err = c.ReadNext()
		if err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
				break
			}
		}
		if !c.RegistroAtual.IsDead() {
			ids = append(ids, c.RegistroAtual.Pokemon.Numero)
		}
	}

	utils.InsertionSort(ids)

	return
}

func GetList(idList []int64, method int) (pokeList []models.Pokemon, duration int64, err error) {
	c, _ := binManager.InicializarControleLeitura(binManager.BIN_FILE)
	defer c.Close()

	start := time.Now()
	switch method {
	case 0: // Sequencial
		for _, id := range idList {
			pokemon, _, _ := binManager.ReadBinToPoke(int(id))
			pokeList = append(pokeList, pokemon)
		}
	case 1, -1: // Hash
		hash, _ := hashing.Load(binManager.FILES_PATH, "hashIndex")
		for _, id := range idList {
			pos, err := hash.Read(id)
			if err == nil {
				pokeList = append(pokeList, c.ReadTarget(pos))
			}
		}
		hash.Close()
	case 2: // Arvore B
		btree, _ := btree.ReadBTree(binManager.FILES_PATH)
		for _, id := range idList {
			pos := btree.Find(id)
			if pos != nil {
				pokeList = append(pokeList, c.ReadTarget(pos.Ptr))
			}
		}
		btree.Close()
	case 3: // Arvore B+
		bptreeeee, _ := bplustree.ReadBPlusTree(binManager.FILES_PATH, "id")
		for _, id := range idList {
			pos := bptreeeee.Find(float64(id))
			if err == nil {
				pokeList = append(pokeList, c.ReadTarget(pos.Ptr))
			}
		}
	}
	duration = time.Since(start).Milliseconds()

	return
}

// Create adiciona um novo pokemon ao banco de dados.
//
// Recebe um modelo pokemon e serializa para inserir
// Por fim retorna o ID do pokemon criado e erro se houver.
//
// tambem realiza: HashCreate
func Create(pokemon models.Pokemon) (int, error) {
	// Recupera o ultimo ID para gerar o proximo
	ultimoID := binManager.GetLastPokemon()
	ultimoID++
	pokemon.Numero = ultimoID

	// Prepara, serializa e insere
	pokemon.CalculateSize()
	pokeBytes := pokemon.ToBytes()
	address, err := binManager.AppendPokemon(pokeBytes)

	// Indice invertido
	invertedIndex.Create(pokemon, binManager.FILES_PATH, models.PokeStrings()...)

	// Tabela Hash
	hashing.HashCreate(int64(pokemon.Numero), address, binManager.FILES_PATH, "hashIndex")

	// Arvore B
	bTree, _ := btree.ReadBTree(binManager.FILES_PATH)
	bTree.Insert(&btree.Key{Id: int64(pokemon.Numero), Ptr: address})
	bTree.Close()

	// Arvore B+
	bplustree.Create(pokemon, address, binManager.FILES_PATH, []string{"id"})
	bplustree.Create(pokemon, int64(pokemon.Numero), binManager.FILES_PATH, models.PokeNumbers())

	return int(ultimoID), err
}

// Read recebe o ID de um pokemon, procura no banco de dados atraves do
// indice hash e o retorna, se nao achar gera um erro
func Read(id int) (models.Pokemon, error) {
	pos, err := hashing.HashRead(int64(id), binManager.FILES_PATH, "hashIndex")
	pokemon := binManager.ReadTargetPokemon(pos)
	return pokemon, err
}

// Update atualiza um registro no arquivo binário de acordo com o número do pokemon informado.
// Recebe uma struct do tipo models.Pokemon a ser atualizada.
// Retorna um erro caso ocorra algum problema ao atualizar o registro.
//
// O update é feito deletando um valor e adicionando outro ao final do arquivo.
//
// tambem realiza: HashUpdate
func Update(pokemon models.Pokemon) (err error) {
	// Recupera a posição do id no arquivo
	pos, err := hashing.HashRead(int64(pokemon.Numero), binManager.FILES_PATH, "hashIndex")
	if err != nil {
		return
	}
	old := binManager.ReadTargetPokemon(pos)

	// Serializa os dados
	pokemon.CalculateSize()
	pokeBytes := pokemon.ToBytes()

	// Deleta o antigo e insere o novo registro
	err = binManager.DeletarPokemon(pos)
	if err != nil {
		return
	}

	newAddress, err := binManager.AppendPokemon(pokeBytes)
	if err != nil {
		return
	}

	// Indice invertido
	invertedIndex.Update(pokemon, binManager.FILES_PATH, models.PokeStrings()...)

	// Tabela Hash
	err = hashing.HashUpdate(int64(pokemon.Numero), newAddress, binManager.FILES_PATH, "hashIndex")

	// Arvore B
	btree, _ := btree.ReadBTree(binManager.FILES_PATH)
	btree.Update(int64(pokemon.Numero), newAddress)
	btree.Close()

	// Arvore B+
	bplustree.Update(old, pokemon, pos, newAddress, binManager.FILES_PATH, []string{"id"})
	bplustree.Update(old, pokemon, int64(old.Numero), int64(pokemon.Numero), binManager.FILES_PATH, models.PokeNumbers())

	return
}

// Delete recebe um ID, procura no arquivo e gera a remoçao logica do mesmo
//
// tambem realiza: HashDelete
func Delete(id int) (pokemon models.Pokemon, err error) {
	// Tenta encontrar a posiçao do pokemon no arquivo binario
	var pos int64
	pos, err = hashing.HashRead(int64(id), binManager.FILES_PATH, "hashIndex")
	pokemon = binManager.ReadTargetPokemon(pos)
	if err != nil {
		return
	}

	// Efetiva a remoção logica
	if err = binManager.DeletarPokemon(pos); err != nil {
		return
	}

	// Indice invertido
	invertedIndex.Delete(pokemon, binManager.FILES_PATH, models.PokeStrings()...)

	// Tabela Hash
	hashing.HashDelete(int64(pokemon.Numero), binManager.FILES_PATH, "hashIndex")

	// Arvore B
	btree, err := btree.ReadBTree(binManager.FILES_PATH)
	if err != nil {
		return
	}
	btree.Remove(int64(id))
	btree.Close()

	// Arvore B+
	bplustree.Delete(pokemon, pos, binManager.FILES_PATH, []string{"id"})
	bplustree.Delete(pokemon, int64(pokemon.Numero), binManager.FILES_PATH, models.PokeNumbers())

	return
}

func MergeSearch(req SearchRequest) (idList []int64, duration int64, err error) {
	getFieldScDoc := func(field, text string) []invertedIndex.ScoredDocument {
		switch req.PatternMatch {
		case "1": // KMP
			return kmp.SearchPokemon(text, field)
		case "2": // Rabin Karp
			return rabinKarp.SearchPokemon(text, field)
		default:
			return invertedIndex.Read(binManager.FILES_PATH, field, strings.Fields(text)...)
		}
	}

	getIdsBPTree := func(start string, end string, field string) []invertedIndex.ScoredDocument {
		tree, _ := bplustree.ReadBPlusTree(binManager.FILES_PATH, field)
		defer tree.Close()
		startf64, _ := strconv.ParseFloat(start, 64)
		endf64, _ := strconv.ParseFloat(end, 64)
		result, _ := tree.FindRange(startf64, endf64)
		docs := make([]invertedIndex.ScoredDocument, len(result))
		for i, id := range result {
			docs[i] = invertedIndex.ScoredDocument{DocumentID: id, Score: 1}
		}
		return docs
	}

	hash, _ := hashing.Load(binManager.FILES_PATH, "hashIndex")
	defer hash.Close()

	req.LancamentoI = utils.FormatDate(req.LancamentoI)
	req.LancamentoF = utils.FormatDate(req.LancamentoF)
	req.JapName = utils.ToKatakana(req.JapName)

	start := time.Now()
	nomeScDoc := getFieldScDoc("nome", req.Nome)
	especieScDoc := getFieldScDoc("especie", req.Especie)
	tipoScDoc := getFieldScDoc("tipo", req.Tipo)
	descricaoScDoc := getFieldScDoc("descricao", req.Descricao)
	japNameScDoc := getFieldScDoc("nomeJap", req.JapName)
	duration = time.Since(start).Milliseconds()

	ID := getIdsBPTree(req.IDI, req.IDF, "numero")
	Geracao := getIdsBPTree(req.GeracaoI, req.GeracaoF, "geracao")
	Lancamento := getIdsBPTree(req.LancamentoI, req.LancamentoF, "lancamento")
	Atk := getIdsBPTree(req.AtkI, req.AtkF, "atk")
	Def := getIdsBPTree(req.DefI, req.DefF, "def")
	Hp := getIdsBPTree(req.HpI, req.HpF, "hp")
	Altura := getIdsBPTree(req.AlturaI, req.AlturaF, "altura")
	Peso := getIdsBPTree(req.PesoI, req.PesoF, "peso")

	var Lendario []invertedIndex.ScoredDocument
	var Mitico []invertedIndex.ScoredDocument
	if req.Lendario == "1" {
		Lendario = getIdsBPTree("1", "2", "lendario")
	}
	if req.Mitico == "1" {
		Mitico = getIdsBPTree("1", "2", "mitico")
	}

	scDoc := invertedIndex.Merge(nomeScDoc, especieScDoc, tipoScDoc, descricaoScDoc, japNameScDoc, ID, Geracao, Lancamento, Atk, Def, Hp, Altura, Peso, Lendario, Mitico)

	for _, tmp := range scDoc {
		idList = append(idList, tmp.DocumentID)
	}

	return
}

func Encrypt(method int) (key string) {
	utils.Create_verifier()

	aes := func(k aescbc.Key, file string) {
		iv, _ := aescbc.RandBytes(aescbc.BLOCK_SIZE)
		data, _ := os.ReadFile(file)
		data = aescbc.Encrypt(k, iv, data)
		os.WriteFile(file, data, 0644)
	}

	switch method {
	case 0:
		fallthrough
	case 1:
		t := trivium.New()
		t.Encrypt(utils.VERIFIER, utils.VERIFIER)

		t2 := trivium.New(t.Key)
		t2.Encrypt(binManager.BIN_FILE, binManager.BIN_FILE)

		key = utils.ByteArrayToAscii(t.Key)
	case 2:
		k, _ := aescbc.NewKey(128)
		aes(k, utils.VERIFIER)
		aes(k, binManager.BIN_FILE)
		key = utils.SliceToAscii(k.Key)
	case 3:
		k, _ := aescbc.NewKey(192)
		aes(k, utils.VERIFIER)
		aes(k, binManager.BIN_FILE)
		key = utils.SliceToAscii(k.Key)
	case 4:
		k, _ := aescbc.NewKey(256)
		aes(k, utils.VERIFIER)
		aes(k, binManager.BIN_FILE)
		key = utils.SliceToAscii(k.Key)
	}

	return
}

func Decrypt(method int, key string) (success bool) {
	switch method {
	case 0:
		fallthrough
	case 1:
		newKey, _ := utils.StringToByteArray(key)

		t := trivium.New(newKey)
		success = utils.Verify(t.VirtualDecrypt(utils.VERIFIER))
		if success {
			t2 := trivium.New(newKey)
			t2.Decrypt(binManager.BIN_FILE, binManager.BIN_FILE)
			utils.Create_verifier()
		}
	case 2:
		fallthrough
	case 3:
		fallthrough
	case 4:
		newKey, _ := utils.StringToSlice(key)
		k, err := aescbc.NewKeyFrom(newKey)
		if err != nil {
			return false
		}
		verifier, _ := os.ReadFile(utils.VERIFIER)
		raw := utils.Verify(string(verifier))
		success = !raw && utils.Verify(string(aescbc.Decrypt(k, verifier)))

		if success {
			utils.Create_verifier()
			data, _ := os.ReadFile(binManager.BIN_FILE)
			data = aescbc.Decrypt(k, data)
			os.WriteFile(binManager.BIN_FILE, data, 0644)
		}
	}

	return
}

func Zip(method int) {
	switch method {
	case 1:
		huffman.Zip(binManager.BIN_FILE)
	case 2:
		lzw.Zip(binManager.CSV_PATH)
	default:
		lzw.Zip(binManager.BIN_FILE)
	}
}

func Unzip(method int) {
	switch method {
	case 1:
		huffman.Unzip(binManager.BIN_FILE, "bin")
	case 2:
		lzw.Unzip(binManager.CSV_PATH)
	default:
		lzw.Unzip(binManager.BIN_FILE)
	}
}
