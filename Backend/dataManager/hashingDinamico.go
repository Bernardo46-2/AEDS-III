package dataManager

import "fmt"

type diretorio struct {
	p             int
	bucketPointer []int
	loadFactor    int
}

type Bucket struct {
	TamanhoAtual int
	Registros    []RegistroBucket
}

type RegistroBucket struct {
	ID       int32
	Endereco int64
}

func (d *diretorio) numeroDeBuckets() int {
	return 1 << d.p
}

func criarDiretorio() *diretorio {
	numRegistros, _, _ := NumRegistros()
	return &diretorio{
		p:             1,
		bucketPointer: make([]int, 0, 2),
		loadFactor:    int(float64(numRegistros) * 0.05),
	}
}

/* func (d *diretorio) aumentarP() {
	d.p++
	novoTamanhoBucket := 1 << d.p
	novoBucket := make([]int, novoTamanhoBucket)
	copy(novoBucket, d.bucketPointer)
	d.bucketPointer = novoBucket
} */

func CriarHashingEstendido() {
	d := criarDiretorio()
	fmt.Printf("p = %d, numero de buckets = %d, fator de carga = %d\n\n", d.p, d.numeroDeBuckets(), d.loadFactor)

	c, err := inicializarControleLeitura(BIN_FILE)
	b := Bucket{
		TamanhoAtual: 0,
		Registros:    make([]RegistroBucket, 0),
	}

	for i := 0; i < d.loadFactor && err == nil; i++ { // i < int(c.TotalRegistros)
		err = c.ReadNext()
		if c.RegistroAtual.Lapide != 1 {
			r := RegistroBucket{
				ID:       c.RegistroAtual.Pokemon.Numero,
				Endereco: c.RegistroAtual.Endereco,
			}
			b.TamanhoAtual++
			b.Registros = append(b.Registros, r)
		}
	}

	for i := 0; i < b.TamanhoAtual; i++ {
		fmt.Printf("ID = %d, Posicao = %x\n", b.Registros[i].ID, b.Registros[i].Endereco)
	}
}
