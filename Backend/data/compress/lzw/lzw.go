package lzw

import (
	"os"
	"fmt"
    "math"

	"github.com/Bernardo46-2/AEDS-III/utils"
)

// ================================================ consts ================================================ //

const DICT_MAX_SIZE uint16 = 4096

// ================================================ Dict ================================================ //

type Dict struct {
    dict map[string]uint16
    size uint16
}

func initDict() Dict {
    return Dict{
        dict: make(map[string]uint16, DICT_MAX_SIZE),
        size: 0,
    }
}

func (d *Dict) isFull() bool {
    return d.size == DICT_MAX_SIZE - 1
}

func (d *Dict) push(s string) {
    if !d.isFull() {
        d.size++
        d.dict[s] = d.size
    }
}

func (d *Dict) get(s string) (uint16, bool) {
    bytes := []byte(s)

    if len(bytes) == 1 {
        value, _ := utils.BytesToUint16(bytes, 0)
        return value, true
    } else {
        value, contains := d.dict[s]
        return value + math.MaxUint8, contains
    }
}

// ================================================ LZW ================================================ //

func parseValue(dict *Dict, content []byte) (uint16, int) {
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

func zip(dict Dict, content []byte) []byte {
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

    return nil
}
