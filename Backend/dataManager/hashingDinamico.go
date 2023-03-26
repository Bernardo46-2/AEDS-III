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
		loadFactor: 5, // int64(float64(numRecords) * 0.05),
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

func (hash *DinamicHash) Add(r BucketRecord) {
	fmt.Printf("\nRegistro a ser adicionado = %+v\n", r)

	// Recuperar e dar parsing no bucket a ser editado
	pos := int64(r.ID) % int64(hash.GetBucketCount())
	fmt.Printf("id = %d | hash = %d | Posicao da hash no bucket = %d ou 0x%X\n", r.ID, pos, hash.directory.bucketPointer[pos], hash.directory.bucketPointer[pos])
	hash.bucketFile.Seek(hash.directory.bucketPointer[pos], io.SeekStart)
	data := make([]byte, hash.bucketSize)
	hash.bucketFile.Read(data)
	bucket := binToBucket(data, hash.loadFactor)

	// Atualiza o bucket com o novo valor
	if bucket.CurrentSize == hash.loadFactor {
		// se o bucket tiver apenas 1 ponteiro aumentar p em +1, se nao so atualiza o bucket
		if bucket.ActualPower == hash.directory.p {
			fmt.Printf("Bucket estourado! Repartindo! aumentando p! p total = %d | p local = %d\n", hash.directory.p, bucket.ActualPower)
			hash.increasePower()
			for i, j := (hash.GetBucketCount() >> 1), 0; i < hash.GetBucketCount(); i, j = i+1, j+1 {
				hash.directory.bucketPointer[i] = hash.directory.bucketPointer[j]
			}

			address := hash.InitializeNewBucket(1)
			hash.directory.bucketPointer[int64(hash.GetBucketCount()>>1)+pos] = address[0]

			// limpeza e reinsercao
			hash.bucketFile.Seek(hash.directory.bucketPointer[pos], io.SeekStart)
			binary.Write(hash.bucketFile, binary.LittleEndian, bucket.ActualPower+1)
			binary.Write(hash.bucketFile, binary.LittleEndian, int64(0))
			binary.Write(hash.bucketFile, binary.LittleEndian, make([]BucketRecord, hash.loadFactor))

			for i := 0; i < len(bucket.Records); i++ {
				hash.Add(bucket.Records[i])
			}
			hash.Add(r)
		} else {
			fmt.Printf("Bucket estourado! Repartindo! p nao aumenta! p total = %d | p local = %d\n", hash.directory.p, bucket.ActualPower)
			address := hash.InitializeNewBucket(1)
			hash.directory.bucketPointer[int64(hash.GetBucketCount()>>1)+pos] = address[0]

			// redivisao do bucket
			hash.bucketFile.Seek(hash.directory.bucketPointer[pos], io.SeekStart)
			binary.Write(hash.bucketFile, binary.LittleEndian, bucket.ActualPower+1)
			binary.Write(hash.bucketFile, binary.LittleEndian, int64(0))
			binary.Write(hash.bucketFile, binary.LittleEndian, make([]BucketRecord, hash.loadFactor))

			for i := 0; i < len(bucket.Records); i++ {
				hash.Add(bucket.Records[i])
			}
			hash.Add(r)
		}
	} else {
		bucket.Records[bucket.CurrentSize] = r
		bucket.CurrentSize++

		fmt.Printf("conteudo = %+v\n", bucket)

		// Grava novamente
		hash.bucketFile.Seek(hash.directory.bucketPointer[pos], io.SeekStart)
		binary.Write(hash.bucketFile, binary.LittleEndian, bucket.ActualPower)
		binary.Write(hash.bucketFile, binary.LittleEndian, bucket.CurrentSize)
		binary.Write(hash.bucketFile, binary.LittleEndian, bucket.Records)
	}
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
}
