package bplustree

import (
    "encoding/binary"
    "errors"
    "fmt"
    "io"
    "os"

    "github.com/Bernardo46-2/AEDS-III/data/binManager"
    "github.com/Bernardo46-2/AEDS-III/utils"
)

const B_PLUS_TREE_FILE string = "bplustree/BPlusTree.dat"
const B_PLUS_TREE_NODES_FILE string = "bplustree/BPlusTreeNodes.dat"
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
    Id  int64
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

// ====================================== Key ====================================== //

func newKey(register *binManager.Registro) Key {
    return Key{
        Id:  int64(register.Pokemon.Numero),
        Ptr: register.Endereco,
    }
}

func newEmptyKey() Key {
    return Key{NULL, NULL}
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
//    |
//    v
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
        carryUp = n.keys[n.numberOfKeys - 1]
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

func (n *BPlusTreeNode) find(id int64) int64 {
    i := n.numberOfKeys - 1
    
    for i >= 0 && n.keys[i].Id > id {
        i--
    }

    return i
}

func (n *BPlusTreeNode) update(newKey Key) {
    i := int64(0)
    
    for i < n.numberOfKeys && n.keys[i].Id < newKey.Id {
        i++
    }

    if n.keys[i].Id == newKey.Id {
        n.keys[i].Ptr = newKey.Ptr
    } 
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

func (n *BPlusTreeNode) removeKeyNonLeaf(index int64) *Key {
    if index == NULL {
        index = n.numberOfKeys - 1
    }

    return &n.keys[index]
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

func NewBPlusTree(order int, dir string) (*BPlusTree, error) {
    if order < 3 {
        return nil, errors.New("invalid order")
    }

    nodesFile, _ := os.Create(dir + B_PLUS_TREE_NODES_FILE)
    root := newNode(order, 1, NULL)
    tree := &BPlusTree{
        root:      0,
        order:     order,
        file:      dir + B_PLUS_TREE_FILE,
        nodesFile: nodesFile,
        emptyNodes: make([]int64, 0),
    }

    root.write(nodesFile)

    return tree, nil
}

func ReadBPlusTree(dir string) (*BPlusTree, error) {
    file, err := os.ReadFile(dir + B_PLUS_TREE_FILE)
    if err != nil {
        return nil, err
    }
    
    nodesFile, _ := os.OpenFile(dir+B_PLUS_TREE_NODES_FILE, os.O_RDWR|os.O_CREATE, 0644)
    root, ptr := utils.BytesToInt64(file, 0)
    order, ptr := utils.BytesToInt64(file, ptr)
    len, ptr := utils.BytesToInt64(file, ptr)
    emptyNodes := make([]int64, len)

    for i := int64(0); i < len; i++ {
        emptyNodes[i], ptr = utils.BytesToInt64(file, ptr)
    }

    return &BPlusTree{
        root:      root,
        order:     int(order),
        file:      dir + B_PLUS_TREE_FILE,
        nodesFile: nodesFile,
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
        keys[i].Id, ptr = utils.BytesToInt64(buf, ptr)
        keys[i].Ptr, ptr = utils.BytesToInt64(buf, ptr)
    }
    child[len(child)-2], ptr = utils.BytesToInt64(buf, ptr)
    child[len(child)-1] = NULL
    keys[len(keys)-1] = newEmptyKey()

    next, _ := utils.BytesToInt64(buf, ptr);

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
    for i < node.numberOfKeys && data.Id > node.keys[i].Id {
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

func (b *BPlusTree) printFile() {
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
            fmt.Printf("[%4x] ", tmp)
        } else {
            fmt.Printf("[    ] ")
        }

        b.nodesFile.Read(reader64)
        tmp, _ = utils.BytesToInt64(reader64, 0)

        if tmp != -1 {
            fmt.Printf("|| Next: %4x }\n", tmp)
        } else {
            fmt.Printf("|| Next: %4x }\n", NULL)
        }
    }
    fmt.Printf("\n\n")
}

func (b *BPlusTree) maxLeft(node *BPlusTreeNode, index int64, k *Key) *Key {
    node = b.readNode(node.child[index])

    for node.leaf == 0 {
        node = b.readNode(node.child[node.numberOfKeys])
    }

    return node.max()
}

func (b *BPlusTree) replace(new *Key, old *Key, node *BPlusTreeNode) {
    i := node.numberOfKeys - 1
    for i > 0 && node.keys[i].Id > old.Id {
        i--
    }

    if node.keys[i].Id == old.Id {
        node.keys[i] = *new
        node.write(b.nodesFile)
    } else if old.Id < node.keys[i].Id {
        b.replace(new, old, b.readNode(node.child[i]))
    } else {
        b.replace(new, old, b.readNode(node.child[i]))
    }
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

    if flag == LEAF | CAN_LEND { // replace
        b.borrowFromLeaf(node, l, r, index)
    } else if flag == LEAF { // merge
        b.mergeLeaf(node, l, r, index)
    } else if flag == CAN_LEND { // borrow B
        b.borrowFromNonLeaf(node, l, r, index);
    } else { // downgrade parent
        b.borrowFromParentOG(node, l, r, index);
    }

    if node.canLendKey() {
        flag = OK
    } else {
        flag = HELP
    }

    return flag
}

func (b *BPlusTree) replaceKey(n *BPlusTreeNode, k *Key, kk *Key) (flag int) {
    if REPLACE & flag != 0 {
        i := n.find(k.Id);
        if i != NULL {
            n.keys[i] = *kk
            n.write(b.nodesFile)
        } else {
            flag = REPLACE
        }
    }

    return
}

func (b *BPlusTree) parseFlag(flag int, node *BPlusTreeNode, index int64, k *Key, kk *Key) int {
    walter := b.replaceKey(node, k, kk)

    if HELP & flag != 0 {
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

    if index == node.numberOfKeys - 1 {
        flag |= REPLACE
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

func (b *BPlusTree) remove(address int64, id int64) (*Key, *Key, int, int64) {
    if address == NULL {
        return nil, nil, OK, NULL
    }

    var k, kk *Key
    node := b.readNode(address)
    flag := 0

    i := node.numberOfKeys - 1
    for i > 0 && node.keys[i].Id > id {
        i--
    }

    if node.leaf == 1 && node.keys[i].Id == id {
        k, kk, flag = b.removeFromNode(i, node)
    } else if node.keys[i].Id < id {
        k, kk, flag, _ = b.remove(node.child[i+1], id)
        flag = b.parseFlag(flag, node, i, k, kk)
    } else {
        k, kk, flag, _ = b.remove(node.child[i], id)
        flag = b.parseFlag(flag, node, i, k, kk)
    }


    return k, kk, flag, i
}

func (b *BPlusTree) Remove(id int64) *Key {
    k, kk, flag, _ := b.remove(b.root, id)

    root := b.readNode(b.root)

    if flag & REPLACE != 0 {
        b.replaceKey(root, k, kk)
    } 

    if flag & EMPTY != 0 {
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


// ====================================== Tests ====================================== //

func StartBPlusTreeFile(dir string) {
    order := 8
    tree, err := NewBPlusTree(order, dir)
    reader, err := binManager.InicializarControleLeitura(binManager.BIN_FILE)

    n := int(reader.TotalRegistros)
    for i := 0; i < n && err == nil; i++ {
        err = reader.ReadNext()
        if reader.RegistroAtual.Lapide != 1 {
            r := newKey(reader.RegistroAtual)
            tree.Insert(&r)
        }
    }

    for i := 1; i <= n; i++ {
        tree.Remove(int64(i))
    }

    reader.Close()
    reader2, err2 := binManager.InicializarControleLeitura(binManager.BIN_FILE)
    for i := 0; i < n && err2 == nil; i++ {
        err = reader2.ReadNext()
        if reader2.RegistroAtual.Lapide != 1 {
            g := newKey(reader2.RegistroAtual)
            tree.Insert(&g)
        }
    }
    
    tree.printFile()
    tree.Close()
}
