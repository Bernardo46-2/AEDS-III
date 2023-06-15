package aescbc

import (
    "errors"
)

type Key struct {
    Key []byte
    Expanded []byte
    SizeBits int
    SizeBytes int
}

func NewKey(sizeBits int) (key Key, err error) {
    if sizeBits != 128 && sizeBits != 192 && sizeBits != 256 {
        err = errors.New("Invalid key size")
        return
    }
    
    sizeBytes := sizeBits / 8
    k, err := randBytes(sizeBytes)
    if err != nil {
        return
    }

    expanded := expand(k)

    key = Key {
        SizeBits: sizeBits,
        SizeBytes: sizeBytes,
        Key: k,
        Expanded: expanded,
    }

    return
}

func NewKeyFrom(k []byte) (key Key, err error) {
    sizeBytes := len(k)
    sizeBits := sizeBytes * 8

    if sizeBits != 128 && sizeBits != 192 && sizeBits != 256 {
        err = errors.New("Invalid key size")
        return
    }

    expanded := expand(k)

    key = Key {
        SizeBits: sizeBits,
        SizeBytes: sizeBytes,
        Key: k,
        Expanded: expanded,
    }

    return
}

func expand(key []byte) (keySchedule []byte) {
    keyLen := len(key)
    var nr, nk int

    switch keyLen {
    case 16:
        nr, nk = 10, 4
    case 24:
        nr, nk = 12, 6
    case 32:
        nr, nk = 14, 8
    }

    keyScheduleLen := 16 * (nr + 1)
    keySchedule = make([]byte, 0, keyScheduleLen)
    keySchedule = append(keySchedule, key...)

    for i := 0; i < 4 * nr + 4 - nk; i++ {
        word := make([]byte, 4)
        copy(word, keySchedule[i * 4 + keyLen - 4 : i * 4 + keyLen])

        if i % nk == 0 {
            RotateLeft(word)
            subBytes(word)
            word[0] ^= RCON[i / nk]
        } else if nk > 6 && i % nk == 4 {
            subBytes(word)
        }

        prevWord := keySchedule[i * 4 : i * 4 + 4]
        xorBlock(word, prevWord)
        keySchedule = append(keySchedule, word...)
    }
    
    return
}