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
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Bernardo46-2/AEDS-III/handlers"
	"github.com/Bernardo46-2/AEDS-III/logger"
	"github.com/Bernardo46-2/AEDS-III/middlewares"
	"github.com/Bernardo46-2/AEDS-III/service"
)

func Servidor() {
	// Inicializa o servidor de log
	logger.LigarServidor()

	// Crud
	http.HandleFunc("/getPagesNumber/", middlewares.EnableCORS(handlers.GetPagesNumber))
	http.HandleFunc("/getIdList", middlewares.EnableCORS(handlers.GetIdList))
	http.HandleFunc("/getList/", middlewares.EnableCORS(handlers.GetList))
	http.HandleFunc("/get/", middlewares.EnableCORS(handlers.GetPokemon))
	http.HandleFunc("/post/", middlewares.EnableCORS(handlers.PostPokemon))
	http.HandleFunc("/put/", middlewares.EnableCORS(handlers.PutPokemon))
	http.HandleFunc("/delete/", middlewares.EnableCORS(handlers.DeletePokemon))
	http.HandleFunc("/loadDatabase", middlewares.EnableCORS(handlers.LoadDatabase))
	http.HandleFunc("/toKatakana/", middlewares.EnableCORS(handlers.ToKatakana))
	http.HandleFunc("/mergeSearch/", middlewares.EnableCORS(handlers.MergeSearch))

	// Ordenação externa
	http.HandleFunc("/ordenacao/", middlewares.EnableCORS(handlers.Ordenacao))

	// Criptografia
	http.HandleFunc("/Encrypt/", middlewares.EnableCORS(handlers.Encrypt))

	// Inicializa o servidor HTTP na porta 8080 e escreve no log eventuais erros
	logger.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	switch os.Args[1] {
	case "0":
		key := service.Encrypt(1)
		fmt.Printf("key = %s", key)
	case "1":
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		key, _ := reader.ReadString('\n')
		key = strings.TrimSpace(key)
		service.Decrypt(1, key)
	}
	// servidor()
}
