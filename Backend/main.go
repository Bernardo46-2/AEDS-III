// Autores: Marcos Lommez / Bernardo Marques
// Data de criação: 15/03/2023
//
// Programa para gerenciar uma base de dados de Pokemons,
// com suporte a operações crud, diferentes métodos de ordenação externa,
// sistemas de indexação, pattern matching, compressao e criptografia.
// Seu funcionamento é feito atraves de uma comunicação JSON com um frontend
// O servidor HTTP é inicializado na porta 8080.

package main

import (
	"net/http"

	h "github.com/Bernardo46-2/AEDS-III/handlers"
	l "github.com/Bernardo46-2/AEDS-III/logger"
	m "github.com/Bernardo46-2/AEDS-III/middlewares"
)

func Servidor() {
	// Inicializa o servidor de log
	l.LigarServidor()

	// Ordenação externa - TP1
	http.HandleFunc("/ordenacao/", m.EnableCORS(h.Ordenacao))

	// Indexação - TP2
	http.HandleFunc("/getPagesNumber/", m.EnableCORS(h.GetPagesNumber))
	http.HandleFunc("/getIdList", m.EnableCORS(h.GetIdList))
	http.HandleFunc("/getList/", m.EnableCORS(h.GetList))
	http.HandleFunc("/get/", m.EnableCORS(h.GetPokemon))
	http.HandleFunc("/post/", m.EnableCORS(h.PostPokemon))
	http.HandleFunc("/put/", m.EnableCORS(h.PutPokemon))
	http.HandleFunc("/delete/", m.EnableCORS(h.DeletePokemon))
	http.HandleFunc("/loadDatabase", m.EnableCORS(h.LoadDatabase))
	http.HandleFunc("/toKatakana/", m.EnableCORS(h.ToKatakana))

	// Compressao - TP3
	http.HandleFunc("/zip/", m.EnableCORS(h.Zip))
	http.HandleFunc("/unzip/", m.EnableCORS(h.Unzip))

	// Indexacao - TP4
	http.HandleFunc("/mergeSearch/", m.EnableCORS(h.MergeSearch))

	// Criptografia - TP5
	http.HandleFunc("/encrypt/", m.EnableCORS(h.Encrypt))
	http.HandleFunc("/decrypt/", m.EnableCORS(h.Decrypt))

	// Inicializa o servidor HTTP na porta 8080 e escreve no log eventuais erros
	l.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	Servidor()
}
