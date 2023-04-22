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

const BTREE_FILE string = "btree/BTree.dat"
const BTREE_NODES_FILE string = "btree/BTreeNodes.dat"
const NULL int64 = -1

// TODO:
// - change dir concatenation to use os.join
// - mkdir

// ====================================== Bit-Flags ====================================== //

// Bit-Flags used for removing an element from the B Tree
const (
    // Value removed without any complications
    FLAG1 = 1 << 0

    // Node has less than 50% ocupation
    FLAG2 = 1 << 1

    // Value is in a non-leaf node
    FLAG3 = 1 << 2

    // Empty node
    FLAG4 = 1 << 3
)

// ====================================== Structs ====================================== //

// Key contem os valores que a árvore carrega
type Key struct {
    Id  int64
    Ptr int64
}

// BTreeNode é o struct que representa cada nó que
// está presente na árvore
type BTreeNode struct {
    address      int64
    numberOfKeys int64
    child        []int64
    keys         []Key
    leaf         int64
}

// BTree é.. a Árvore B
type BTree struct {
    file       string
    nodesFile  *os.File
    root       int64
    order      int
    emptyNodes []int64
}

// ====================================== Key ====================================== //

// newKey inicializa uma nova chave a partir de um struct Registro
func newKey(register *binManager.Registro) Key {
    return Key{
        Id:  int64(register.Pokemon.Numero),
        Ptr: register.Endereco,
    }
}

// newEmptyKey inicializa uma chave vazia
func newEmptyKey() Key {
    return Key{NULL, NULL}
}

// ====================================== Node ====================================== //

// newNode inicializa um novo no' vazio, preenchendo os campos com NULL
// ou com uma chave vazia (vide funcao acima)
func newNode(order int, leaf int64, address int64) *BTreeNode {
    node := BTreeNode{
        child:        make([]int64, order+1),
        keys:         make([]Key, order),
        numberOfKeys: 0,
        leaf:         leaf,
        address:      address,
    }

    for i := 0; i < order; i++ {
        node.child[i] = NULL
        node.keys[i] = newEmptyKey()
    }
    node.child[len(node.child)-1] = NULL

    return &node
}

// write escreve o no' no arquivo, diretamente no endereco do no',
// se existente, caso contrario, escreve no final do arquivo
func (n *BTreeNode) write(file *os.File) {
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
}

// self: * 2 l 3 r 5 * 9 *
//
//    |
//    v
//
// self: * 2 l _ * _ * _ *
// new:  r 5 * 9 * _ * _ *
//
// return (self, 3, new)
func (n *BTreeNode) split(tree *BTree) (int64, *Key, int64) {
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

    n.numberOfKeys = int64(middle - 1)
    carryUp := n.keys[n.numberOfKeys]
    n.keys[n.numberOfKeys] = newEmptyKey()

    n.write(tree.nodesFile)
    new.write(tree.nodesFile)

    return n.address, &carryUp, new.address
}

// insert insere uma nova chave no no', recebendo a chave a ser inserida
// e o indice de onde o valor vai ser inserido.
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

    return NULL, nil, NULL
}

// canLendKey testa se um no' tem chaves o suficiente para emprestar
// para o no' irmao, durante a remocao
func (n *BTreeNode) canLendKey() bool {
    return n.numberOfKeys >= int64(len(n.keys)/2)
}

// removeKeyLeaf remove um elemento de um no' que for uma folha, retornando
// o elemento removido e um ponteiro para o filho a esquerda, que sera' usado
// durante a remocao
func (n *BTreeNode) removeKeyLeaf(index int64) (*Key, int64) {
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

// removeKeyNonLeaf remove um elemento de um no' nao folha
func (n *BTreeNode) removeKeyNonLeaf(index int64) *Key {
    if index == NULL {
        index = n.numberOfKeys - 1
    }

    return &n.keys[index]
}

// max retorna o maior elemento presente em um no'
func (n *BTreeNode) max() Key {
    return n.keys[n.numberOfKeys-1]
}

// find procura um elemento no no', retornando o elemento, se encontrado,
// ou um ponteiro para o proximo no' a ser procurado, caso nao encontrado
func (n *BTreeNode) find(id int64) (*Key, int64) {
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

// update atualiza o valor de uma chave, pesquisando por ela e 
// alterando seu valor
func (n *BTreeNode) update(newKey Key) {
    i := int64(0)
    
    for i < n.numberOfKeys && n.keys[i].Id < newKey.Id {
        i++
    }

    if n.keys[i].Id == newKey.Id {
        n.keys[i].Ptr = newKey.Ptr
    } 
}


// ====================================== B Tree ====================================== //

// NewBTree inicializa uma arvore vazia, recebendo o endereco do
// arquivo a ser gravado e a ordem da arvore
func NewBTree(order int, dir string) (*BTree, error) {
    if order < 3 {
        return nil, errors.New("invalid order")
    }

    nodesFile, _ := os.Create(dir + BTREE_NODES_FILE)
    root := newNode(order, 1, NULL)
    tree := &BTree{
        root:      0,
        order:     order,
        file:      dir + BTREE_FILE,
        nodesFile: nodesFile,
    }

    root.write(nodesFile)

    return tree, nil
}

// ReadBTree lê uma arvore de um arquivo header, extraindo
// a ordem da arvore, o endereco da raiz, a quantidade de nós 
// vazios e os endereços dos nós vazios
func ReadBTree(dir string) (*BTree, error) {
    file, err := os.ReadFile(dir + BTREE_FILE)
    if err != nil {
        return nil, err
    }
    
    nodesFile, _ := os.OpenFile(dir+BTREE_NODES_FILE, os.O_RDWR|os.O_CREATE, 0644)
    root, ptr := utils.BytesToInt64(file, 0)
    order, ptr := utils.BytesToInt64(file, ptr)
    len, ptr := utils.BytesToInt64(file, ptr)
    emptyNodes := make([]int64, len)

    for i := int64(0); i < len; i++ {
        emptyNodes[i], ptr = utils.BytesToInt64(file, ptr)
    }

    return &BTree{
        root:      root,
        order:     int(order),
        file:      dir + BTREE_FILE,
        nodesFile: nodesFile,
    }, nil
}

// Close fecha o arquivo da árvore, salvando o endereço da raiz, 
// a ordem, quantidade de nós vazios e os endereços dos nós vazios
// (Esta função deve ser chamada para salvar qualquer alteração feita
// na base de dados. Nao cumprimento disso poderá ocasionar em dados
// corrompidos)
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

// popEmptyNode busca um endereço de um nó vazio, o remove
// da lista de nós vazios e por fim o retorna
func (b *BTree) popEmptyNode() int64 {
    if len(b.emptyNodes) > 0 {
        address := b.emptyNodes[0]
        b.emptyNodes = b.emptyNodes[1:]
        return address
    }

    return NULL
}

// nodeSize calcula o tamanho de um nó em bytes para ser salvo
// em um arquivo
func (b *BTree) nodeSize() int64 {
    node := BTreeNode{}
    s := int64(0)
    s += int64(binary.Size(node.numberOfKeys))
    s += int64(binary.Size(node.leaf))
    s += int64(binary.Size(int64(0)) * b.order)
    s += int64(binary.Size(Key{}) * (b.order - 1))
    return s
}

// readNode lê um nó do arquivo dado seu endereço
func (b *BTree) readNode(address int64) *BTreeNode {
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
    child[len(child)-2], _ = utils.BytesToInt64(buf, ptr)
    child[len(child)-1] = NULL
    keys[len(keys)-1] = newEmptyKey()

    return &BTreeNode{
        numberOfKeys: numberOfKeys,
        leaf:         leaf,
        child:        child,
        keys:         keys,
        address:      address,
    }
}

// insert insere uma chave na árvore e atualiza o arquivo
// com os nós 
func (b *BTree) insert(node *BTreeNode, data *Key) (int64, *Key, int64) {
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

// Insert insere uma chave na árvore e atualiza o arquivo
// com os nós
func (b *BTree) Insert(data *Key) {
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

// printFile abre o arquivo com os nós e printa todos eles
// na ordem que aparecem (ideal para debug)
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

            if tmp != NULL {
                fmt.Printf("[%5x] ", tmp)
            } else {
                fmt.Printf("[     ] ")
            }

            b.nodesFile.Read(reader64)
            tmp, _ = utils.BytesToInt64(reader64, 0)

            if tmp != NULL {
                fmt.Printf("%3d ", tmp)
            } else {
                fmt.Printf("    ")
            }

            b.nodesFile.Read(reader64)
        }

        b.nodesFile.Read(reader64)
        tmp, _ = utils.BytesToInt64(reader64, 0)

        if tmp != NULL {
            fmt.Printf("[%4x] }\n", tmp)
        } else {
            fmt.Printf("[    ] }\n")
        }
    }
    fmt.Printf("\n\n")
}

// maxLeft procura o maior elemento a esquerda de uma chave,
// para substituir a chave que está sendo removida
func (b *BTree) maxLeft(node *BTreeNode, index int64, k *Key) Key {
    node = b.readNode(node.child[index])

    for node.leaf == 0 {
        node = b.readNode(node.child[node.numberOfKeys])
    }

    return node.max()
}

// replace substitui um elemento da árvore por outro, procurando
// na subárvore inteira abaixo do nó enviado para a função
func (b *BTree) replace(new *Key, old *Key, node *BTreeNode) {
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

// concatNodes junta dois nós que estao com o tamanho pequeno
// em um novo nó, retornando o nó resultante e enviando o nó
// perdido para os nós vazios salvos na árvore
func (b *BTree) concatNodes(left *BTreeNode, right *BTreeNode, key *Key) *BTreeNode {
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

    left.write(b.nodesFile)
    right.write(b.nodesFile)

    return left
}

// borrowFromParent desce um elemento do nó pai que será usado
// para concatenar os nós filhos durante a remoção, retornando
// uma flag indicando se ocorreu algum problema durante o
// processo
func (b *BTree) borrowFromParent(node *BTreeNode, l *BTreeNode, r *BTreeNode, keyIndex int64) int {
    flag := 0

    if node.numberOfKeys == 1 {
        flag = FLAG4
    } else if node.canLendKey() {
        flag = FLAG1
    } else {
        flag = FLAG2
    }

    b.concatNodes(l, r, &node.keys[keyIndex])

    for i := keyIndex; i < node.numberOfKeys; {
        node.keys[i] = node.keys[i+1]
        i++
        node.child[i] = node.child[i+1]
    }

    b.emptyNodes = append(b.emptyNodes, r.address)

    node.numberOfKeys--
    node.write(b.nodesFile)

    return flag
}

// borrowFromSibling busca um elemento de um nó irmão, o
// envia para o nó pai, desce o elemento do pai e envia para
// o nó que esta com menos de 50% de ocupação
func (b *BTree) borrowFromSibling(parent *BTreeNode, left *BTreeNode, right *BTreeNode, index int64) {
    k, child := right.removeKeyLeaf(0)

    left.keys[left.numberOfKeys] = parent.keys[index]
    left.numberOfKeys++
    left.child[left.numberOfKeys] = child

    parent.keys[index] = *k

    parent.write(b.nodesFile)
    left.write(b.nodesFile)
    right.write(b.nodesFile)
}

// tryBorrowKey testa se (durante a remoção) a chave pode 
// ser emprestada de um nó irmão ou se este vai ficar com
// menos de 50% de ocupação, se pode, pega a chave do irmão,
// se nao, pega a chave do pai e concatena os irmãos, retornando
// uma flag indicando se teve erro no processo
func (b *BTree) tryBorrowKey(node *BTreeNode, index int64) int {
    var l, r, lefter *BTreeNode

    l = b.readNode(node.child[index])
    r = b.readNode(node.child[index+1])

    if index != 0 {
        lefter = b.readNode(node.child[index-1])
    }

    if r.canLendKey() {
        b.borrowFromSibling(node, l, r, index)
    } else if lefter != nil && lefter.canLendKey() {
        b.borrowFromSibling(node, lefter, l, index)
    } else {
        return b.borrowFromParent(node, l, r, index)
    }

    return FLAG1
}

// parseFlag resolve qualquer problema que as flags podem
// indicar que ocorreu, chamando as respectivas funções que 
// cuidam desses cenários e retornando uma nova flag que 
// será resolvida no nó anterior na recursão
func (b *BTree) parseFlag(flag int, node *BTreeNode, index int64, k *Key) int {
    if FLAG1 & flag != 0 {
        // nothing to see here
    } else if FLAG2 & flag != 0 {
        flag = b.tryBorrowKey(node, index)
    } else if FLAG3 & flag != 0 {
        maxLeft := b.maxLeft(node, index, k)
        _, flag, _ = b.remove(node.address, maxLeft.Id)
        b.replace(&maxLeft, k, node)
    }

    return flag
}

// removeFromNode remove um elemento que está presente
// em um nó, retornando a flag correspondente a o que
// aconteceu durante a remoção
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
        k, _ = node.removeKeyLeaf(index)
        node.write(b.nodesFile)
    } else {
        k = node.removeKeyNonLeaf(index)
    }

    return k, flag
}

// remove remove um elemento da árvore e o retorna,
// ou NULL, caso nao encontrado. Retorna o elemento,
// uma flag e o indice do elemento removido no nó
func (b *BTree) remove(address int64, id int64) (*Key, int, int64) {
    if address == NULL {
        return nil, FLAG1, NULL
    }

    var k *Key
    node := b.readNode(address)
    flag := 0

    i := node.numberOfKeys - 1
    for i > 0 && node.keys[i].Id > id {
        i--
    }

    if node.keys[i].Id == id {
        k, flag = b.removeFromNode(i, node)

        if node.leaf == 1 {
            return k, flag, i
        }

        maxLeft := b.maxLeft(node, i, k)
        _, flag, _ = b.remove(address, maxLeft.Id)

        node = b.readNode(node.address)
        b.replace(&maxLeft, k, node)

        return k, flag, i
    } else if node.keys[i].Id < id {
        k, flag, _ = b.remove(node.child[i+1], id)
    } else {
        k, flag, _ = b.remove(node.child[i], id)
    }

    flag = b.parseFlag(flag, node, i, k)

    return k, flag, i
}

// Remove remove um elemento da árvore, pesquisando
// recursivamente pelo elemento e por fim retornando-o,
// ou nil, caso não encontrado
func (b *BTree) Remove(id int64) *Key {
    k, flag, index := b.remove(b.root, id)
    var root *BTreeNode

    if FLAG3 & flag != 0 {
        root = b.readNode(b.root)
        node := b.readNode(root.child[index])
        maxLeft := node.max()
        b.remove(root.address, maxLeft.Id)
        root = b.readNode(b.root)
        b.replace(&maxLeft, k, root)
    } else if FLAG4 & flag != 0 {
        root = b.readNode(b.root)
        tmp := root.child[0]
        root.child[0] = NULL
        root.write(b.nodesFile)
        b.root = tmp
        root = b.readNode(b.root)
        root.write(b.nodesFile)
    }

    return k
}

// Find pesquisa por um elemento presente na árvore
// e o retorna se encontrado, ou nil, caso contrário
func (b *BTree) Find(id int64) *Key {
    var k *Key
    var address int64
    node := b.readNode(b.root)
    
    for node != nil && k == nil {
        k, address = node.find(id)
        node = b.readNode(address)
    }
    
    return k
}

// Update atualiza um elemento na árvore, pesquisando-o recursivamente
// e, se encontrado, atualiza seu valor
func (b *BTree) Update(id int64, ptr int64) {
    var k *Key
    var address int64
    node := b.readNode(b.root)
    k, address = node.find(id)

    for node != nil && k == nil {
        node = b.readNode(address)
        k, address = node.find(id)
    }

    if node != nil {
        node.update(Key {id, ptr})
        node.write(b.nodesFile)
    }
}

// StartBTreeFile inicializa a Árvore B e insere todos os elementos
// contidos no arquivo informado na árvore, escrevendo-os em um novo arquivo
// e escrevendo as informações gerais da arvore (como ordem, endereço da raíz,
// etc) em outro arquivo
func StartBTreeFile(dir string) {
    order := 8
    tree, _ := NewBTree(order, dir)
    reader, err := binManager.InicializarControleLeitura(binManager.BIN_FILE)

    n := int(reader.TotalRegistros)

    for i := 0; i < n && err == nil; i++ {
        err = reader.ReadNext()
        if reader.RegistroAtual.Lapide != 1 {
            r := newKey(reader.RegistroAtual)
            tree.Insert(&r)
        }
    }

    tree.Close()
}
