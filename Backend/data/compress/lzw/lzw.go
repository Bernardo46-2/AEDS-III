package lzw

import (
	"fmt"
	"math"
	"os"

	"github.com/Bernardo46-2/AEDS-III/utils"
)

// ================================================ consts ================================================ //

const DICT_MAX_SIZE uint16 = 4096

// ================================================ ZipDict ================================================ //

type ZipDict struct {
    dict map[string]uint16
    size uint16
}

func initDict() ZipDict {
    return ZipDict{
        dict: make(map[string]uint16, DICT_MAX_SIZE),
        size: 0,
    }
}

func (d *ZipDict) isFull() bool {
    return d.size == DICT_MAX_SIZE - 1
}

func (d *ZipDict) push(s string) {
    if !d.isFull() {
        d.size++
        d.dict[s] = d.size
    }
}

func (d *ZipDict) get(s string) (uint16, bool) {
    bytes := []byte(s)

	if len(bytes) == 1 {
		value, _ := utils.BytesToUint16(bytes, 0)
		return value, true
	} else {
		value, contains := d.dict[s]
		return value + math.MaxUint8, contains
	}
}

// ================================================ Bit Compress ================================================ //

func compress12bitArray(bytes []uint16) []byte {
    numBytes := (len(bytes) * 12) / 8
	if (len(bytes)*12)%8 != 0 {
		numBytes++
	}
    bs := make([]byte, numBytes)
    
    offset := 0
    i := 0
    for _, b := range bytes {
        for offset < 12 {
            bits := (b >> (12 - offset - 8)) & 0xFF
            bs[i] |= byte(bits)
            if offset%8 == 4 {
                i++
			}
            offset += 8
        }

        offset -= 12
    }

    return bs
}

// ================================================ LZW ================================================ //

func parseValue(dict *ZipDict, content []byte) (uint16, int) {
    offset := 0
    var value uint16

	for j := 1; j <= len(content); j++ {
		tmp, contains := dict.get(string(content[:j]))

		if contains {
			offset = len(content[:j])
			value = tmp
		} else {
			dict.push(string(content[:j]))
			break
		}
	}

	return value, offset
}

func zip(dict ZipDict, content []byte) []byte {
    zipped := make([]byte, 0, len(content))

	for i := 0; i < len(content); i++ {
		value, offset := parseValue(&dict, content[i:])
		i += offset - 1
		zipped = append(zipped, utils.Uint16ToBytes(value)...)
	}

	return zipped
}

func Zip(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	dict := initDict()
	fmt.Println(content)

    zipped := zip(dict, content)
    fmt.Println(zipped)

    numbers := []uint16{0x123, 0xABC, 0x789}
    fmt.Println(compress12bitArray(numbers))

	return nil
}
