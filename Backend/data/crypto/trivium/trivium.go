package trivium

import (
	"bufio"
	"crypto/rand"
	"log"
	"os"
)

// Trivium representa o estado de 288 bits da cifra Trivium.
type Trivium struct {
	state [5]uint64
	Key   [KeyLength]byte
	iv    [KeyLength]byte
}

// Definições de constantes para o funcionamento da cifra Trivium.
const (
	// KeyLength bytes na chave e no valor de inicialização (IV), 10 bytes = 80 bits
	KeyLength = 10

	// lgWordSize é a base do logaritmo do tamanho da palavra usada na matriz de respaldo.
	lgWordSize = 6

	// Os índices na matriz correspondentes às células que são acionadas para o processamento.
	i66  = 65 >> lgWordSize
	i93  = 92 >> lgWordSize
	i162 = 161 >> lgWordSize
	i177 = 176 >> lgWordSize
	i243 = 242 >> lgWordSize
	i288 = 287 >> lgWordSize
	i91  = 90 >> lgWordSize
	i92  = 91 >> lgWordSize
	i171 = 170 >> lgWordSize
	i175 = 174 >> lgWordSize
	i176 = 175 >> lgWordSize
	i264 = 263 >> lgWordSize
	i286 = 285 >> lgWordSize
	i287 = 286 >> lgWordSize
	i69  = 68 >> lgWordSize
	i94  = 93 >> lgWordSize
	i178 = 177 >> lgWordSize

	// A posição dentro da palavra, a mudança dentro da palavra começando da esquerda.
	wordSize = 1 << lgWordSize
	mask     = wordSize - 1

	// Os shifts correspondentes a cada célula acionada.
	sh66  = mask - (65 & mask)
	sh93  = mask - (92 & mask)
	sh162 = mask - (161 & mask)
	sh177 = mask - (176 & mask)
	sh243 = mask - (242 & mask)
	sh288 = mask - (287 & mask)
	sh91  = mask - (90 & mask)
	sh92  = mask - (91 & mask)
	sh171 = mask - (170 & mask)
	sh175 = mask - (174 & mask)
	sh176 = mask - (175 & mask)
	sh264 = mask - (263 & mask)
	sh286 = mask - (285 & mask)
	sh287 = mask - (286 & mask)
	sh69  = mask - (68 & mask)
	sh94  = mask - (93 & mask)
	sh178 = mask - (177 & mask)
)

// nextBit obtém o próximo bit do fluxo Trivium.
func (t *Trivium) nextBit() uint64 {
	return t.nextBits(1)
}

// nextBits obtém os próximos 1 a 63 bits do Trivium Stream.
func (t *Trivium) nextBits(n uint) uint64 {
	// A máscara de bits é calculada para cobrir n bits
	var bitmask uint64 = (1 << n) - 1

	// Aqui, obtemos os "taps" de nossa matriz de estado
	s66 := (t.state[i66] >> sh66) | (t.state[i66-1] << (wordSize - sh66))
	s93 := (t.state[i93] >> sh93) | (t.state[i93-1] << (wordSize - sh93))
	s162 := (t.state[i162] >> sh162) | (t.state[i162-1] << (wordSize - sh162))
	s177 := (t.state[i177] >> sh177) | (t.state[i177-1] << (wordSize - sh177))
	s243 := (t.state[i243] >> sh243) | (t.state[i243-1] << (wordSize - sh243))
	s288 := (t.state[i288] >> sh288) | (t.state[i288-1] << (wordSize - sh288))

	t1 := s66 ^ s93
	t2 := s162 ^ s177
	t3 := s243 ^ s288

	// armazenada a saida
	z := (t1 ^ t2 ^ t3) & bitmask

	// Agora, processamos os taps
	s91 := (t.state[i91] >> sh91) | (t.state[i91-1] << (wordSize - sh91))
	s92 := (t.state[i92] >> sh92) | (t.state[i92-1] << (wordSize - sh92))
	s171 := (t.state[i171] >> sh171) | (t.state[i171-1] << (wordSize - sh171))
	s175 := (t.state[i175] >> sh175) | (t.state[i175-1] << (wordSize - sh175))
	s176 := (t.state[i176] >> sh176) | (t.state[i176-1] << (wordSize - sh176))
	s264 := (t.state[i264] >> sh264) | (t.state[i264-1] << (wordSize - sh264))
	s286 := (t.state[i286] >> sh286) | (t.state[i286-1] << (wordSize - sh286))
	s287 := (t.state[i287] >> sh287) | (t.state[i287-1] << (wordSize - sh287))
	s69 := (t.state[i69] >> sh69) | (t.state[i69-1] << (wordSize - sh69))

	t1 ^= ((s91 & s92) ^ s171)
	t2 ^= ((s175 & s176) ^ s264)
	t3 ^= ((s286 & s287) ^ s69)
	t1 &= bitmask
	t2 &= bitmask
	t3 &= bitmask

	// Rotação do estado
	t.state[4] = (t.state[4] >> n) | (t.state[3] << (wordSize - n))
	t.state[3] = (t.state[3] >> n) | (t.state[2] << (wordSize - n))
	t.state[2] = (t.state[2] >> n) | (t.state[1] << (wordSize - n))
	t.state[1] = (t.state[1] >> n) | (t.state[0] << (wordSize - n))
	t.state[0] = (t.state[0] >> n) | (t3 << (wordSize - n))

	// Atualiza os valores finais

	n94 := 92 + n
	n178 := 176 + n
	ni94 := n94 >> lgWordSize
	nsh94 := mask - (n94 & mask)
	ni178 := n178 >> lgWordSize
	nsh178 := mask - (n178 & mask)

	t.state[ni94] = t.state[ni94] &^ (bitmask << nsh94)
	t.state[ni94] |= t1 << nsh94

	// Lidando com a sobreposição entre os limites das palavras
	t.state[i94] = t.state[i94] &^ (bitmask >> (wordSize - nsh94))
	t.state[i94] |= t1 >> (wordSize - nsh94)

	t.state[ni178] = t.state[ni178] &^ (bitmask << nsh178)
	t.state[ni178] |= t2 << nsh178
	// Lidando com a sobreposição entre os limites das palavras
	t.state[i178] = t.state[i178] &^ (bitmask >> (wordSize - nsh178))
	t.state[i178] |= t2 >> (wordSize - nsh178)

	return z
}

// nextByte retorna o próximo byte do fluxo de chaves com o MSB como o último bit produzido.
// O primeiro byte produzido terá bits [76543210] da keystream
func (t *Trivium) nextByte() byte {
	return byte(t.nextBits(8))
}

// reverseByte inverte os bits em um byte
func reverseByte(b byte) byte {
	return ((b & 0x1) << 7) | ((b & 0x80) >> 7) |
		((b & 0x2) << 5) | ((b & 0x40) >> 5) |
		((b & 0x4) << 3) | ((b & 0x20) >> 3) |
		((b & 0x8) << 1) | ((b & 0x10) >> 1)
}

// openFile convenience method to open a file or stdin and fatally log on failure
func openFile(filename string) *os.File {
	var file *os.File
	var err error
	if filename == "" {
		file = os.Stdin
	} else {
		file, err = os.Open(filename)
		if err != nil {
			log.Fatalf("error opening %v: %v", filename, err)
		}
	}
	return file
}

// createFile convenience method to create a file or stdout and fatally log on failure
func createFile(filename string) *os.File {
	var file *os.File
	var err error
	if filename == "" {
		file = os.Stdout
	} else {
		file, err = os.Create(filename)
		if err != nil {
			log.Fatalf("error creating %v: %v", filename, err)
		}
	}
	return file
}

// keyGen gera uma chave aleatória de comprimento definido pela constante KeyLength.
// Ele cria um buffer de bytes e usa a função Read da biblioteca "crypto/rand"
// para preenchê-lo com bytes aleatórios.
func (t *Trivium) keyGen() []byte {
	keybuffer := make([]byte, KeyLength)
	rand.Read(keybuffer)
	for i := 0; i < KeyLength; i++ {
		t.Key[i] = keybuffer[i]
	}
	return keybuffer
}

// NewTrivium retorna uma cifra Trivium inicializada com uma chave e um valor de inicialização (IV).
// Tanto a chave quanto o IV são de 80 bits (10 bytes). O processo de inicialização processa a cifra durante
// 4 * 288 ciclos para "aquecer" e tentar eliminar qualquer dependência usável na chave e no IV.
func New(Key ...[KeyLength]byte) *Trivium {
	t := Trivium{}

	// Criar chave
	if len(Key) > 0 {
		t.Key = Key[0]
	} else {
		t.keyGen()
	}

	if len(Key) > 1 {
		t.iv = Key[1]
	} else {
		// Criar IV (vetor inicial)
		ivbuffer := make([]byte, KeyLength)
		rand.Read(ivbuffer)
		for i := 0; i < KeyLength; i++ {
			t.iv[i] = ivbuffer[i]
		}
	}

	// Inicializando o estado com um array de 5 inteiros de 64 bits
	var state [5]uint64

	// Preenchendo o estado com a chave e o IV, cada byte é revertido antes de ser adicionado ao estado
	state[0] |= (uint64(reverseByte(t.Key[0])) << 56) | (uint64(reverseByte(t.Key[1])) << 48) | (uint64(reverseByte(t.Key[2])) << 40) | (uint64(reverseByte(t.Key[3])) << 32)
	state[0] |= (uint64(reverseByte(t.Key[4])) << 24) | (uint64(reverseByte(t.Key[5])) << 16) | (uint64(reverseByte(t.Key[6])) << 8) | uint64(reverseByte(t.Key[7]))
	state[1] |= (uint64(reverseByte(t.Key[8])) << 56) | (uint64(reverseByte(t.Key[9])) << 48)
	state[1] |= (uint64(reverseByte(t.iv[4])) >> 5) | (uint64(reverseByte(t.iv[3])) << 3) | (uint64(reverseByte(t.iv[2])) << 11) | (uint64(reverseByte(t.iv[1])) << 19) | (uint64(reverseByte(t.iv[0])) << 27)
	state[2] |= (uint64(reverseByte(t.iv[7])) << 35) | (uint64(reverseByte(t.iv[6])) << 43) | (uint64(reverseByte(t.iv[5])) << 51) | (uint64(reverseByte(t.iv[4])) << 59)
	state[2] |= (uint64(reverseByte(t.iv[9])) << 19) | (uint64(reverseByte(t.iv[8])) << 27)

	// state[3] é inicializado com todos os zeros
	state[4] |= uint64(7) << 32

	// Inicializando a cifra Trivium com o estado
	t.state = state

	// "Aquecendo" a cifra Trivium processando 4 * 288 bits
	for i := 0; i < 4*288; i++ {
		t.nextBit()
	}

	return &t
}

func (t *Trivium) Encrypt(dst string, src string) {
	// Abrir source
	inputFile := openFile(src)
	defer inputFile.Close()
	reader := bufio.NewReader(inputFile)

	// Cifrando
	var result []byte
	for b, err := reader.ReadByte(); err == nil; b, err = reader.ReadByte() {
		kb := t.nextByte()                // proximo byte da chave de cifra
		result = append(result, (b ^ kb)) // escrever o XOR
	}

	// Salvando em arquivo
	outputFile := createFile(dst)
	defer outputFile.Close()
	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()
	for i := 0; i < KeyLength; i++ {
		writer.WriteByte(t.iv[i])
	}
	writer.Write(result)
}

func (t *Trivium) Decrypt(dst string, src string) {
	// Abrir source
	inputFile := openFile(src)
	defer inputFile.Close()
	reader := bufio.NewReader(inputFile)

	// Recuperar IV
	ivbuffer := make([]byte, KeyLength)
	reader.Read(ivbuffer)
	for i := 0; i < KeyLength; i++ {
		t.iv[i] = ivbuffer[i]
	}

	t = New(t.Key, t.iv)

	// Descifrando
	var result []byte
	for b, err := reader.ReadByte(); err == nil; b, err = reader.ReadByte() {
		kb := t.nextByte()                // next byte of the keystream
		result = append(result, (b ^ kb)) // write the xor out
	}

	// Salvando o arquivo
	outputFile := createFile(dst)
	defer outputFile.Close()
	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()
	writer.Write(result)
}

func (t *Trivium) VirtualDecrypt(src string) string {
	// Abrir source
	inputFile := openFile(src)
	defer inputFile.Close()
	reader := bufio.NewReader(inputFile)

	// Recuperar IV
	ivbuffer := make([]byte, KeyLength)
	reader.Read(ivbuffer)
	for i := 0; i < KeyLength; i++ {
		t.iv[i] = ivbuffer[i]
	}

	t = New(t.Key, t.iv)

	// Descifrando
	var result []byte
	for b, err := reader.ReadByte(); err == nil; b, err = reader.ReadByte() {
		kb := t.nextByte()                // next byte of the keystream
		result = append(result, (b ^ kb)) // write the xor out
	}

	// Salvando o arquivo
	return string(result)
}
