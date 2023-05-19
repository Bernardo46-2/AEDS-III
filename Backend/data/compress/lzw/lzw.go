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

func zip(dict Dict, content []byte) []byte {
    zipped := make([]byte, 0, len(content))

    for i := 0; i < len(content); i++ {
        var value uint16
        contains := true
        var j int
        size := 0
        
        for j = i + 1; j <= len(content) && contains; j++ {
            var tmp uint16
            tmp, contains = dict.get(string(content[i:j]))
            
            if contains {
                size = len(content[i:j])
                value = tmp
            } else {
                dict.push(string(content[i:j]))
                contains = false
            }
        }

        i += size-1

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
