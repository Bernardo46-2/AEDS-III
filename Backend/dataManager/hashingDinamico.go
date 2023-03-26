package dataManager

import (
	"os"
)

const BUCKETS_FILE string = "data/Buckets.csv"

type DinamicHash struct {
	directory  Directory
	loadFactor int
	bucketFile *os.File
}

type Directory struct {
	p             int
	bucketPointer []int64
}

type Bucket struct {
	CurrentSize int
	Registros   []BucketRecord
}

type BucketRecord struct {
	ID      int32
	Address int64
}

func (d *Directory) GetBucketCount() int {
	return 1 << d.p
}

func newHash(fileAddress string) *DinamicHash {
	// numRegistros, _, _ := NumRegistros()
	arquivo, _ := os.Create(fileAddress)

	d := Directory{
		p:             1,
		bucketPointer: make([]int64, 0, 2),
	}

	hash := DinamicHash{
		directory:  d,
		loadFactor: 5, // int(float64(numRegistros) * 0.05),
		bucketFile: arquivo,
	}

	return &hash
}

/* func (d *Directory) increasePower() {
	d.p++
	novoTamanhoBucket := 1 << d.p
	novoBucket := make([]int, novoTamanhoBucket)
	copy(novoBucket, d.bucketPointer)
	d.bucketPointer = novoBucket
} */

func newBucketRecord(registro Registro) BucketRecord {
	return BucketRecord{
		ID:      registro.Pokemon.Numero,
		Address: registro.Endereco,
	}
}

func (hash *DinamicHash) Add(r BucketRecord) {
	/*
		b := Bucket{
			CurrentSize: 0,
			Registros:   make([]BucketRecord, 0),
		}

		pos := int(r.ID) % d.GetBucketCount()

		b.CurrentSize++
		b.Registros = append(b.Registros, r)
	*/
}

func CriarHashingEstendido() {
	hash := newHash(BUCKETS_FILE)
	c, err := inicializarControleLeitura(BIN_FILE)

	for i := 0; i < 15 && err == nil; i++ { // i < int(c.TotalRegistros)
		err = c.ReadNext()
		if c.RegistroAtual.Lapide != 1 {
			r := newBucketRecord(*c.RegistroAtual)
			hash.Add(r)
		}
	}
}
