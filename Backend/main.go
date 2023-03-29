// Autores: Marcos Lommez / Bernardo Marques
// Data de criação: 15/03/2023
//
// Programa para gerenciar uma base de dados de Pokemons,
// com suporte a operações crud e diferentes métodos de ordenação externa.
// Seu funcionamento é feito atraves de uma comunicação JSON com um frontend
// O servidor HTTP é inicializado na porta 8080.

package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/Bernardo46-2/AEDS-III/dataManager"
	"github.com/Bernardo46-2/AEDS-III/handlers"
	"github.com/Bernardo46-2/AEDS-III/logger"
	"github.com/Bernardo46-2/AEDS-III/middlewares"
)

func main() {
	var opcao string
	if len(os.Args) > 1 {
		opcao = os.Args[1]
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("1 | server")
		fmt.Println("2 | csv")
		fmt.Println("3 | hash")
		fmt.Println("4 | btree")
		fmt.Println("5 | searchHash")
		fmt.Print("\n> ")
		scanner.Scan()
		opcao = scanner.Text()
		fmt.Println()
	}

	switch opcao {
	case "1", "server":
		fmt.Println("Servidor Iniciado")
		servidor()
	case "2", "csv":
		dataManager.ImportCSV().CsvToBin()
	case "3", "hash":
		dataManager.StartHashFile()
	case "4", "btree":
		dataManager.StartBTreeFile()
	case "5", "searchHash":
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("\nID: ")
		scanner.Scan()
		str := scanner.Text()
		id, _ := strconv.Atoi(str)
		fmt.Printf("%+v", dataManager.RecoverRegisterAddress(int64(id)))
	default:
		fmt.Println("Opção inválida")
	}
}

func servidor() {
	// Inicializa o servidor de log
	logger.LigarServidor()

	// Crud
	http.HandleFunc("/getPagesNumber/", middlewares.EnableCORS(handlers.GetPagesNumber))
	http.HandleFunc("/getAll/", middlewares.EnableCORS(handlers.GetAllPokemon))
	http.HandleFunc("/get/", middlewares.EnableCORS(handlers.GetPokemon))
	http.HandleFunc("/post/", middlewares.EnableCORS(handlers.PostPokemon))
	http.HandleFunc("/put/", middlewares.EnableCORS(handlers.PutPokemon))
	http.HandleFunc("/delete/", middlewares.EnableCORS(handlers.DeletePokemon))
	http.HandleFunc("/loadDatabase", middlewares.EnableCORS(handlers.LoadDatabase))
	http.HandleFunc("/toKatakana/", middlewares.EnableCORS(handlers.ToKatakana))

	// Ordenação externa
	http.HandleFunc("/intercalacaoComum/", middlewares.EnableCORS(handlers.IntercalacaoComum))
	http.HandleFunc("/intercalacaoVariavel/", middlewares.EnableCORS(handlers.IntercalacaoVariavel))
	http.HandleFunc("/selecaoPorSubstituicao/", middlewares.EnableCORS(handlers.SelecaoPorSubstituicao))

	// Indexação
	http.HandleFunc("/criarHashingEstendido/", middlewares.EnableCORS(handlers.CriarHashingEstendido))

	// Inicializa o servidor HTTP na porta 8080 e escreve no log eventuais erros
	logger.Fatal(http.ListenAndServe(":8080", nil))
}
