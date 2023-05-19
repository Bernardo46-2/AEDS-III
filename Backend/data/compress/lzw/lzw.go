package lzw

import (
	// "os"
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
    value, contains := d.dict[s]
    if contains {
        value += math.MaxUint8
    } else {
        value, _ = utils.BytesToUint16([]byte(s), 0)
    }
    return value, contains
}

// ================================================ LZW ================================================ //

func zip(dict Dict, content []byte) []byte {
    zipped := make([]byte, 0, len(content))

    for i := 0; i < len(content); i++ {
        j := i+1
        value, contains := dict.get(string(content[i:j]))
        fmt.Printf("content[%d:%d]: %+v\n", i, j, content[i:j])
        fmt.Println("value:", value)
        
        for j = i + 2; j <= len(content) && contains; j++ {
            fmt.Printf("searching content[%d:%d]: %+v\n", i, j, content[i:j])
            tmp, contains := dict.get(string(content[i:j]))
            if contains { value = tmp }
        }
        
        if j <= len(content) {
            dict.push(string(content[i:j]))
        }

        if contains {
            zipped = append(zipped, utils.Uint16ToBytes(value + math.MaxUint8)...)
        } else if j < len(content) {
            zipped = append(zipped, content[i])
        }

        fmt.Println("")
    }

    return zipped
}

func Zip(path string) error {
    // content, err := os.ReadFile(path)
    // if err != nil {
    //     return err
    // }

    content := []byte{101, 102, 101, 102}
    dict := initDict()
    fmt.Println(content)

    zipped := zip(dict, content)
    fmt.Println(zipped)

    return nil
}
