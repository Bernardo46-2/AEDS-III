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

	"github.com/Bernardo46-2/AEDS-III/data/binManager"
	"github.com/Bernardo46-2/AEDS-III/data/indexes/bplustree"
	"github.com/Bernardo46-2/AEDS-III/data/indexes/btree"
	"github.com/Bernardo46-2/AEDS-III/data/indexes/hashing"
	"github.com/Bernardo46-2/AEDS-III/handlers"
	"github.com/Bernardo46-2/AEDS-III/logger"
	"github.com/Bernardo46-2/AEDS-III/middlewares"
	"github.com/Bernardo46-2/AEDS-III/service"
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
		fmt.Println("5 | criarIndiceInvertido")
		fmt.Println("6 | lerMultiplosIndiceInvertido")
		fmt.Println("7 | b+ tree")
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
		binManager.ImportCSV().CsvToBin()
	case "3", "hash":
		controler, _ := binManager.InicializarControleLeitura(binManager.BIN_FILE)
		defer controler.Close()
		hashing.StartHashFile(controler, 8, binManager.FILES_PATH, "hashIndex")
	case "4", "btree":
		btree.StartBTreeFile(binManager.FILES_PATH)
	case "7":
		controler, _ := binManager.InicializarControleLeitura(binManager.BIN_FILE)
		bplustree.StartBPlusTreeFile(binManager.FILES_PATH, "altura", controler)
	case "8":
		tree, _ := bplustree.ReadBPlusTree(binManager.FILES_PATH, "numero")
		tree.FindRange(256, 300)
		tree.PrintFile()
	case "9":
		var req service.SearchRequest
		req.Nome = ""
		req.JapName = ""
		req.Especie = ""
		req.Tipo = ""
		req.Descricao = ""
		req.IDI = ""
		req.IDF = ""
		req.GeracaoI = ""
		req.GeracaoF = ""
		req.LancamentoI = ""
		req.LancamentoF = ""
		req.AtkI = ""
		req.AtkF = "1000"
		req.DefI = ""
		req.DefF = ""
		req.HpI = ""
		req.HpF = ""
		req.AlturaI = ""
		req.AlturaF = ""
		req.PesoI = ""
		req.PesoF = ""

		service.MergeSearch(req)
	default:
		fmt.Println("Opção inválida")
	}
}

func servidor() {
	// Inicializa o servidor de log
	logger.LigarServidor()

	// Teste
	http.HandleFunc("/ping/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "pong")
	})

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
	http.HandleFunc("/mergeIndex/", middlewares.EnableCORS(handlers.MergeSearch))

	// Ordenação externa
	http.HandleFunc("/ordenacao/", middlewares.EnableCORS(handlers.Ordenacao))

	// Inicializa o servidor HTTP na porta 8080 e escreve no log eventuais erros
	logger.Fatal(http.ListenAndServe(":8080", nil))
}
