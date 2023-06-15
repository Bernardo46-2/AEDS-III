package aescbc

import "os"

const BLOCK_SIZE int = 16

func gmul(a, b byte) (p byte) {
    for a != 0 && b != 0 {
        if b & 1 != 0 {
            p ^= a
        }

        c := a & 0x80 != 0
        a <<= 1
        if c {
            a ^= 0x1b
        }
        b >>= 1
    }
    
    return   
}

func addPadding(blocks [][]byte) ([][]byte) {
    last := len(blocks) - 1
    lastLen := len(blocks[last])
    b := byte(BLOCK_SIZE - lastLen)

    if b == 0 {
        b = byte(BLOCK_SIZE)
        lastLen = 0
        last++
        tmp := make([]byte, 0, BLOCK_SIZE)
        blocks = append(blocks, tmp)
    }

    for i := lastLen; i < BLOCK_SIZE; i++ {
        blocks[last] = append(blocks[last], b)
    }

    return blocks
}

func removePadding(state []byte) []byte {
    paddingLen := int(state[len(state)-1])
    return state[:len(state) - paddingLen]
}

func splitBlocks(bytes []byte) (blocks [][]byte) {
    blocksLen := len(bytes) / BLOCK_SIZE
    if len(bytes) % BLOCK_SIZE != 0 {
        blocks = make([][]byte, 0, blocksLen+1)
    } else {
        blocks = make([][]byte, 0, blocksLen)
    }
    
    for i := 0; i < blocksLen; i++ {
        tmp := make([]byte, BLOCK_SIZE)
        copy(tmp, bytes[i * BLOCK_SIZE : i * BLOCK_SIZE + BLOCK_SIZE])
        blocks = append(blocks, tmp)
    }
    
    if len(bytes) % BLOCK_SIZE != 0 {
        tmp := make([]byte, 0, BLOCK_SIZE)
        tmp = append(tmp, bytes[blocksLen * BLOCK_SIZE:]...)
        blocks = append(blocks, tmp)
    }

    return
}

func rotateBytesLeft(bytes []byte, start, end int) {
    tmp := bytes[start]
    for i := start; i < end - 1; i++ {
        bytes[i] = bytes[i + 1]
    }
    bytes[end - 1] = tmp
}

func rotateBytesRight(bytes []byte, start, end int) {
    tmp := bytes[end - 1]
    for i := end - 1; i > start; i-- {
        bytes[i] = bytes[i - 1]
    }
    bytes[start] = tmp
}

func subBytes(array []byte) {
    for i := 0; i < len(array); i++ {
        array[i] = S_BOX[array[i]]
    }
}

func invSubBytes(array []byte) {
    for i := 0; i < len(array); i++ {
        array[i] = INV_S_BOX[array[i]]
    }
}

func shiftRows(block []byte) {
    rotateBytesLeft(block, 4, 8)
    rotateBytesLeft(block, 8, 12)
    rotateBytesLeft(block, 8, 12)
    rotateBytesLeft(block, 12, 16)
    rotateBytesLeft(block, 12, 16)
    rotateBytesLeft(block, 12, 16)
}

func invShiftRows(block []byte) {
    rotateBytesRight(block, 4, 8)
    rotateBytesRight(block, 8, 12)
    rotateBytesRight(block, 8, 12)
    rotateBytesRight(block, 12, 16)
    rotateBytesRight(block, 12, 16)
    rotateBytesRight(block, 12, 16)
}

func mixColumns(bytes []byte) {
    for i := 0; i < 4; i++ {
        i0, i1, i2, i3 := i, i + 4, i + 8, i + 12
        s0, s1, s2, s3 := bytes[i0], bytes[i1], bytes[i2], bytes[i3]

        bytes[i0] = gmul(2, s0) ^ gmul(3, s1) ^ s2 ^ s3
        bytes[i1] = s0 ^ gmul(2, s1) ^ gmul(3, s2) ^ s3
        bytes[i2] = s0 ^ s1 ^ gmul(2, s2) ^ gmul(3, s3)
        bytes[i3] = gmul(3, s0) ^ s1 ^ s2 ^ gmul(2, s3)
    }    
}

func invMixColumns(bytes []byte) {
    for i := 0; i < 4; i++ {
        i0, i1, i2, i3 := i, i + 4, i + 8, i + 12
        s0, s1, s2, s3 := bytes[i0], bytes[i1], bytes[i2], bytes[i3]

        bytes[i0] = gmul(0x0e, s0) ^ gmul(0x0b, s1) ^ gmul(0x0d, s2) ^ gmul(0x09, s3)
        bytes[i1] = gmul(0x09, s0) ^ gmul(0x0e, s1) ^ gmul(0x0b, s2) ^ gmul(0x0d, s3)
        bytes[i2] = gmul(0x0d, s0) ^ gmul(0x09, s1) ^ gmul(0x0e, s2) ^ gmul(0x0b, s3)
        bytes[i3] = gmul(0x0b, s0) ^ gmul(0x0d, s1) ^ gmul(0x09, s2) ^ gmul(0x0e, s3)
    }
}

func addRoundKey(key, bytes []byte) {
    xorBlock(bytes, key)
}

func Encrypt(k Key, iv, data []byte) (cipherText []byte) {
    nr := len(k.Expanded) / BLOCK_SIZE
    state := splitBlocks(data)
    state = addPadding(state)
    numBlocks := len(state)
    cipherText = make([]byte, 0, numBlocks * (BLOCK_SIZE + 1))
    previousBlock := iv
    cipherText = append(cipherText, previousBlock...)
    
    for i := 0; i < numBlocks; i++ {
        xorBlock(state[i], previousBlock)
        addRoundKey(k.Expanded[:BLOCK_SIZE], state[i])

        for j := 1; j < nr - 1; j++ {
            subBytes(state[i])
            transposeMatrix(state[i], 4, 4)
            shiftRows(state[i])
            mixColumns(state[i])
            transposeMatrix(state[i], 4, 4)
            addRoundKey(k.Expanded[j * BLOCK_SIZE : (j + 1) * BLOCK_SIZE], state[i])
        }
        
        subBytes(state[i])
        transposeMatrix(state[i], 4, 4)
        shiftRows(state[i])
        transposeMatrix(state[i], 4, 4)
        addRoundKey(k.Expanded[(nr - 1) * BLOCK_SIZE:], state[i])

        previousBlock = state[i]
        cipherText = append(cipherText, state[i]...)
    }

    return
}

func Decrypt(k Key, data []byte) (decryptedText []byte) {
    keyLen := len(k.Expanded)
    numRounds := keyLen / BLOCK_SIZE
    state := splitBlocks(data)
    numBlocks := len(state)
    decryptedText = make([]byte, 0, (numBlocks - 1) * BLOCK_SIZE)

    for i := numBlocks - 1; i > 0; i-- {
        addRoundKey(k.Expanded[(numRounds - 1) * BLOCK_SIZE:], state[i])
        transposeMatrix(state[i], 4, 4)
        invShiftRows(state[i])
        transposeMatrix(state[i], 4, 4)
        invSubBytes(state[i])
        
        for j := 1; j < numRounds - 1; j++ {
            addRoundKey(k.Expanded[keyLen - (j + 1) * BLOCK_SIZE : keyLen - j * BLOCK_SIZE], state[i])
            transposeMatrix(state[i], 4, 4)
            invMixColumns(state[i])
            invShiftRows(state[i])
            transposeMatrix(state[i], 4, 4)
            invSubBytes(state[i])
        }

        previousBlock := state[i - 1]
        addRoundKey(k.Expanded[:BLOCK_SIZE], state[i])
        xorBlock(state[i], previousBlock)
        decryptedText = prependSlice(decryptedText, state[i])
    }
    
    decryptedText = removePadding(decryptedText)
    return
}

func main() {
    key, _ := NewKey(128)
    iv, _ := randBytes(BLOCK_SIZE)
    file, _ := os.ReadFile("pokedex.csv")
    
    encrypted := Encrypt(key, iv, file)
    os.WriteFile("encrypted.txt", encrypted, 0644)
    
    decrypted := Decrypt(key, encrypted)
    os.WriteFile("decrypted.txt", decrypted, 0644)
}