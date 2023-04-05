package btree

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Bernardo46-2/AEDS-III/data/binManager"
	"github.com/Bernardo46-2/AEDS-III/utils"
)

const BTREE_FILE string = "/btree/BTree.dat"
const BTREE_NODES_FILE string = "/btree/BTreeNodes.dat"

// TODO:
// - remove *BTree argument from functions
// - remove extra indentation on printFile()
// - create remove function
// - create find function
// - change dir concatenation to use os.join
// - store empty nodes in b.emptyNodes
// - documentation

// ====================================== Bit-Flags ====================================== //

// Bit-Flags used for removing an element from the B Tree
const (
	// Value removed without any complications
	FLAG1 = 1 << 0

	// Node has less than 50% ocupation
	FLAG2 = 1 << 1

	// Sibling node can't lend key
	FLAG3 = 1 << 2

	// Empty node
	FLAG4 = 1 << 3
)

// ====================================== Structs ====================================== //

type Key struct {
	id  int64
	ptr int64
}

type BTreeNode struct {
	address      int64
	numberOfKeys int64
	child        []int64
	keys         []Key
	leaf         int64
}

type BTree struct {
	file       string
	nodesFile  *os.File
	root       int64
	order      int
	emptyNodes []int64
}

// ====================================== Key ====================================== //

func newKey(register *binManager.Registro) Key {
	return Key{
		id:  int64(register.Pokemon.Numero),
		ptr: register.Endereco,
	}
}

func newEmptyKey() Key {
	return Key{-1, -1}
}

// ====================================== Node ====================================== //

func newNode(order int, leaf int64) *BTreeNode {
	node := BTreeNode{
		child:        make([]int64, order+1),
		keys:         make([]Key, order),
		numberOfKeys: 0,
		leaf:         leaf,
		address:      -1,
	}

	for i := 0; i < order; i++ {
		node.child[i] = -1
		node.keys[i] = newEmptyKey()
	}
	node.child[len(node.child)-1] = -1

	return &node
}

func (n *BTreeNode) write(file *os.File) {
	if n.address == -1 {
		n.address, _ = file.Seek(0, io.SeekEnd)
	} else {
		file.Seek(n.address, io.SeekStart)
	}

	binary.Write(file, binary.LittleEndian, n.numberOfKeys)
	binary.Write(file, binary.LittleEndian, n.leaf)

	for i := 0; i < len(n.keys)-1; i++ {
		binary.Write(file, binary.LittleEndian, n.child[i])
		binary.Write(file, binary.LittleEndian, n.keys[i].id)
		binary.Write(file, binary.LittleEndian, n.keys[i].ptr)
	}
	binary.Write(file, binary.LittleEndian, n.child[len(n.child)-2])
}

// self: * 2 l 3 r 5 * 9 *
//
//	|
//	v
//
// self: * 2 l _ * _ * _ *
// new:  r 5 * 9 * _ * _ *
//
// return (self, 3, new)
func (n *BTreeNode) split(tree *BTree) (int64, *Key, int64) {
	new := newNode(len(n.keys), n.leaf)
	order := len(n.keys)
	middle := order / 2

	for i := middle; i < order; i++ {
		new.keys[i-middle] = n.keys[i]
		new.child[i-middle] = n.child[i]
		n.keys[i] = newEmptyKey()
		n.child[i] = -1
		new.numberOfKeys++
	}
	new.child[order-middle] = n.child[len(n.child)-1]
	n.child[len(n.child)-1] = -1

	n.numberOfKeys = int64(middle - 1)
	carryUp := n.keys[n.numberOfKeys]
	n.keys[n.numberOfKeys] = newEmptyKey()

	n.write(tree.nodesFile)
	new.write(tree.nodesFile)

	return n.address, &carryUp, new.address
}

func (n *BTreeNode) insert(index int64, left int64, data *Key, right int64, tree *BTree) (int64, *Key, int64) {
	if index < 0 || index >= int64(len(n.keys)) {
		panic("B Tree insert error: Invalid index")
	}

	for i := int64(len(n.keys)) - 1; i > index; i-- {
		n.keys[i] = n.keys[i-1]
		n.child[i+1] = n.child[i]
	}

	n.keys[index] = *data
	n.child[index] = left
	n.child[index+1] = right
	n.numberOfKeys++

	if n.numberOfKeys == int64(len(n.keys)) {
		return n.split(tree)
	}

	n.write(tree.nodesFile)

	return -1, nil, -1
}

// ====================================== B Tree ====================================== //

func NewBTree(order int, dir string) (*BTree, error) {
	if order < 3 {
		return nil, errors.New("invalid order")
	}

	nodesFile, _ := os.Create(dir + BTREE_NODES_FILE)
	root := newNode(order, 1)
	tree := &BTree{
		root:      0,
		order:     order,
		file:      dir + BTREE_FILE,
		nodesFile: nodesFile,
	}

	root.write(nodesFile)

	return tree, nil
}

func ReadBTree(dir string) *BTree {
	file, _ := os.ReadFile(dir + BTREE_FILE)
	nodesFile, _ := os.OpenFile(dir+BTREE_NODES_FILE, os.O_RDWR|os.O_CREATE, 0644)
	root, ptr := utils.BytesToInt64(file, 0)
	order, ptr := utils.BytesToInt64(file, ptr)
	len, _ := utils.BytesToInt64(file, ptr)

	for i := int64(0); i < len; i++ {
		// do stuff
	}

	return &BTree{
		root:      root,
		order:     int(order),
		file:      dir + BTREE_FILE,
		nodesFile: nodesFile,
	}
}

func (b *BTree) Close() {
	file, _ := os.Create(b.file)

	binary.Write(file, binary.LittleEndian, b.root)
	binary.Write(file, binary.LittleEndian, int64(b.order))
	binary.Write(file, binary.LittleEndian, int64(len(b.emptyNodes)))

	for i := 0; i < len(b.emptyNodes); i++ {
		binary.Write(file, binary.LittleEndian, b.emptyNodes[i])
	}

	file.Close()
	b.nodesFile.Close()
}

func (b *BTree) nodeSize() int64 {
	node := BTreeNode{}
	s := int64(0)
	s += int64(binary.Size(node.numberOfKeys))
	s += int64(binary.Size(node.leaf))
	s += int64(binary.Size(int64(0)) * b.order)
	s += int64(binary.Size(Key{}) * (b.order - 1))
	return s
}

func (b *BTree) readNode(address int64) *BTreeNode {
	if address == -1 {
		return nil
	}

	b.nodesFile.Seek(address, io.SeekStart)
	buf := make([]byte, b.nodeSize())
	b.nodesFile.Read(buf)

	child := make([]int64, b.order+1)
	keys := make([]Key, b.order)

	numberOfKeys, ptr := utils.BytesToInt64(buf, 0)
	leaf, ptr := utils.BytesToInt64(buf, ptr)

	for i := 0; i < b.order-1; i++ {
		child[i], ptr = utils.BytesToInt64(buf, ptr)
		keys[i].id, ptr = utils.BytesToInt64(buf, ptr)
		keys[i].ptr, ptr = utils.BytesToInt64(buf, ptr)
	}
	child[len(child)-2], _ = utils.BytesToInt64(buf, ptr)
	child[len(child)-1] = -1
	keys[len(keys)-1] = newEmptyKey()

	return &BTreeNode{
		numberOfKeys: numberOfKeys,
		leaf:         leaf,
		child:        child,
		keys:         keys,
		address:      address,
	}
}

func (b *BTree) insert(node *BTreeNode, data *Key) (int64, *Key, int64) {
	l, r := int64(-1), int64(-1)
	i := int64(0)
	for i < node.numberOfKeys && data.id > node.keys[i].id {
		i++
	}

	if node.leaf == 0 {
		child := b.readNode(node.child[i])
		l, data, r = b.insert(child, data)
	}

	if data != nil {
		l, data, r = node.insert(i, l, data, r, b)
	}

	return l, data, r
}

func (b *BTree) Insert(data *Key) {
	l, m, r := b.insert(b.readNode(b.root), data)

	if m != nil {
		root := newNode(b.order, 0)
		root.child[0], root.child[1] = l, r
		root.keys[0] = *m
		root.numberOfKeys++
		root.write(b.nodesFile)
		b.root = root.address
	}
}

func (b *BTree) printFile() {
	fileEnd, _ := b.nodesFile.Seek(0, io.SeekEnd)
	b.nodesFile.Seek(0, io.SeekStart)
	reader64 := make([]byte, binary.Size(int64(0)))

	fileStart, _ := b.nodesFile.Seek(0, io.SeekStart)
	numberOfNodes := (fileEnd - fileStart) / b.nodeSize()
	tmp := int64(0)

	fmt.Printf("Root Address: %5x\n", b.root)

	for i := int64(0); i < numberOfNodes; i++ {
		currentAddress, _ := b.nodesFile.Seek(0, io.SeekCurrent)
		fmt.Printf("[%5x] || ", currentAddress)
		b.nodesFile.Read(reader64)
		tmp, _ = utils.BytesToInt64(reader64, 0)
		fmt.Printf("Size: %1d | ", tmp)
		b.nodesFile.Read(reader64)
		tmp, _ = utils.BytesToInt64(reader64, 0)
		fmt.Printf("Leaf: %1d || {", tmp)

		for i := 0; i < b.order-1; i++ {
			b.nodesFile.Read(reader64)
			tmp, _ = utils.BytesToInt64(reader64, 0)

			if tmp != -1 {
				fmt.Printf("[%5x] ", tmp)
			} else {
				fmt.Printf("[     ] ")
			}

			b.nodesFile.Read(reader64)
			tmp, _ = utils.BytesToInt64(reader64, 0)

			if tmp != -1 {
				fmt.Printf("%3d ", tmp)
			} else {
				fmt.Printf("    ")
			}

			b.nodesFile.Read(reader64)
		}

		b.nodesFile.Read(reader64)
		tmp, _ = utils.BytesToInt64(reader64, 0)

		if tmp != -1 {
			fmt.Printf("[%4x] }\n", tmp)
		} else {
			fmt.Printf("[    ] }\n")
		}
	}
	fmt.Printf("\n\n")
}

// ====================================== Tests ====================================== //

func StartBTreeFile() {
	dir := "./data/files/"
	order := 8
	tree, _ := NewBTree(order, dir)
	reader, err := binManager.InicializarControleLeitura(binManager.BIN_FILE)

	// n := 10
	n := 32
	// n := int(reader.TotalRegistros)

	for i := 0; i < n && err == nil; i++ {
		err = reader.ReadNext()
		if reader.RegistroAtual.Lapide != 1 {
			r := newKey(reader.RegistroAtual)
			tree.Insert(&r)
		}
	}

	tree.printFile()
	tree.Close()

	tree = ReadBTree(dir)

	fmt.Println("removed:", tree.Remove(8))
	fmt.Println()
	tree.printFile()

	fmt.Println("removed:", tree.Remove(7))
	tree.printFile()
	fmt.Println()

	fmt.Println("removed:", tree.Remove(6))
	tree.printFile()
	fmt.Println()

	tree.Close()
}

// ====================================== Remove ====================================== //

/*

1 case: Leaf
    remove and return flags
2 case: Not leaf
    clone max left 'A'
    clone key to be removed 'B'
    call remove function for 'A'
    call replace function, swap 'B' for 'A'

Flags:
    1 - ok
    2 - node len < 50%
    3 - empty node
    4 - key not found

Flag actions:
    1 - just return flag
    2 - check if sibling can lend key.
        if so   - rotate elements
        if no   - concat 2 siblings and get key from parent
        finally - check for errors and return flag
    3 - should happen only on the root. Updates root address
    4 - cascade back to root

*/

func (n *BTreeNode) canLendKey() bool {
	return n.numberOfKeys-1 >= int64(len(n.keys)/2)
}

func (n *BTreeNode) removeKeyLeaf(index int64) *Key {
	if index == -1 {
		index = n.numberOfKeys - 1
	}

	k := n.keys[index]

	for i := index; i < n.numberOfKeys; i++ {
		n.child[i] = n.child[i+1]
		n.keys[i] = n.keys[i+1]
	}
	n.keys[n.numberOfKeys-1] = newEmptyKey()
	n.child[n.numberOfKeys] = n.child[n.numberOfKeys+1]

	n.numberOfKeys--

	return &k
}

func (n *BTreeNode) removeKeyNonLeaf(index int64) *Key {
	if index == -1 {
		index = n.numberOfKeys - 1
	}

	return &n.keys[index]
}

func (n *BTreeNode) findMax() Key {
	return n.keys[n.numberOfKeys-1]
}

// receive node and key index
// go to key left and start iterating
func (b *BTree) maxLeft(parent *BTreeNode, childIndex int64, k *Key) Key {
	node := b.readNode(parent.child[childIndex])
	i := node.numberOfKeys - 1

	for i > 0 && k.id != node.keys[i].id {
		i--
	}
	node = b.readNode(node.child[i])

	for node.leaf == 0 {
		node = b.readNode(node.child[node.numberOfKeys])
	}

	return node.findMax()
}

func (b *BTree) replace(new *Key, old *Key, node *BTreeNode) {
	i := node.numberOfKeys - 1
	for i > 0 && node.keys[i].id > old.id {
		i--
	}

	if node.keys[i].id == old.id {
		node.keys[i] = *new
		node.write(b.nodesFile)
	} else if old.id < node.keys[i].id {
		b.replace(new, old, b.readNode(node.child[i]))
	} else {
		b.replace(new, old, b.readNode(node.child[i]))
	}
}

func (b *BTree) concatNodes(left *BTreeNode, right *BTreeNode, key *Key) *BTreeNode {
	left.keys[left.numberOfKeys] = *key
	left.numberOfKeys++

	for i := int64(0); i < right.numberOfKeys; i++ {
		left.keys[left.numberOfKeys] = right.keys[i]
		left.child[left.numberOfKeys] = right.child[i]
		left.numberOfKeys++
		right.keys[i] = newEmptyKey()
	}
	left.child[left.numberOfKeys] = right.child[right.numberOfKeys]

	left.write(b.nodesFile)

	return left
}

// remember to set the node as empty and place it in the b.emptyNodes list
func (b *BTree) borrowFromParent(parent *BTreeNode, sibling *BTreeNode, childIndex int64, left bool) int {
	flag := 0

	if parent.numberOfKeys == 1 {
		flag = FLAG4
	} else if parent.canLendKey() {
		flag = FLAG1
	} else {
		flag = FLAG2
	}

	child := b.readNode(parent.child[childIndex])

	if left {
		b.concatNodes(sibling, child, &parent.keys[childIndex])
	} else {
		b.concatNodes(child, sibling, &parent.keys[childIndex])
	}

	for i := childIndex; i < parent.numberOfKeys; {
		parent.keys[i] = parent.keys[i+1]
		i++
		parent.child[i] = parent.child[i+1]
	}

	parent.numberOfKeys--

	parent.write(b.nodesFile)

	return flag
}

func (b *BTree) borrowFromSibling(parent *BTreeNode, sibling *BTreeNode, childIndex int64, left bool) {
	var k *Key
	child := b.readNode(childIndex)

	if left {
		k = sibling.removeKeyLeaf(-1)
	} else {
		k = sibling.removeKeyLeaf(0)
	}

	child.keys[child.numberOfKeys] = parent.keys[childIndex]
	parent.keys[childIndex] = *k

	parent.write(b.nodesFile)
	sibling.write(b.nodesFile)
	child.write(b.nodesFile)
}

func (b *BTree) tryBorrowKey(parent *BTreeNode, index int64) int {
	var l, r *BTreeNode

	if index != 0 {
		l = b.readNode(parent.child[index-1])
	}
	r = b.readNode(parent.child[index+1])

	if l != nil && l.canLendKey() {
		b.borrowFromSibling(parent, l, index, true)
	} else if r != nil && r.canLendKey() {
		b.borrowFromSibling(parent, r, index, false)
	} else if l != nil {
		return b.borrowFromParent(parent, l, index, true)
	} else if r != nil {
		return b.borrowFromParent(parent, r, index, false)
	}

	return FLAG1
}

func (b *BTree) parseFlag(flag int, parent *BTreeNode, childIndex int64, k *Key) int {
	if FLAG1&flag != 0 {
		// nothing to see here
	} else if FLAG2&flag != 0 {
		flag = b.tryBorrowKey(parent, childIndex)
	} else if FLAG3&flag != 0 {
		maxLeft := b.maxLeft(parent, childIndex, k)
		_, flag, _ = b.remove(parent.address, maxLeft.id)
		b.replace(&maxLeft, k, parent)
	} else if FLAG4&flag != 0 {
		root := b.readNode(b.root)
		b.root = root.child[0]
		root.write(b.nodesFile)
	}

	return flag
}

func (b *BTree) removeFromNode(index int64, node *BTreeNode) (*Key, int) {
	var k *Key
	flag := 0
	easyRemove := node.leaf == 1 && node.canLendKey()
	midRemove := !easyRemove && node.leaf == 1
	hardRemove := !(midRemove || easyRemove)

	if easyRemove {
		flag |= FLAG1
	} else if midRemove {
		flag |= FLAG2
	} else if hardRemove {
		flag |= FLAG3
	}

	if node.leaf == 1 {
		k = node.removeKeyLeaf(index)
		node.write(b.nodesFile)
	} else {
		k = node.removeKeyNonLeaf(index)
	}

	return k, flag
}

func (b *BTree) remove(nodeAddress int64, id int64) (*Key, int, int64) {
	if nodeAddress == -1 {
		return nil, FLAG1, -1
	}

	var k *Key
	node := b.readNode(nodeAddress)
	flag := 0

	i := node.numberOfKeys - 1
	for i > 0 && node.keys[i].id > id {
		i--
	}

	if node.keys[i].id == id {
		k, flag = b.removeFromNode(i, node)
		return k, flag, i
	} else if node.keys[i].id < id {
		k, flag, _ = b.remove(node.child[i+1], id)
	} else {
		k, flag, _ = b.remove(node.child[i], id)
	}

	flag = b.parseFlag(flag, node, i, k)

	return k, flag, i
}

func (b *BTree) Remove(id int64) *Key {
	k, flag, childIndex := b.remove(b.root, id)

	if FLAG2&flag != 0 {
		root := b.readNode(b.root)
		b.tryBorrowKey(root, childIndex)
		root.write(b.nodesFile)
	} else if FLAG3&flag != 0 {
		root := b.readNode(b.root)
		node := b.readNode(root.child[childIndex])
		maxLeft := node.findMax()
		b.remove(root.address, maxLeft.id)
		b.replace(&maxLeft, k, root)
		root.write(b.nodesFile)
	} else if FLAG4&flag != 0 {
		root := b.readNode(b.root)
		b.root = root.child[0]
		root.write(b.nodesFile)
	}

	return k
}
