package huffman

import (
	"fmt"
	"os"
)

type TreeNode struct {
	Char  byte
	Count int
	Path  byte
	PSize int
	Left  *TreeNode
	Right *TreeNode
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
		h.insert(&TreeNode{char, count, 0, 0, nil, nil})
	}

	return h
}

// getHuffmanTree desmonta o heap organizando a arvore
func getHuffmanTree(h *Heap) *TreeNode {
	for len(h.Nodes) >= 2 {
		a := h.remove()
		b := h.remove()
		h.insert(&TreeNode{0, a.Count + b.Count, 0, 0, a, b})
	}

	return h.Nodes[0]
}

func encode(node *TreeNode, path byte, pos int) {
	if node != nil {
		node.Path = path
		node.PSize = pos
		encode(node.Left, path, pos+1)
		encode(node.Right, path|(1<<pos), pos+1)
	}
}

func preOrder(node *TreeNode) {
	if node != nil {
		if node.Char != 0 {
			fmt.Printf("%1c | %3d | %b \n", node.Char, node.Count, node.Path)
		} else {
			fmt.Printf("  | %3d | %b \n", node.Count, node.Path)
		}
		preOrder(node.Left)
		preOrder(node.Right)
	}
}

func Zip(path string) error {
	// Abertura do arquivo a zipar
	content, err := os.ReadFile(path)
	if err != nil || len(content) == 0 {
		return fmt.Errorf("erro do tipo: %s", err.Error())
	}

	charMap := getCharMap(content)
	nodeHeap := getNodeHeap(charMap)
	tree := getHuffmanTree(nodeHeap)
	encode(tree, 0, 0)
	preOrder(tree)

	return nil
}
