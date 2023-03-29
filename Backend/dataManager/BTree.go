package dataManager

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
    "fmt"

	"github.com/Bernardo46-2/AEDS-III/utils"
)

const BTREE_FILE string = "data/BTree.dat"

// TODO: 
// - remove tree pointer from functions
// - remove extra indentation on printTree
// - create remove function
// - create find function


// ====================================== Structs ====================================== //

type Key struct {
    id int64
    ptr int64
}

type BTreeNode struct {
    address int64
    numberOfKeys int64
    child []int64
    keys []Key
    leaf int64
}

type BTree struct {
    file *os.File
    root int64
    order int
}

// ====================================== Key ====================================== //

func newKey(register *Registro) Key {
    return Key { 
        id: int64(register.Pokemon.Numero),
        ptr: register.Endereco,
    }
}

func newEmptyKey() Key {
    return Key { -1, -1 }
}


// ====================================== Node ====================================== //

func newNode(order int, leaf int64) *BTreeNode {
    node := BTreeNode {
        child: make([]int64, order + 1),
        keys: make([]Key, order),
        numberOfKeys: 0,
        leaf: leaf,
        address: -1,
    }

    for i := 0; i < order; i++ {
        node.child[i] = -1
        node.keys[i] = newEmptyKey()
    }
    node.child[len(node.child)-1] = -1

    return &node
}

//
// self: * 2 l 3 r 5 * 9 *
//
//            |
//            v
//  
// self: * 2 l _ * _ * _ *
// new:  r 5 * 9 * _ * _ *
//
// return (self, 3, new)
//
func (n *BTreeNode) split(tree *BTree) (int64, *Key, int64) {
    new := newNode(len(n.keys), n.leaf)
    order := len(n.keys)
    middle := order / 2

    for i := middle; i < order; i++ {
        new.keys[i - middle] = n.keys[i]
        new.child[i - middle] = n.child[i]
        n.keys[i] = newEmptyKey()
        n.child[i] = -1
        new.numberOfKeys++
    }
    new.child[order - middle] = n.child[len(n.child)-1]
    n.child[len(n.child)-1] = -1
    
    n.numberOfKeys = int64(middle - 1)
    carryUp := n.keys[n.numberOfKeys]
    n.keys[n.numberOfKeys] = newEmptyKey()

    tree.writeNode(n)
    tree.writeNode(new)
    
    return n.address, &carryUp, new.address
}

func (n *BTreeNode) insert(index int64, left int64, data *Key, right int64, tree *BTree) (int64, *Key, int64) {
    if index < 0 || index >= int64(len(n.keys)) {
        panic("B Tree insert error: Invalid index")
    }

    for i := int64(len(n.keys)) - 1; i > index; i-- {
        n.keys[i] = n.keys[i - 1]
        n.child[i + 1] = n.child[i]
    }

    n.keys[index] = *data
    n.child[index] = left
    n.child[index + 1] = right
    n.numberOfKeys++
    
    if n.numberOfKeys == int64(len(n.keys)) {
        return n.split(tree)
    }

    tree.writeNode(n)

    return -1, nil, -1
}


// ====================================== B Tree ====================================== //

func NewBTree(order int) (*BTree, error) {
    if order < 3 {
        return nil, errors.New("Invalid order")
    }

    file, _ := os.Create(BTREE_FILE)
    binary.Write(file, binary.LittleEndian, int64(8))
    root := newNode(order, 1)
    tree := &BTree {
        root: 8,
        order: order,
        file: file,
    }
    
    tree.writeNode(root)
    
    return tree, nil
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
    b.file.Seek(address, io.SeekStart)
    buf := make([]byte, b.nodeSize())
    b.file.Read(buf)

    child := make([]int64, b.order + 1)
    keys := make([]Key, b.order)

    numberOfKeys, ptr := utils.BytesToInt64(buf, 0)
    leaf, ptr := utils.BytesToInt64(buf, ptr)

    for i := 0; i < b.order - 1; i++ {
        child[i], ptr = utils.BytesToInt64(buf, ptr)
        keys[i].id, ptr = utils.BytesToInt64(buf, ptr)
        keys[i].ptr, ptr = utils.BytesToInt64(buf, ptr)
    }
    child[len(child)-2], _ = utils.BytesToInt64(buf, ptr)
    child[len(child)-1] = -1
    keys[len(keys)-1] = newEmptyKey()
    
    return &BTreeNode {
        numberOfKeys: numberOfKeys,
        leaf: leaf,
        child: child,
        keys: keys,
        address: address,
    }
}

func (b *BTree) writeNode(node *BTreeNode){
    if node.address == -1 {
        node.address, _ = b.file.Seek(0, io.SeekEnd)
    } else {
        b.file.Seek(node.address, io.SeekStart)
    }
    
    binary.Write(b.file, binary.LittleEndian, node.numberOfKeys)
    binary.Write(b.file, binary.LittleEndian, node.leaf)
    
    for i := 0; i < b.order - 1; i++ {
        binary.Write(b.file, binary.LittleEndian, node.child[i])
        binary.Write(b.file, binary.LittleEndian, node.keys[i].id)
        binary.Write(b.file, binary.LittleEndian, node.keys[i].ptr)
    }
    binary.Write(b.file, binary.LittleEndian, node.child[len(node.child)-2])
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
        b.writeNode(root)
        b.file.Seek(0, io.SeekStart)
        binary.Write(b.file, binary.LittleEndian, root.address)
        b.root = root.address
    }
}

func (b *BTree) printFile() {
    fileEnd, _ := b.file.Seek(0, io.SeekEnd)
    b.file.Seek(0, io.SeekStart)
    reader64 := make([]byte, binary.Size(int64(0)))
    
    b.file.Read(reader64)
    fileStart, _ := b.file.Seek(0, io.SeekCurrent)
    numberOfNodes := (fileEnd-fileStart)/b.nodeSize()
    tmp, _ := utils.BytesToInt64(reader64, 0)
    
    fmt.Printf("Root Address: %5X\n", tmp)
    
    for i := int64(0); i < numberOfNodes; i++ {
        currentAddress, _ := b.file.Seek(0, io.SeekCurrent)
        fmt.Printf("[%5x] || ", currentAddress)
        b.file.Read(reader64)
        tmp, _ = utils.BytesToInt64(reader64, 0)
        fmt.Printf("Size: %1d | ", tmp)
        b.file.Read(reader64)
        tmp, _ = utils.BytesToInt64(reader64, 0)
        fmt.Printf("Leaf: %1d || {", tmp)

        for i := 0; i < b.order-1; i++ {
            b.file.Read(reader64)
            tmp, _ = utils.BytesToInt64(reader64, 0)
        
            if tmp != -1 {
                fmt.Printf("[%5x] ", tmp)
            } else {
                fmt.Printf("[     ] ")
            }
        
            b.file.Read(reader64)
            tmp, _ = utils.BytesToInt64(reader64, 0)
        
            if tmp != -1 {
                fmt.Printf("%3d ", tmp)
            } else {
                fmt.Printf("    ")
            }
        
            b.file.Read(reader64)
        }
        
        b.file.Read(reader64)
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
    order := 8
    tree, _ := NewBTree(order)
    reader, err := inicializarControleLeitura(BIN_FILE)
    
    for i := 0; i < int(reader.TotalRegistros) && err == nil; i++ {
		err = reader.ReadNext()
		if reader.RegistroAtual.Lapide != 1 {
			r := newKey(reader.RegistroAtual)
			tree.Insert(&r)

            fmt.Printf("Id adicionado = %d\n", r.id)
		}
	}
    
    tree.printFile()
    // fmt.Printf("%+v\n", tree.readNode(tree.root))
}
