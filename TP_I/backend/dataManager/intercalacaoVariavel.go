package dataManager

import (
	"encoding/binary"
	// "fmt"
	"io"
	"os"
	"sort"

	"github.com/Bernardo46-2/AEDS-III/models"
	"github.com/Bernardo46-2/AEDS-III/utils"
)

func IntercalacaoTamanhoVariavel() {
    sortBlocksToFile(BIN_FILE, 8192, "data/tmp/first_sort.dat")
}

func extractBlockFixedSize(f *os.File, registerStart int64, registersRead *int, numRegisters int, blockSize int64) ([]models.Pokemon, int64) {
    ptr := registerStart
    currentBlockSize := int64(0)
    block := []models.Pokemon{}
    full := false

    for !full && *registersRead < numRegisters {
        ptr, _ = f.Seek(0, io.SeekCurrent)
        registerSize, dead, _ := tamanhoProxRegistro(f, ptr)
        
        if dead != 0 {
            if registerSize + currentBlockSize > blockSize {
                f.Seek(-8, io.SeekCurrent)
                full = true
            } else {
                currentBlockSize += registerSize
                pokemon, _, _ := readRegistro(f, ptr)
                pokemon.CalculateSize()
                *registersRead += 1
                block = append(block, pokemon)
            }
        } else {
            readRegistro(f, ptr)
        }
    }
    
    return block, currentBlockSize
}

func extractBlockAnySize(f *os.File, registerStart int64, registersRead *int, numRegisters int) ([]models.Pokemon, int64) {
    ptr := registerStart
    currentBlockSize := int64(0)
    block := []models.Pokemon{}
    full := false
    // nextPokemonvalue := 0

    for !full && *registersRead < numRegisters {
        ptr, _ = f.Seek(0, io.SeekCurrent)
        registerSize, dead, _ := tamanhoProxRegistro(f, ptr)
        
        if dead != 0 {
            if registerSize + currentBlockSize > blockSize {
                f.Seek(-8, io.SeekCurrent)
                full = true
            } else {
                currentBlockSize += registerSize
                pokemon, _, _ := readRegistro(f, ptr)
                pokemon.CalculateSize()
                *registersRead += 1
                block = append(block, pokemon)
            }
        } else {
            readRegistro(f, ptr)
        }
    }
    
    return block, currentBlockSize
}

func sortBlocksToFile(inputFile string, blockSize int64, outputFile string) {
    inFile, err1 := os.Open(inputFile)
    outputFiles := make([]*os.File, 2)
    tmp, err2 := os.Create(outputFile + "1")
    outputFiles[0] = tmp 
    tmp, err3 := os.Create(outputFile + "1")
    outputFiles[1] = tmp
    outFileSizes := make([]int64, 2)
    whichFile := 0

    if err1 != nil || err2 != nil || err3 != nil {
        return
    }
    defer inFile.Close()
    defer outputFiles[0].Close()
    defer outputFiles[1].Close()

    numRegisters, start, _ := NumRegistros()
    start, _ = inFile.Seek(start, io.SeekStart)
    binary.Write(outputFiles[0], binary.LittleEndian, utils.IntToBytes(int32(outFileSizes[0])))
    binary.Write(outputFiles[1], binary.LittleEndian, utils.IntToBytes(int32(outFileSizes[1])))

    for i := 0; i < numRegisters; i++ {
        registerStart, _ := inFile.Seek(0, io.SeekCurrent)
        block, currentBlockSize := extractBlockFixedSize(inFile, registerStart, &i, numRegisters, blockSize)
        outFileSizes[whichFile] += currentBlockSize
        
        sort.Slice(block, func(i , j int) bool {
            return block[i].Numero < block[j].Numero
        })
        
        binary.Write(outputFiles[whichFile], binary.LittleEndian, utils.IntToBytes(int32(len(block))))

		for i := 0; i < len(block); i++ {
			tmp := block[i].ToBytes()
			binary.Write(outputFiles[whichFile], binary.LittleEndian, tmp)
		}

        if whichFile == 0 { whichFile = 1 } else { whichFile = 0 }
    }

    outputFiles[0].Seek(0, io.SeekStart)
    binary.Write(outputFiles[0], binary.LittleEndian, utils.IntToBytes(int32(outFileSizes[0])))
    outputFiles[1].Seek(0, io.SeekStart)
    binary.Write(outputFiles[1], binary.LittleEndian, utils.IntToBytes(int32(outFileSizes[1])))
}
