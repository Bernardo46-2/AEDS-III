package dataManager

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/Bernardo46-2/AEDS-III/utils"
)

const BUCKETS_FILE string = "data/Buckets.bin"

type DinamicHash struct {
	directory  Directory
	bucketFile *os.File
	loadFactor int64
	bucketSize int64
}

type Directory struct {
	p             int64
	bucketPointer []int64
}

type Bucket struct {
	ActualPower int64
	CurrentSize int64
	Records     []BucketRecord
}

type BucketRecord struct {
	ID      int64
	Address int64
}

func (hash *DinamicHash) GetBucketCount() int {
	return 1 << hash.directory.p
}

func GetBucketPower(power int64) int64 {
	return 1 << power
}

func (hash *DinamicHash) InitializeNewBucket(numberOfBuckets int) []int64 {
	hash.bucketFile.Seek(0, io.SeekEnd)
	bucketAddress := make([]int64, numberOfBuckets)

	for i := 0; i < numberOfBuckets; i++ {
		bucketAddress[i], _ = hash.bucketFile.Seek(0, io.SeekCurrent)
		binary.Write(hash.bucketFile, binary.LittleEndian, hash.directory.p)
		binary.Write(hash.bucketFile, binary.LittleEndian, int64(0))
		binary.Write(hash.bucketFile, binary.LittleEndian, make([]BucketRecord, hash.loadFactor))
	}

	return bucketAddress
}

func newHash(fileAddress string, numRecords int32) *DinamicHash {
	arquivo, _ := os.Create(fileAddress)

	d := Directory{
		p:             1,
		bucketPointer: make([]int64, 2),
	}

	hash := DinamicHash{
		directory:  d,
		loadFactor: 11, // 10 + 1 a mais para tratar overflow
		bucketFile: arquivo,
	}

	ActualPowerSize := int64(binary.Size(Bucket{}.ActualPower))
	currentSizeSize := int64(binary.Size(Bucket{}.CurrentSize))
	bucketRecordSize := int64(binary.Size(BucketRecord{}))
	hash.bucketSize = ActualPowerSize + currentSizeSize + (hash.loadFactor * bucketRecordSize)

	bucketAddress := hash.InitializeNewBucket(hash.GetBucketCount())
	for i := 0; i < hash.GetBucketCount(); i++ {
		hash.directory.bucketPointer[i] = bucketAddress[i]
	}

	return &hash
}

func newBucketRecord(registro Registro) BucketRecord {
	return BucketRecord{
		ID:      int64(registro.Pokemon.Numero),
		Address: registro.Endereco,
	}
}

func (hash *DinamicHash) increasePower() {
	hash.directory.p++
	novoTamanhoBucket := 1 << hash.directory.p
	novoBucket := make([]int64, novoTamanhoBucket)
	copy(novoBucket, hash.directory.bucketPointer)
	hash.directory.bucketPointer = novoBucket

	for i, j := (hash.GetBucketCount() >> 1), 0; i < hash.GetBucketCount(); i, j = i+1, j+1 {
		hash.directory.bucketPointer[i] = hash.directory.bucketPointer[j]
	}
}

func binToBucket(byteArray []byte, numRecords int64) Bucket {
	var ID int64
	var Address int64

	ptr := 0
	bucket := Bucket{}
	bucket.ActualPower, ptr = utils.BytesToInt64(byteArray, ptr)
	bucket.CurrentSize, ptr = utils.BytesToInt64(byteArray, ptr)

	bucket.Records = make([]BucketRecord, numRecords)
	for i := int64(0); i < numRecords; i++ {
		ID, ptr = utils.BytesToInt64(byteArray, ptr)
		Address, ptr = utils.BytesToInt64(byteArray, ptr)
		bucket.Records[i] = BucketRecord{
			ID:      ID,
			Address: Address,
		}
	}

	return bucket
}

func (hash *DinamicHash) insertIntoBucket(pos int64, power int64, currentSize int64, records []BucketRecord) {
	hash.bucketFile.Seek(hash.directory.bucketPointer[pos], io.SeekStart)
	binary.Write(hash.bucketFile, binary.LittleEndian, power)
	binary.Write(hash.bucketFile, binary.LittleEndian, currentSize)
	binary.Write(hash.bucketFile, binary.LittleEndian, records)
}

func (hash *DinamicHash) Add(r BucketRecord) {
	// Recuperar e dar parsing no bucket a ser editado
	pos := int64(r.ID) % int64(hash.GetBucketCount())
	// fmt.Printf("ID   = %3d | Address = %5x\nHash = %3d | Posicao da hash no bucket = %3d ou 0x%4X\n", r.ID, r.Address, pos, hash.directory.bucketPointer[pos], hash.directory.bucketPointer[pos])
	hash.bucketFile.Seek(hash.directory.bucketPointer[pos], io.SeekStart)
	data := make([]byte, hash.bucketSize)
	hash.bucketFile.Read(data)
	bucket := binToBucket(data, hash.loadFactor)

	// Atualiza o bucket com o novo valor
	if bucket.CurrentSize == hash.loadFactor-1 {
		// se o bucket tiver apenas 1 ponteiro aumentar p em +1, se nao so atualiza o bucket
		if bucket.ActualPower == hash.directory.p {
			hash.increasePower()
		}

		address := hash.InitializeNewBucket(1)
		newPos := GetBucketPower(bucket.ActualPower) + pos
		if newPos >= int64(hash.GetBucketCount()) {
			newPos = pos
			pos %= GetBucketPower(bucket.ActualPower)
		}
		hash.directory.bucketPointer[newPos] = address[0]
		bucket.ActualPower++

		// limpeza e reinsercao
		//hash.insertIntoBucket(pos, bucket.ActualPower, 0, make([]BucketRecord, hash.loadFactor))
		bucket1 := Bucket{
			ActualPower: bucket.ActualPower,
			CurrentSize: 0,
			Records:     make([]BucketRecord, hash.loadFactor),
		}
		bucket2 := Bucket{
			ActualPower: bucket.ActualPower,
			CurrentSize: 0,
			Records:     make([]BucketRecord, hash.loadFactor),
		}

		bucket.Records[bucket.CurrentSize] = r
		for i, b1, b2 := 0, 0, 0; i < len(bucket.Records); i++ {
			if bucket.Records[i].ID%GetBucketPower(bucket.ActualPower) == pos {
				bucket1.Records[b1] = bucket.Records[i]
				bucket1.CurrentSize++
				b1++
			} else {
				bucket2.Records[b2] = bucket.Records[i]
				bucket2.CurrentSize++
				b2++
			}
		}

		hash.insertIntoBucket(pos, bucket1.ActualPower, bucket1.CurrentSize, bucket1.Records)
		hash.insertIntoBucket(newPos, bucket2.ActualPower, bucket2.CurrentSize, bucket2.Records)
	} else {
		bucket.Records[bucket.CurrentSize] = r
		hash.insertIntoBucket(pos, bucket.ActualPower, bucket.CurrentSize+1, bucket.Records)
	}
}

func (hash *DinamicHash) PrintHash() {
	seen := make(map[int64]bool)
	fmt.Println()
	for i := 0; i < len(hash.directory.bucketPointer); i++ {
		if !seen[hash.directory.bucketPointer[i]] {
			hash.bucketFile.Seek(hash.directory.bucketPointer[i], io.SeekStart)
			data := make([]byte, hash.bucketSize)
			hash.bucketFile.Read(data)
			bucket := binToBucket(data, hash.loadFactor)
			fmt.Printf("[%3d] %5X = ActualPower:%d | CurrentSize:%d | [ ", i, hash.directory.bucketPointer[i], bucket.ActualPower, bucket.CurrentSize)
			for i := 0; i < len(bucket.Records); i++ {
				if bucket.Records[i].ID != 0 {
					fmt.Printf("{ID:%3d Address:%5X} ", bucket.Records[i].ID, bucket.Records[i].Address)
				} else {
					fmt.Printf("                       ")
				}
			}
			fmt.Printf("]\n")
			seen[hash.directory.bucketPointer[i]] = true
		} else {
			fmt.Printf("[%3d] %5X\n", i, hash.directory.bucketPointer[i])
		}
	}
	fmt.Printf("\n")
}

func StartHashFile() {
	c, err := inicializarControleLeitura(BIN_FILE)
	hash := newHash(BUCKETS_FILE, c.TotalRegistros)

	for i := 0; i < int(c.TotalRegistros) && err == nil; i++ {
		err = c.ReadNext()
		if c.RegistroAtual.Lapide != 1 {
			r := newBucketRecord(*c.RegistroAtual)
			hash.Add(r)
		}
	}

	hash.PrintHash()
}
