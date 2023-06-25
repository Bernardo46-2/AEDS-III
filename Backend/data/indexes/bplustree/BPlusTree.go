// Package bplustree implementa uma árvore B+ em memoria secundaria, uma estrutura
// de dados em árvore que é usada para armazenar dados ordenados de maneira eficiente.
// A árvore B+ é uma versão especializada da árvore B que permite inserções, exclusões
// e pesquisas rápidas.
//
// As árvores B+ são comumente usadas em sistemas de banco de dados e sistemas de
// arquivos devido à sua eficiência em lidar com grandes quantidades de dados e
// sua habilidade de manter os dados ordenados, tornando a busca de intervalos e a
// varredura sequencial mais rápidas.
//
// O pacote oferece funções para inserir, excluir e procurar dados na árvore,
// verificar a integridade da árvore e visualizar a árvore em uma representação
// gráfica.
//
// Este pacote não fornece suporte a transações ou persistência, embora esses
// recursos possam ser adicionados se necessário.
//
// O desempenho do pacote bplustree pode variar dependendo do hardware e da carga
// de trabalho, mas em geral, é esperado que ofereça desempenho similar ou
// superior a outras estruturas de dados ordenadas, como árvores AVL e árvores
// vermelho-preto, para grandes conjuntos de dados.
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

// Path dos arquivos necessarios
const (
	PATH   string = "bplustree/"
	NODES  string = "BPlusTreeNodes.bin"
	HEADER string = "BPlusTree.bin"
)

// NULL padrao
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

// Key contem os valores que vão estar presente
// na árvore
type Key struct {
	Id  float64
	Ptr int64
}

// BPlusTreeNode representa cada nó que contém
// a árvore
type BPlusTreeNode struct {
	address      int64
	numberOfKeys int64
	child        []int64
	keys         []Key
	leaf         int64
	next         int64
}

// BPlusTree é.. a Árvore B+
type BPlusTree struct {
	file       string
	nodesFile  *os.File
	root       int64
	order      int
	emptyNodes []int64
}

// Interface para leitura da database
type Reader interface {
	ReadNextGeneric() (any, bool, int64, error)
}

// Interface para recuperacao do campo do objeto indexavel
type IndexableObject interface {
	GetFieldF64(fieldName string) (float64, int64)
}

// ====================================== Key ====================================== //

// newEmptyKey inicializa uma chave vazia
func newEmptyKey() Key {
	return Key{float64(NULL), NULL}
}

// compareTo compara uma chave à outra, podendo
// retornar valores que representam maior, menor ou igual
func (k *Key) compareTo(other *Key) float64 {
	diff := k.Id - other.Id
	if diff == 0 {
		diff = float64(k.Ptr - other.Ptr)
	}
	return diff
}

// ====================================== Node ====================================== //

// newNode inicializa um novo no' vazio, preenchendo os campos com NULL
// ou com uma chave vazia
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

// write escreve o no' no arquivo, diretamente no endereco do no',
// se existente, caso contrario, escreve no final do arquivo
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

// insert insere uma nova chave no no', recebendo a chave a ser inserida
// e o indice de onde o valor vai ser inserido.
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

// find pesquisa no nó pela presença de uma
// chave, retornando o indice onde ela está,
// ou NULL
func (n *BPlusTreeNode) find(k *Key) int64 {
	i := n.numberOfKeys - 1

	for i >= 0 && n.keys[i].compareTo(k) > 0 {
		i--
	}

	if n.keys[i].compareTo(k) != 0 {
		i = NULL
	}

	return i
}

// canLendKey testa se um no' tem chaves o suficiente para emprestar
// para o no' irmao, durante a remocao
func (n *BPlusTreeNode) canLendKey() bool {
	return n.numberOfKeys >= int64(len(n.keys)/2)
}

// removeKeyLeaf remove um elemento de um no' que for uma folha, retornando
// o elemento removido e um ponteiro para o filho a esquerda, que sera' usado
// durante a remocao
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

// max retorna o maior elemento presente em um no'
func (n *BPlusTreeNode) max() *Key {
	if n.numberOfKeys == 0 {
		return nil
	}
	return &n.keys[n.numberOfKeys-1]
}

// getStatus testa se o nó é uma folha e se
// pode emprestar uma chave e retorna uma flag
// indicando as duas coisas
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

// NewBPlusTree inicializa uma arvore vazia, recebendo o endereco do
// arquivo a ser gravado, a ordem da arvore e o campo que vai ser
// guardado
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

// ReadBPlusTree lê uma arvore de um arquivo header, extraindo
// a ordem da arvore, o endereco da raiz, a quantidade de nós
// vazios e os endereços dos nós vazios
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

// Close fecha o arquivo da árvore, salvando o endereço da raiz,
// a ordem, quantidade de nós vazios e os endereços dos nós vazios
// (Esta função deve ser chamada para salvar qualquer alteração feita
// na base de dados. Nao cumprimento disso poderá ocasionar em dados
// corrompidos)
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

// popEmptyNode busca um endereço de um nó vazio, o remove
// da lista de nós vazios e por fim o retorna
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

// nodeSize calcula o tamanho de um nó em bytes para ser salvo
// em um arquivo
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

// readNode lê um nó do arquivo dado seu endereço
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

// insert insere uma chave na árvore e atualiza o arquivo
// com os nós
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

// Insert insere uma chave na árvore e atualiza o arquivo
// com os nós
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

// printFile abre o arquivo com os nós e printa todos eles
// na ordem que aparecem (ideal para debug)
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
				fmt.Printf("%5d ", i64)
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

// concatLeaf concatena dois nós em um, dado que estes são folhas
// e retorna o nó resultante
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

// mergeLeaf chama a função concatLeaf para concatenar dois nós e
// desloca os elementos do pai para comportar 1 elemento a menos
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

// borrowFromLeaf pega um elemento emprestado do nó irmao (folha)
// e faz todos os ajustes para passar pro nó com tamanho
// invalido
func (b *BPlusTree) borrowFromLeaf(node *BPlusTreeNode, left *BPlusTreeNode, right *BPlusTreeNode, index int64) {
	k, _ := right.removeKeyLeaf(0)

	left.keys[left.numberOfKeys] = *k
	left.numberOfKeys++
	node.keys[index] = *k

	node.write(b.nodesFile)
	left.write(b.nodesFile)
	right.write(b.nodesFile)
}

// borrowFromNonLeaf pega um elemento emprestado do nó irmao (nao folha)
// e faz todos os ajustes para passar pro nó com tamanho invalido
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

// concatNodesOG concatena dois nós não folha e retorna o
// nó resultante
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

// borrowFromParentOG pega um valor do pai e usa para concatenar
// os dois nós irmãos
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

// tryBorrowKey testa se (durante a remoção) a chave pode
// ser emprestada de um nó irmão ou se este vai ficar com
// menos de 50% de ocupação, se pode, pega a chave do irmão,
// se nao, pega a chave do pai e concatena os irmãos, retornando
// uma flag indicando se teve erro no processo
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

// replaceKey testa se a chave removida tem sua cópia
// no nó atual e, se sim, substitui ela pela nova
// chave (usado durante a remoção) e retorna uma nova
// flag
func (b *BPlusTree) replaceKey(n *BPlusTreeNode, k *Key, kk *Key, flag int) int {
	walter := 0
	if REPLACE&flag != 0 {
		i := n.find(k)
		if i != NULL {
			n.keys[i] = *kk
			n.write(b.nodesFile)
		} else {
			walter = REPLACE
		}
	}

	return walter
}

// parseFlag resolve qualquer problema que as flags podem
// indicar que ocorreu, chamando as respectivas funções que
// cuidam desses cenários e retornando uma nova flag que
// será resolvida no nó anterior na recursão
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

// removeFromNode remove um elemento que está presente
// em um nó, retornando a flag correspondente a o que
// aconteceu durante a remoção
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

// remove remove um elemento da árvore e o retorna,
// ou NULL, caso nao encontrado. Retorna o elemento,
// uma flag e o indice do elemento removido no nó
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

// Remove remove um elemento da árvore, pesquisando
// recursivamente pelo elemento e por fim retornando-o,
// ou nil, caso não encontrado
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

// find2 procura um elemento no no', retornando o elemento, se encontrado,
// ou um ponteiro para o proximo no' a ser procurado, caso nao encontrado
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

// Find pesquisa por um elemento presente na árvore
// e o retorna se encontrado, ou nil, caso contrário
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

// findNode pesquisa na árvore em busca de um nó que tem
// o elemento procurado e retorna esse nó ou o nó com o proximo
// elemento ordenadamente
func (b *BPlusTree) findNode(id float64) *BPlusTreeNode {
	node := b.readNode(b.root)

	for i := int64(0); node.leaf == 0; i++ {
		if i < node.numberOfKeys-1 {
			if node.keys[i].Id >= id {
				node = b.readNode(node.child[i])
				i = NULL
			}
		} else {
			if node.keys[i].Id >= id {
				node = b.readNode(node.child[i])
			} else {
				node = b.readNode(node.child[i+1])
			}
			i = NULL
		}
	}

	return node
}

// FindRange pesquisa na árvore por todos os valores contidos em um
// intervalo [a, b) e os retorna
func (b *BPlusTree) FindRange(start float64, end float64) ([]int64, error) {
	if start > end || end < 1 {
		return nil, nil
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

	for start <= end && node != nil {

		start = node.keys[index].Id
		if start > end {
			break
		}
		addresses = append(addresses, node.keys[index].Ptr)

		if index == node.numberOfKeys-1 {
			node = b.readNode(node.next)
			index = NULL
		}

		index++
	}

	return addresses, nil
}

// Create insere um elemento na árvore
func Create(pokemon models.Pokemon, pokeAddress int64, path string, fields []string) {
	for _, field := range fields {
		tree, _ := ReadBPlusTree(path, field)
		id, _ := pokemon.GetFieldF64(field)
		k := Key{id, pokeAddress}
		tree.Insert(&k)
		tree.Close()
	}
}

// Update atualiza um elemento da árvore
func Update(old models.Pokemon, new models.Pokemon, kAddress int64, kkAddress int64, path string, fields []string) {
	for _, field := range fields {
		tree, _ := ReadBPlusTree(path, field)
		kId, _ := old.GetFieldF64(field)
		kkId, _ := new.GetFieldF64(field)
		k := Key{kId, kAddress}
		kk := Key{kkId, kkAddress}
		tree.Remove(&k)
		tree.Insert(&kk)
		tree.Close()
	}
}

// Delete remove um elemento da árvore
func Delete(pokemon models.Pokemon, address int64, path string, fields []string) {
	for _, field := range fields {
		tree, _ := ReadBPlusTree(path, field)
		id, _ := pokemon.GetFieldF64(field)
		tree.Remove(&Key{id, address})

		tree.Close()
	}
}

// ====================================== Init ====================================== //

// StartBPlusTreeFile inicializa a Árvore B e insere todos os elementos
// contidos no arquivo informado na árvore, escrevendo-os em um novo arquivo
// e escrevendo as informações gerais da arvore (como ordem, endereço da raíz,
// etc) em outro arquivo
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

// StartBPlusTreeFile inicializa a Árvore B e insere todos os elementos
// contidos no arquivo informado na árvore, escrevendo-os em um novo arquivo
// e escrevendo as informações gerais da arvore (como ordem, endereço da raíz,
// etc) em outro arquivo
func StartBPlusTreeFilesSearch(dir string, field string, controler Reader) error {
	order := 8
	tree, _ := NewBPlusTree(order, dir, field)

	for {
		objInterface, isDead, address, err := controler.ReadNextGeneric()
		if err != nil {
			break
		}

		obj, ok := objInterface.(IndexableObject)
		if !ok {
			return fmt.Errorf("failed to convert object to IndexableObject\n%+v", objInterface)
		}

		if !isDead {
			id, _ := obj.GetFieldF64(field)
			r := Key{id, address}
			tree.Insert(&r)
		}
	}

	tree.Close()
	return nil
}
