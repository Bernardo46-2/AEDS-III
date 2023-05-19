package huffman

type Heap struct {
	Nodes []*TreeNode
}

func NewHeap() *Heap {
	return &Heap{
		Nodes: []*TreeNode{},
	}
}

func (h *Heap) insert(node *TreeNode) {
	h.Nodes = append(h.Nodes, node)
	h.heapifyUp(len(h.Nodes) - 1)
}

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

func (h *Heap) swap(i, j int) {
	h.Nodes[i], h.Nodes[j] = h.Nodes[j], h.Nodes[i]
}
