package dataManager

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/Bernardo46-2/AEDS-III/models"
	"github.com/Bernardo46-2/AEDS-III/utils"
)

const INTERCALACAO_VARIAVEL_MIN_BLOCK_SIZE int64 = 8192

func IntercalacaoTamanhoVariavel() {
	// blockSize := INTERCALACAO_VARIAVEL_MIN_BLOCK_SIZE
	sortBlocksToFile(BIN_FILE, INTERCALACAO_VARIAVEL_MIN_BLOCK_SIZE, "data/tmp/variadic_sort")

	/* 	for sorted := false; !sorted; {
		blockSize *= 2
		sorted = sortFiles("data/tmp/variadic_sort1.dat", "data/tmp/variadic_sort2.dat", blockSize, "data/tmp/variadic_sort3.dat", "data/tmp/variadic_sort4.dat")

		if !sorted {
			blockSize *= 2
			sorted = sortFiles("data/tmp/variadic_sort3.dat", "data/tmp/variadic_sort4.dat", blockSize, "data/tmp/variadic_sort1.dat", "data/tmp/variadic_sort2.dat")
		}
	} */
}

func extractBlockFixedSize(f *os.File, registerStart int64, registersRead *int, numRegisters int, blockSize int64, count *int) ([]models.Pokemon, int64) {
	ptr := registerStart
	currentBlockSize := int64(0)
	block := []models.Pokemon{}
	full := false

	for !full && *registersRead < numRegisters {
		ptr, _ = f.Seek(0, io.SeekCurrent)
		registerSize, valid, _ := tamanhoProxRegistro(f, ptr)

		if valid != 0 {
			if registerSize+currentBlockSize > blockSize {
				f.Seek(-8, io.SeekCurrent)
				full = true
			} else {
				currentBlockSize += registerSize
				pokemon, _, _ := readRegistro(f, ptr)
				pokemon.CalculateSize()
				*registersRead += 1
				*count++
				block = append(block, pokemon)
			}
		} else {
			readRegistro(f, ptr)
		}
	}

	return block, currentBlockSize
}

func extractBlockAnySize(f *os.File, registerStart int64, registersRead *int, numRegisters int, maxBlockSize int64) ([]models.Pokemon, int64) {
	ptr := registerStart
	currentBlockSize := int64(0)
	block := []models.Pokemon{}
	full := false
	lastPokemon := models.Pokemon{Numero: -1}
	for !full && *registersRead < numRegisters {
		ptr, _ = f.Seek(0, io.SeekCurrent)
		registerSize, valid, _ := tamanhoProxRegistro(f, ptr)
		fmt.Println(ptr)
		pokemon, positionBackup, _ := readRegistro(f, ptr)

		if valid != 0 {
			if registerSize+currentBlockSize > maxBlockSize || pokemon.Numero > lastPokemon.Numero {
				f.Seek(positionBackup-8, io.SeekStart)
				full = true
			} else {
				currentBlockSize += registerSize
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
	tmp, err2 := os.Create(outputFile + "1.dat")
	outputFiles[0] = tmp
	tmp, err3 := os.Create(outputFile + "2.dat")
	outputFiles[1] = tmp
	outFileSizes := make([]int64, 2)
	whichFile := 0
	numRegistros := make([]int, 2)

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
		block, currentBlockSize := extractBlockFixedSize(inFile, registerStart, &i, numRegisters, blockSize, &numRegistros[whichFile])
		outFileSizes[whichFile] += currentBlockSize

		sort.Slice(block, func(i, j int) bool {
			return block[i].Numero < block[j].Numero
		})

		for i := 0; i < len(block); i++ {
			tmp := block[i].ToBytes()
			binary.Write(outputFiles[whichFile], binary.LittleEndian, tmp)
		}

		if whichFile == 0 {
			whichFile = 1
		} else {
			whichFile = 0
		}
	}

	outputFiles[0].Seek(0, io.SeekStart)
	binary.Write(outputFiles[0], binary.LittleEndian, utils.IntToBytes(int32(numRegistros[0])))
	outputFiles[1].Seek(0, io.SeekStart)
	binary.Write(outputFiles[1], binary.LittleEndian, utils.IntToBytes(int32(numRegistros[1])))
}

func mergeBlocksToFile(file *os.File, numRegistros int, block1, block2 []models.Pokemon) {
	var n int

	if len(block1) < len(block2) {
		n = len(block1)
	} else {
		n = len(block2)
	}

	for i := 0; i < n; i++ {
		if i < len(block1) && i < len(block2) {
			if block1[i].Numero < block2[i].Numero {
				binary.Write(file, binary.LittleEndian, block1[i].ToBytes())
				binary.Write(file, binary.LittleEndian, block2[i].ToBytes())
			} else {
				binary.Write(file, binary.LittleEndian, block2[i].ToBytes())
				binary.Write(file, binary.LittleEndian, block1[i].ToBytes())
			}
		} else if i < len(block1) {
			binary.Write(file, binary.LittleEndian, block1[i].ToBytes())
		} else if i < len(block2) {
			binary.Write(file, binary.LittleEndian, block2[i].ToBytes())
		}
	}

	file.Seek(0, io.SeekStart)
	binary.Write(file, binary.LittleEndian, int32(n))
}

func sortFiles(inputFile1, inputFile2 string, maxBlockSize int64, outputFile1, outputFile2 string) bool {
	file1, err1 := os.Open(inputFile1)
	file2, err2 := os.Open(outputFile2)
	file3, err3 := os.Create(outputFile1)
	file4, err4 := os.Create(outputFile2)
	file3Size := int64(0)
	file4Size := int64(0)
	whichFile := 0
	sorted := false

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		panic("Error opening files")
	}
	defer file1.Close()
	defer file2.Close()
	defer file3.Close()
	defer file4.Close()

	numRegisters1, start1, _ := NumRegistros()
	numRegisters2, start2, _ := NumRegistros()
	numRegisters := numRegisters1 + numRegisters2

	start1, _ = file1.Seek(start1, io.SeekStart)
	start2, _ = file2.Seek(start2, io.SeekStart)
	binary.Write(file3, binary.LittleEndian, utils.IntToBytes(int32(file3Size)))
	binary.Write(file4, binary.LittleEndian, utils.IntToBytes(int32(file4Size)))

	for i := 0; i < numRegisters; i++ {
		registerStart1, _ := file1.Seek(0, io.SeekCurrent)
		registerStart2, _ := file2.Seek(0, io.SeekCurrent)
		block1, currentBlockSize1 := extractBlockAnySize(file1, registerStart1, &i, numRegisters1, maxBlockSize)
		file3Size += currentBlockSize1
		block2, currentBlockSize2 := extractBlockAnySize(file2, registerStart2, &i, numRegisters2, maxBlockSize)
		file4Size += currentBlockSize2

		binary.Write(file3, binary.LittleEndian, utils.IntToBytes(int32(len(block1))))
		binary.Write(file4, binary.LittleEndian, utils.IntToBytes(int32(len(block2))))

		if whichFile == 0 {
			mergeBlocksToFile(file3, numRegisters, block1, block2)
			whichFile = 1
		} else {
			mergeBlocksToFile(file4, numRegisters, block1, block2)
			whichFile = 0
		}

		sorted = len(block1) == 1 && len(block2) == 1
		if whichFile == 0 {
			whichFile = 1
		} else {
			whichFile = 0
		}
	}

	return sorted
}
