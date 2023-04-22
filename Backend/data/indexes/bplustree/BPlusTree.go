package bplustree

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Bernardo46-2/AEDS-III/models"
	"github.com/Bernardo46-2/AEDS-III/utils"
)

const PATH string = "bplustree/"
const NODES string = "BPlusTreeNodes.bin"
const HEADER string = "BPlusTree.bin"
const NULL int64 = -1

// Bit-Flags used for removing an element from the B+ Tree
const (
	// Value removed without any complications
	OK = 1 << 0

	// Update parent
	REPLACE = 1 << 1

	// Size problem
	HELP = 1 << 2

	// Is leaf
	LEAF = 1 << 3

	// ... Can lend..
	CAN_LEND = 1 << 4

	// Empty node
	EMPTY = 1 << 5
)

// ====================================== Structs ====================================== //

type Key struct {
	Id  float64
	Ptr int64
}

type BPlusTreeNode struct {
	address      int64
	numberOfKeys int64
	child        []int64
	keys         []Key
	leaf         int64
	next         int64
}

type BPlusTree struct {
	file       string
	nodesFile  *os.File
	root       int64
	order      int
	emptyNodes []int64
}

type Reader interface {
	ReadNextGeneric() (any, bool, int64, error)
}

type IndexableObject interface {
	GetFieldF64(fieldName string) (float64, int64)
}

// ====================================== Key ====================================== //

func newEmptyKey() Key {
	return Key{float64(NULL), NULL}
}

func (k *Key) compareTo(other *Key) float64 {
	diff := k.Id - other.Id
	if diff == 0 {
		diff = float64(k.Ptr - other.Ptr)
	}
	return diff
}

// ====================================== Node ====================================== //

func newNode(order int, leaf int64, address int64) *BPlusTreeNode {
	node := BPlusTreeNode{
		child:        make([]int64, order+1),
		keys:         make([]Key, order),
		numberOfKeys: 0,
		leaf:         leaf,
		address:      address,
		next:         NULL,
	}

	for i := 0; i < order; i++ {
		node.child[i] = NULL
		node.keys[i] = newEmptyKey()
	}
	node.child[len(node.child)-1] = NULL

	return &node
}

func (n *BPlusTreeNode) write(file *os.File) {
	if n.address == NULL {
		n.address, _ = file.Seek(0, io.SeekEnd)
	} else {
		file.Seek(n.address, io.SeekStart)
	}

	binary.Write(file, binary.LittleEndian, n.numberOfKeys)
	binary.Write(file, binary.LittleEndian, n.leaf)

	for i := 0; i < len(n.keys)-1; i++ {
		binary.Write(file, binary.LittleEndian, n.child[i])
		binary.Write(file, binary.LittleEndian, n.keys[i].Id)
		binary.Write(file, binary.LittleEndian, n.keys[i].Ptr)
	}
	binary.Write(file, binary.LittleEndian, n.child[len(n.child)-2])

	binary.Write(file, binary.LittleEndian, n.next)
}

// self: * 2 l 3 r 5 * 9 *
//
//	|
//	v
//
// self: * 2 l 3 * _ * _ *
// new:  r 5 * 9 * _ * _ *
//
// return (self, 3, new)
func (n *BPlusTreeNode) split(tree *BPlusTree) (int64, *Key, int64) {
	new := newNode(len(n.keys), n.leaf, tree.popEmptyNode())

	order := len(n.keys)
	middle := order / 2

	for i := middle; i < order; i++ {
		new.keys[i-middle] = n.keys[i]
		new.child[i-middle] = n.child[i]
		n.keys[i] = newEmptyKey()
		n.child[i] = NULL
		new.numberOfKeys++
	}
	new.child[order-middle] = n.child[len(n.child)-1]
	n.child[len(n.child)-1] = NULL

	var carryUp Key
	if n.leaf == 1 {
		n.numberOfKeys = int64(middle)
		carryUp = n.keys[n.numberOfKeys-1]
	} else {
		n.numberOfKeys = int64(middle - 1)
		carryUp = n.keys[n.numberOfKeys]
		n.keys[n.numberOfKeys] = newEmptyKey()
	}

	new.next = n.next
	new.write(tree.nodesFile)
	n.next = new.address
	n.write(tree.nodesFile)

	return n.address, &carryUp, new.address
}

func (n *BPlusTreeNode) insert(index int64, left int64, data *Key, right int64, tree *BPlusTree) (int64, *Key, int64) {
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

	return NULL, nil, NULL
}

func (n *BPlusTreeNode) find(k *Key) int64 {
	i := n.numberOfKeys - 1

	for i >= 0 && n.keys[i].compareTo(k) > 0 {
		i--
	}

	if n.keys[i].compareTo(k) != 0 {
		i = -1
	}

	return i
}

func (n *BPlusTreeNode) canLendKey() bool {
	return n.numberOfKeys >= int64(len(n.keys)/2)
}

func (n *BPlusTreeNode) removeKeyLeaf(index int64) (*Key, int64) {
	if index == NULL {
		index = n.numberOfKeys - 1
	}

	k := n.keys[index]
	leftChild := n.child[index]

	for i := index; i < n.numberOfKeys; i++ {
		n.child[i] = n.child[i+1]
		n.keys[i] = n.keys[i+1]
	}
	n.keys[n.numberOfKeys-1] = newEmptyKey()
	n.child[n.numberOfKeys] = n.child[n.numberOfKeys+1]

	n.numberOfKeys--

	return &k, leftChild
}

func (n *BPlusTreeNode) max() *Key {
	if n.numberOfKeys == 0 {
		return nil
	}
	return &n.keys[n.numberOfKeys-1]
}

func (n *BPlusTreeNode) getStatus() int {
	flag := 0

	if n.leaf == 1 {
		flag |= LEAF
	}

	if n.canLendKey() {
		flag |= CAN_LEND
	}

	return flag
}

// ====================================== B+ Tree ====================================== //

func NewBPlusTree(order int, path string, field string) (*BPlusTree, error) {
	if order < 3 {
		return nil, errors.New("invalid order")
	}

	tree_path := filepath.Join(path, PATH)
	tree_nodes := filepath.Join(tree_path, field+"_"+NODES)
	tree_header := filepath.Join(tree_path, field+"_"+HEADER)
	os.MkdirAll(tree_path, 0755)
	nodesFile, _ := os.Create(tree_nodes)
	root := newNode(order, 1, NULL)
	tree := &BPlusTree{
		root:       0,
		order:      order,
		file:       tree_header,
		nodesFile:  nodesFile,
		emptyNodes: make([]int64, 0),
	}

	root.write(nodesFile)

	return tree, nil
}

func ReadBPlusTree(dir string, field string) (*BPlusTree, error) {
	tree_path := filepath.Join(dir, PATH)
	tree_nodes := filepath.Join(tree_path, field+"_"+NODES)
	tree_header := filepath.Join(tree_path, field+"_"+HEADER)

	file, err := os.ReadFile(tree_header)
	if err != nil {
		return nil, err
	}

	nodesFile, _ := os.OpenFile(tree_nodes, os.O_RDWR|os.O_CREATE, 0644)
	root, ptr := utils.BytesToInt64(file, 0)
	order, ptr := utils.BytesToInt64(file, ptr)
	len, ptr := utils.BytesToInt64(file, ptr)
	emptyNodes := make([]int64, len)

	for i := int64(0); i < len; i++ {
		emptyNodes[i], ptr = utils.BytesToInt64(file, ptr)
	}

	return &BPlusTree{
		root:       root,
		order:      int(order),
		file:       tree_header,
		nodesFile:  nodesFile,
		emptyNodes: emptyNodes,
	}, nil
}

func (b *BPlusTree) Close() {
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

func (b *BPlusTree) popEmptyNode() int64 {
	if len(b.emptyNodes) > 0 {
		address := b.emptyNodes[0]
		if len(b.emptyNodes) > 1 {
			b.emptyNodes = b.emptyNodes[1:]
		} else {
			b.emptyNodes = []int64{}
		}
		return address
	}

	return NULL
}

func (b *BPlusTree) nodeSize() int64 {
	node := BPlusTreeNode{}
	s := int64(0)
	s += int64(binary.Size(node.numberOfKeys))
	s += int64(binary.Size(node.leaf))
	s += int64(binary.Size(int64(0)) * b.order)
	s += int64(binary.Size(Key{}) * (b.order - 1))
	s += int64(binary.Size(node.next))
	return s
}

func (b *BPlusTree) readNode(address int64) *BPlusTreeNode {
	if address == NULL {
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
		keys[i].Id, ptr = utils.BytesToFloat64(buf, ptr)
		keys[i].Ptr, ptr = utils.BytesToInt64(buf, ptr)
	}
	child[len(child)-2], ptr = utils.BytesToInt64(buf, ptr)
	child[len(child)-1] = NULL
	keys[len(keys)-1] = newEmptyKey()

	next, _ := utils.BytesToInt64(buf, ptr)

	return &BPlusTreeNode{
		numberOfKeys: numberOfKeys,
		leaf:         leaf,
		child:        child,
		keys:         keys,
		address:      address,
		next:         next,
	}
}

func (b *BPlusTree) insert(node *BPlusTreeNode, data *Key) (int64, *Key, int64) {
	l, r := NULL, NULL
	i := int64(0)
	for i < node.numberOfKeys && data.compareTo(&node.keys[i]) > 0 {
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

func (b *BPlusTree) Insert(data *Key) {
	l, m, r := b.insert(b.readNode(b.root), data)

	if m != nil {
		root := newNode(b.order, 0, b.popEmptyNode())
		root.child[0], root.child[1] = l, r
		root.keys[0] = *m
		root.numberOfKeys++
		root.write(b.nodesFile)
		b.root = root.address
	}
}

func (b *BPlusTree) PrintFile() {
	fileEnd, _ := b.nodesFile.Seek(0, io.SeekEnd)
	b.nodesFile.Seek(0, io.SeekStart)
	reader64 := make([]byte, binary.Size(int64(0)))

	fileStart, _ := b.nodesFile.Seek(0, io.SeekStart)
	numberOfNodes := (fileEnd - fileStart) / b.nodeSize()
	i64 := int64(0)
	f64 := float64(0)

	fmt.Printf("Root Address: %5x\n", b.root)

	for i := int64(0); i < numberOfNodes; i++ {
		currentAddress, _ := b.nodesFile.Seek(0, io.SeekCurrent)
		fmt.Printf("[%5x] || ", currentAddress)
		b.nodesFile.Read(reader64)
		i64, _ = utils.BytesToInt64(reader64, 0)
		fmt.Printf("Size: %1d | ", i64)
		b.nodesFile.Read(reader64)
		i64, _ = utils.BytesToInt64(reader64, 0)
		fmt.Printf("Leaf: %1d || {", i64)

		for i := 0; i < b.order-1; i++ {
			b.nodesFile.Read(reader64)
			i64, _ = utils.BytesToInt64(reader64, 0)

			if i64 != NULL {
				fmt.Printf("[%5x] ", i64)
			} else {
				fmt.Printf("[     ] ")
			}

			b.nodesFile.Read(reader64)
			f64, _ = utils.BytesToFloat64(reader64, 0)

			if f64 != float64(NULL) {
				fmt.Printf("%6.2f ", f64)
			} else {
				fmt.Printf("       ")
			}

			b.nodesFile.Read(reader64)
			i64, _ = utils.BytesToInt64(reader64, 0)

			if i64 != NULL {
				fmt.Printf("%5x ", i64)
			} else {
				fmt.Printf("      ")
			}
		}

		b.nodesFile.Read(reader64)
		i64, _ := utils.BytesToInt64(reader64, 0)

		if i64 != NULL {
			fmt.Printf("[%4x] ", i64)
		} else {
			fmt.Printf("[    ] ")
		}

		b.nodesFile.Read(reader64)
		i64, _ = utils.BytesToInt64(reader64, 0)

		if i64 != NULL {
			fmt.Printf("|| Next: %4x }\n", i64)
		} else {
			fmt.Printf("|| Next: %4x }\n", NULL)
		}
	}
	fmt.Printf("\n\n")
}

func (b *BPlusTree) concatLeaf(left *BPlusTreeNode, right *BPlusTreeNode) *BPlusTreeNode {
	for i := int64(0); i < right.numberOfKeys; i++ {
		left.keys[left.numberOfKeys] = right.keys[i]
		left.child[left.numberOfKeys] = right.child[i]
		left.numberOfKeys++
		right.keys[i] = newEmptyKey()
		right.child[i] = NULL
	}
	left.child[left.numberOfKeys] = right.child[right.numberOfKeys]
	right.child[right.numberOfKeys] = NULL
	right.numberOfKeys = 0
	right.leaf = 0
	left.next = right.next
	right.next = NULL

	left.write(b.nodesFile)
	right.write(b.nodesFile)
	b.emptyNodes = append(b.emptyNodes, right.address)

	return left
}

func (b *BPlusTree) mergeLeaf(node *BPlusTreeNode, l *BPlusTreeNode, r *BPlusTreeNode, keyIndex int64) {
	b.concatLeaf(l, r)

	for i := keyIndex; i < node.numberOfKeys; {
		node.keys[i] = node.keys[i+1]
		i++
		node.child[i] = node.child[i+1]
	}

	node.numberOfKeys--
	node.write(b.nodesFile)
}

func (b *BPlusTree) borrowFromLeaf(node *BPlusTreeNode, left *BPlusTreeNode, right *BPlusTreeNode, index int64) {
	k, _ := right.removeKeyLeaf(0)

	left.keys[left.numberOfKeys] = *k
	left.numberOfKeys++
	node.keys[index] = *k

	node.write(b.nodesFile)
	left.write(b.nodesFile)
	right.write(b.nodesFile)
}

func (b *BPlusTree) borrowFromNonLeaf(parent *BPlusTreeNode, left *BPlusTreeNode, right *BPlusTreeNode, index int64) {
	k, child := right.removeKeyLeaf(0)

	left.keys[left.numberOfKeys] = parent.keys[index]
	left.numberOfKeys++
	left.child[left.numberOfKeys] = child

	parent.keys[index] = *k

	parent.write(b.nodesFile)
	left.write(b.nodesFile)
	right.write(b.nodesFile)
}

func (b *BPlusTree) concatNodesOG(left *BPlusTreeNode, right *BPlusTreeNode, key *Key) *BPlusTreeNode {
	left.keys[left.numberOfKeys] = *key
	left.numberOfKeys++
	for i := int64(0); i < right.numberOfKeys; i++ {
		left.keys[left.numberOfKeys] = right.keys[i]
		left.child[left.numberOfKeys] = right.child[i]
		left.numberOfKeys++
		right.keys[i] = newEmptyKey()
		right.child[i] = NULL
	}
	left.child[left.numberOfKeys] = right.child[right.numberOfKeys]
	right.child[right.numberOfKeys] = NULL
	right.numberOfKeys = 0
	left.next = right.next
	right.next = NULL
	b.emptyNodes = append(b.emptyNodes, right.address)

	left.write(b.nodesFile)
	right.write(b.nodesFile)

	return left
}

func (b *BPlusTree) borrowFromParentOG(node *BPlusTreeNode, l *BPlusTreeNode, r *BPlusTreeNode, keyIndex int64) {
	b.concatNodesOG(l, r, &node.keys[keyIndex])

	for i := keyIndex; i < node.numberOfKeys; {
		node.keys[i] = node.keys[i+1]
		i++
		node.child[i] = node.child[i+1]
	}

	node.numberOfKeys--
	node.write(b.nodesFile)
}

func (b *BPlusTree) tryBorrowKey(node *BPlusTreeNode, index int64) int {
	var l, r, lefter *BPlusTreeNode

	l = b.readNode(node.child[index])
	r = b.readNode(node.child[index+1])

	if index != 0 {
		lefter = b.readNode(node.child[index-1])
	}

	if r == nil {
		r = l
		l = lefter
	}

	flag := r.getStatus()

	if flag == LEAF|CAN_LEND { // replace
		b.borrowFromLeaf(node, l, r, index)
	} else if flag == LEAF { // merge
		b.mergeLeaf(node, l, r, index)
	} else if flag == CAN_LEND { // borrow B
		b.borrowFromNonLeaf(node, l, r, index)
	} else { // downgrade parent
		b.borrowFromParentOG(node, l, r, index)
	}

	if node.canLendKey() {
		flag = OK
	} else {
		flag = HELP
	}

	return flag
}

func (b *BPlusTree) replaceKey(n *BPlusTreeNode, k *Key, kk *Key, flag int) int {
	walter := 0
	if REPLACE&flag != 0 {
		i := n.find(k)
		fmt.Println("Found:", i)
		fmt.Println("Searching:", k)
		if i != NULL {
			n.keys[i] = *kk
			n.write(b.nodesFile)
		} else {
			walter = REPLACE
		}
	}

	return walter
}

func (b *BPlusTree) parseFlag(flag int, node *BPlusTreeNode, index int64, k *Key, kk *Key) int {
	walter := b.replaceKey(node, k, kk, flag)

	if HELP&flag != 0 {
		walter |= b.tryBorrowKey(node, index)
	}

	if node.numberOfKeys == 0 {
		walter |= EMPTY
	}

	return walter
}

func (b *BPlusTree) removeFromNode(index int64, node *BPlusTreeNode) (*Key, *Key, int) {
	var k *Key
	flag := 0

	if index == node.numberOfKeys-1 {
		flag = REPLACE
	}

	if node.canLendKey() {
		flag |= OK
	} else {
		flag |= HELP
	}

	k, _ = node.removeKeyLeaf(index)
	max := node.max()
	node.write(b.nodesFile)

	if node.numberOfKeys == 0 {
		flag |= EMPTY
	}

	return k, max, flag
}

func (b *BPlusTree) remove(address int64, old *Key) (*Key, *Key, int, int64) {
	if address == NULL {
		return nil, nil, OK, NULL
	}

	var k, kk *Key
	node := b.readNode(address)
	flag := 0

	i := node.numberOfKeys - 1
	for i > 0 && node.keys[i].compareTo(old) > 0 {
		i--
	}

	if node.leaf == 1 && node.keys[i].compareTo(old) == 0 {
		k, kk, flag = b.removeFromNode(i, node)
	} else if node.keys[i].compareTo(old) < 0 {
		k, kk, flag, _ = b.remove(node.child[i+1], old)
		flag = b.parseFlag(flag, node, i, k, kk)
	} else {
		k, kk, flag, _ = b.remove(node.child[i], old)
		flag = b.parseFlag(flag, node, i, k, kk)
	}

	return k, kk, flag, i
}

func (b *BPlusTree) Remove(old *Key) *Key {
	k, kk, flag, _ := b.remove(b.root, old)
	root := b.readNode(b.root)

	b.replaceKey(root, k, kk, flag)

	if flag&EMPTY != 0 {
		if root.child[0] != NULL {
			b.emptyNodes = append(b.emptyNodes, b.root)
			b.root = root.child[0]
			root.child[0] = NULL
			root.next = NULL
		}
		root.write(b.nodesFile)
	}

	return k
}

func (n *BPlusTreeNode) find2(id float64) (*Key, int64) {
	var k *Key
	address := NULL
	i := n.numberOfKeys - 1

	for i > 0 && n.keys[i].Id > id {
		i--
	}

	if n.keys[i].Id == id {
		k = &n.keys[i]
	} else if n.keys[i].Id < id {
		address = n.child[i+1]
	} else {
		address = n.child[i]
	}

	return k, address
}

func (b *BPlusTree) Find(id float64) *Key {
	var k *Key
	var address int64
	node := b.readNode(b.root)

	for node != nil && k == nil {
		k, address = node.find2(id)
		node = b.readNode(address)
	}

	return k
}

func (b *BPlusTree) findNode(id float64) *BPlusTreeNode {
	node := b.readNode(b.root)

	for i := int64(0); node.leaf == 0; i++ {
		if i < node.numberOfKeys-1 {
			if node.keys[i].Id >= id {
				node = b.readNode(node.child[i])
				i = -1
			}
		} else {
			if node.keys[i].Id >= id {
				node = b.readNode(node.child[i])
			} else {
				node = b.readNode(node.child[i+1])
			}
			i = -1
		}
	}

	return node
}

func (b *BPlusTree) FindRange(start float64, end float64) ([]int64, error) {
	if start > end {
		return nil, errors.New("invalid indexes")
	}

	node := b.findNode(start)
	index := int64(0)

	for node.keys[index].Id < start && index < node.numberOfKeys {
		index++
	}

	if index == node.numberOfKeys {
		return nil, errors.New("not found")
	}

	addresses := make([]int64, 0)

	fmt.Println(node)
	for start < end {
		start = node.keys[index].Id
		addresses = append(addresses, node.keys[index].Ptr)

		if index == node.numberOfKeys-1 {
			node = b.readNode(node.next)
			index = -1
			fmt.Println(node)
		}

		index++
	}

	fmt.Println(addresses)

	return addresses, nil
}

func Create(pokemon models.Pokemon, pokeAddress int64, path string, fields []string) {
	for _, field := range fields {
		tree, _ := ReadBPlusTree(path, field)
		id, address := pokemon.GetFieldF64(field)
		k := Key{id, address}
		tree.Insert(&k)
		tree.Close()
	}
}

func Update(old models.Pokemon, new models.Pokemon, pokeAddress int64, path string, fields []string) {
	for _, field := range fields {
		tree, _ := ReadBPlusTree(path, field)
		kId, kAddress := old.GetFieldF64(field)
		kkId, kkAddress := new.GetFieldF64(field)
		k := Key{kId, kAddress}
		kk := Key{kkId, kkAddress}
		tree.Remove(&k)
		tree.Insert(&kk)
		tree.Close()
	}
}

func Delete(pokemon models.Pokemon, pokeAddress int64, path string, fields []string) {
	for _, field := range fields {
		tree, _ := ReadBPlusTree(path, field)
		id, address := pokemon.GetFieldF64(field)
		removed := tree.Remove(&Key{id, address})
		if removed == nil {
			fmt.Println("not found")
		} else {
			fmt.Printf("removed: %f | %x\n", removed.Id, removed.Ptr)
		}

		fmt.Println(field)
		tree.PrintFile()
		fmt.Println()

		tree.Close()
	}
}

// ====================================== Init ====================================== //

func StartBPlusTreeFile(dir string, field string, controler Reader) error {
	order := 8
	tree, _ := NewBPlusTree(order, dir, field)

	for {
		objInterface, isDead, _, err := controler.ReadNextGeneric()
		if err != nil {
			break
		}

		obj, ok := objInterface.(IndexableObject)
		if !ok {
			return fmt.Errorf("failed to convert object to IndexableObject\n%+v", objInterface)
		}

		if !isDead {
			id, address := obj.GetFieldF64(field)
			r := Key{id, address}
			tree.Insert(&r)
		}
	}

	tree.Close()
	return nil
}
