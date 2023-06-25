// heap fornece a implementacao de um heap minimo adaptado exclusivamente para a construcao
// e manipulacao da arvore de huffman
package huffman

// Heap representa uma estrutura de heap mínimo.
type Heap struct {
	Nodes []*TreeNode // Nodes é uma slice de ponteiros para TreeNode.
}

// NewHeap cria um novo heap vazio e retorna um ponteiro para ele.
func NewHeap() *Heap {
	return &Heap{
		Nodes: []*TreeNode{},
	}
}

// insert adiciona um novo nó ao heap, garantindo que a propriedade de heap mínimo seja mantida.
func (h *Heap) insert(node *TreeNode) {
	h.Nodes = append(h.Nodes, node)
	h.heapifyUp(len(h.Nodes) - 1)
}

// remove remove e retorna o nó com o valor mínimo do heap.
// Se o heap estiver vazio, retorna nil. Após a remoção, garante que a propriedade de heap
// mínimo seja mantida.
func (h *Heap) remove() *TreeNode {
	if len(h.Nodes) == 0 {
		return nil
	}
	min := h.Nodes[0]
	lastIndex := len(h.Nodes) - 1
	h.Nodes[0] = h.Nodes[lastIndex]
	h.Nodes = h.Nodes[:lastIndex]
	h.heapifyDown(0)
	return min
}

// heapifyUp move o nó no índice especificado para cima no heap (se necessário) para garantir
// que a propriedade de heap mínimo seja mantida.
func (h *Heap) heapifyUp(index int) {
	for index > 0 {
		parentIndex := (index - 1) / 2
		if h.Nodes[parentIndex].Count <= h.Nodes[index].Count {
			break
		}
		h.swap(parentIndex, index)
		index = parentIndex
	}
}

// heapifyDown move o nó no índice especificado para baixo no heap (se necessário) para
// garantir que a propriedade de heap mínimo seja mantida.
func (h *Heap) heapifyDown(index int) {
	for {
		leftChildIndex := 2*index + 1
		rightChildIndex := 2*index + 2
		smallest := index

		if leftChildIndex < len(h.Nodes) && h.Nodes[leftChildIndex].Count < h.Nodes[smallest].Count {
			smallest = leftChildIndex
		}
		if rightChildIndex < len(h.Nodes) && h.Nodes[rightChildIndex].Count < h.Nodes[smallest].Count {
			smallest = rightChildIndex
		}

		if smallest == index {
			break
			// Propriedade do heap está satisfeita
		}

		h.swap(index, smallest)
		index = smallest
	}
}

// swap troca os nós nos índices i e j.
func (h *Heap) swap(i, j int) {
	h.Nodes[i], h.Nodes[j] = h.Nodes[j], h.Nodes[i]
}
