package huffman

import (
	"bytes"
	"fmt"
	"os"
)

type TreeNode struct {
	Char  byte
	Count int
	Path  byte
	PSize int
	Leaf  bool
	Left  *TreeNode
	Right *TreeNode
}

type ByteMap struct {
	Path byte
	Size int
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

func createCode(node *TreeNode, path byte, pos int) {
	if node != nil {
		node.Path = path
		node.PSize = pos
		createCode(node.Left, path<<1, pos+1)
		createCode(node.Right, (path<<1)|1, pos+1)
	}
}

func getCodeMap(node *TreeNode, leafMap map[byte]ByteMap) {
	if node != nil {
		if node.Leaf {
			leafMap[node.Char] = ByteMap{node.Path, node.PSize}
		}
		getCodeMap(node.Left, leafMap)
		getCodeMap(node.Right, leafMap)
	}
}

func Zip(path string) error {
	// Abertura do arquivo a zipar
	content, err := os.ReadFile(path)
	if err != nil || len(content) == 0 {
		return fmt.Errorf("erro do tipo: %s", err.Error())
	}

	charMap := getCharMap(content)    // cria a lista de caracteres e ocorrencia
	nodeHeap := getNodeHeap(charMap)  // insere a lista em um Heap
	tree := getHuffmanTree(nodeHeap)  // constroi a arvore de huffman a partir do heap
	createCode(tree, 0, 0)            // serializa o caminhar da arvore nos nos
	codeMap := make(map[byte]ByteMap) // cria um mapa para receber a codificacao
	getCodeMap(tree, codeMap)         // cria o mapa de codificacao a partir da arvore
	// byteCode := compress(content, codeMap)

	fmt.Printf("%+v", byteCode)
	return nil
}
