package huffman

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"

	"github.com/Bernardo46-2/AEDS-III/utils"
)

// TreeNode se refere a um no de arvore de huffman
//
// O no é usado para gerar a codificação atraves da serialização
// do caminhar da arvore
//
// see https://en.wikipedia.org/wiki/Huffman_coding
type TreeNode struct {
	Char  byte
	Count int
	Path  uint
	PSize int
	Leaf  bool
	Left  *TreeNode
	Right *TreeNode
}

// ByteMap guarda a codificação de Huffman e a quantidade
// de bits que o codigo ocupa
type ByteMap struct {
	Path uint
	Size int
}

// Data é um encapsulamento para escrita e leitura do
// conteudo zipado da arvore mais o mapa de caracteres de
// huffman
type Data struct {
	Map  map[byte]ByteMap
	Zip  []byte
	Size uint
}

// getCharMap separa uma Mapa com todos os caracteres existentes
// preparando para a criacao do dicionario
func getCharMap(arr []byte) map[byte]int {
	charMap := make(map[byte]int)

	for _, b := range arr {
		charMap[b]++
	}

	return charMap
}

// getNodeList ira transformar um mapa de caracteres em um Heap
// para a construção da arvore
func getNodeHeap(charMap map[byte]int) *Heap {
	h := NewHeap()

	for char, count := range charMap {
		h.insert(&TreeNode{char, count, 0, 0, true, nil, nil})
	}

	return h
}

// getHuffmanTree desmonta o heap organizando a arvore
func getHuffmanTree(h *Heap) *TreeNode {
	for len(h.Nodes) >= 2 {
		a := h.remove()
		b := h.remove()
		h.insert(&TreeNode{0, a.Count + b.Count, 0, 0, false, a, b})
	}

	return h.Nodes[0]
}

// getCodeMap recebe uma arvore de huffman, gera a serializacao
// de caminhamento e a transfere para um mapa
func getCodeMap(node *TreeNode) map[byte]ByteMap {
	codeMap := make(map[byte]ByteMap)

	var encode func(node *TreeNode, path uint, size int, leafMap map[byte]ByteMap)
	encode = func(node *TreeNode, path uint, size int, leafMap map[byte]ByteMap) {
		if node != nil {
			node.Path = path << ((utils.ByteSize(size) * 8) - size)
			node.PSize = size
			if node.Leaf {
				leafMap[node.Char] = ByteMap{node.Path, node.PSize}
			} else {
				encode(node.Left, path<<1, size+1, leafMap)
				encode(node.Right, (path<<1)|1, size+1, leafMap)
			}
		}
	}

	encode(node, 0, 0, codeMap)

	return codeMap
}

// getSize calcula o tamanho final do arquivo baseado na codificacao
// gerada
func getSize(encoding map[byte]ByteMap, frequencies map[byte]int) int {
	totalBits := 0
	for b, freq := range frequencies {
		encodingByteMap, exists := encoding[b]
		if !exists {
			fmt.Printf("No Huffman code found for byte %v\n", b)
			continue
		}
		totalBits += freq * encodingByteMap.Size
	}
	return totalBits
}

// getZip transforma uma entrada em um arquivo comprimido
// baseado no mapa de codificacao fornecido
func getZip(input []byte, codeMap map[byte]ByteMap, size int) []byte {
	bitSize := utils.ByteSize(codeMap[input[0]].Path) * 8

	data := make([]byte, (size/8)+1)
	tmp := 0
	needle := 0
	i := 0

	for _, b := range input {
		code := codeMap[b]
		tmp += code.Size
		for code.Size > 0 {
			data[i] |= byte(code.Path >> (needle + (bitSize - 8)))

			if code.Size >= 8-needle {
				code.Path <<= 8 - needle
				code.Size -= 8 - needle
				i++
				needle = 0
			} else {
				needle += code.Size
				needle %= 8
				code.Size -= 8
			}
		}
	}

	return data
}

// save escreve os dados zipados e o mapa de huffman em arquivo de maneira
// serializada com o padrao Go
func save(path string, zip []byte, codeMap map[byte]ByteMap, old []byte) error {

	data := Data{codeMap, zip, uint(len(old))} // Dados a serem codificados
	buf := new(bytes.Buffer)                   // Cria um buffer para armazenar a codificação
	enc := gob.NewEncoder(buf)                 // Cria um novo codificador que irá escrever para o buffer
	err := enc.Encode(data)                    // Codifica os dados
	if err != nil {
		return fmt.Errorf("erro do tipo: %+v", err)
	}

	err = os.WriteFile(utils.ChangeExtension(path, ".huffman"), buf.Bytes(), 0644) // Escreve os dados codificados em um arquivo
	if err != nil {
		return fmt.Errorf("erro do tipo: %+v", err)
	}

	return nil
}

// Zip recebe um arquivo e o comprime de acordo com o padrao de
// compressao de Huffman
//
// see https://en.wikipedia.org/wiki/Huffman_coding
func Zip(path string) error {
	data, err := os.ReadFile(path) // abre o arquivo a ser zipado
	if err != nil {
		return fmt.Errorf("erro do tipo: %+v", err)
	} else if len(data) == 0 {
		return fmt.Errorf("arquivo vazio")
	}

	charMap := getCharMap(data)              // cria a lista de caracteres e ocorrencia
	nodeHeap := getNodeHeap(charMap)         // insere a lista em um Heap
	tree := getHuffmanTree(nodeHeap)         // constroi a arvore de huffman a partir do heap
	codeMap := getCodeMap(tree)              // serializa a arvore e gera o mapa de compressao
	size := getSize(codeMap, charMap)        // calcula o tamanho do arquivo comprimido
	zip := getZip(data, codeMap, size)       // transforma o texto em codigo a partir do mapa
	status := save(path, zip, codeMap, data) // sobreescreve com o arquivo zipado

	return status
}

func read(path string) ([]byte, map[byte]ByteMap, uint, error) {
	content, err := os.ReadFile(path) // abre o arquivo a ser deszipado
	if err != nil {
		return nil, nil, 0, fmt.Errorf("erro do tipo: %+v", err)
	} else if len(content) == 0 {
		return nil, nil, 0, fmt.Errorf("arquivo vazio")
	}

	// Cria um novo decodificador
	dec := gob.NewDecoder(bytes.NewBuffer(content))

	// Decodifica os dados
	var data Data
	err = dec.Decode(&data)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("erro na decodificação: %v", err)
	}

	return data.Zip, data.Map, data.Size, nil
}

func invertMap(codeMap map[byte]ByteMap) map[ByteMap]byte {
	charMap := make(map[ByteMap]byte, 0)
	for k, v := range codeMap {
		charMap[v] = k
	}
	return charMap
}

func getUnzip(data []byte, charMap map[ByteMap]byte, limit uint) []byte {
	bitSize := uint(utils.ByteSize(ByteMap{0, 0}.Path) * 8)
	original := make([]byte, 0)
	load := uint(0)
	loadPos := uint(1)
	needle := uint(0)
	size := 0
	shift := uint(0)
	passage := false
	i := uint(0)

	for i < limit {
		if !passage {
			load |= uint(utils.RotateLeft((data[needle/8]&(1<<(7-(needle%8)))), shift)) << (bitSize - (8 * loadPos))
		} else {
			load |= uint(utils.RotateRight((data[needle/8]&(1<<(7-(needle%8)))), shift)) << (bitSize - (8 * loadPos))
		}
		size++
		needle++
		if val, ok := charMap[ByteMap{load, size}]; ok {
			passage = false
			original = append(original, val)
			size = 0
			shift = (needle) % 8
			load = 0
			loadPos = 1
			i++
		} else {
			if size%8 == 0 {
				loadPos++
			}
		}
		if ((needle % 8) == 0) && !passage {
			passage = true
			shift = uint(size)
		}
	}

	return original
}

func Unzip(path string, extension string) error {
	data, codeMap, size, err := read(utils.ChangeExtension(path, ".huffman"))
	if err != nil {
		return err
	}

	charMap := invertMap(codeMap)
	original := getUnzip(data, charMap, size)

	os.WriteFile(utils.ChangeExtension(path, "2"+extension), original, 0644) // Escreve os dados codificados em um arquivo
	return nil
}
