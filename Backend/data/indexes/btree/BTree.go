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

const BTREE_FILE string = "BTree.dat"
const BTREE_NODES_FILE string = "BTreeNodes.dat"

// TODO:
// - remove *BTree argument from functions
// - remove extra indentation on printFile()
// - create remove function
// - create find function

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
	nodesFile, _ := os.OpenFile(dir + BTREE_NODES_FILE, os.O_RDWR|os.O_CREATE, 0644)
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

	for i := 0; i < 10 /* int(reader.TotalRegistros) */ && err == nil; i++ {
		err = reader.ReadNext()
		if reader.RegistroAtual.Lapide != 1 {
			r := newKey(reader.RegistroAtual)
			tree.Insert(&r)
		}
	}

	tree.printFile()
    tree.Close()

	tree = ReadBTree(dir)
    
    fmt.Println("removed:", tree.Remove(7))
    fmt.Println()
	tree.printFile()
    
	tree.Close()
}

// ====================================== Remove ====================================== //

func (n *BTreeNode) maxLeft(id int64, file *os.File) *Key {
    // node := 
    
    for n.leaf == 0 {
        
    }

    return nil
}

func (n *BTreeNode) removeKey(index int64) *Key {
    k := n.keys[index]
    
    for i := index; i < n.numberOfKeys; i++ {
        n.child[i] = n.child[i + 1]
        n.keys[i] = n.keys[i + 1]
    }
    n.keys[n.numberOfKeys] = newEmptyKey()
    n.child[n.numberOfKeys] = n.child[n.numberOfKeys + 1]
    n.child[n.numberOfKeys + 1] = -1

    n.numberOfKeys--

    return &k
}

func (n *BTreeNode) remove(index int64, tree *BTree) *Key {
    fmt.Println("n nodes:", n.numberOfKeys - 1)
    fmt.Println("len:", int64(len(n.keys)))
    fmt.Println("can remove?:", n.numberOfKeys - 1 >= int64(len(n.keys) / 2))

    var k *Key

    if n.leaf == 1 && n.numberOfKeys - 1 >= int64(len(n.keys) / 2) {
        k = n.removeKey(index)
        n.write(tree.nodesFile)
    } else {
        
    }

    return k
}

func (b *BTree) remove(nodeAddress int64 , id int64) *Key {
    if nodeAddress == -1 {
        return nil
    }
    
    var k *Key
    node := b.readNode(nodeAddress)

    i := int64(0)
	for i < node.numberOfKeys && node.keys[i].id < id {
		i++
	}

	if node.keys[i].id == id {
        k = node.remove(i, b)
    } else if node.keys[i].id < id {
        k = b.remove(node.child[i], id)
    } else {
        k = b.remove(node.child[i + 1], id)
    }
    
    return k
}

func (b *BTree) Remove(id int64) *Key {
	return b.remove(b.root, id)
}
