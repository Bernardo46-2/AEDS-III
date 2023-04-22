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
	"strings"

	"github.com/Bernardo46-2/AEDS-III/data/binManager"
	"github.com/Bernardo46-2/AEDS-III/data/indexes/bplustree"
	"github.com/Bernardo46-2/AEDS-III/data/indexes/btree"
	"github.com/Bernardo46-2/AEDS-III/data/indexes/hashing"
	"github.com/Bernardo46-2/AEDS-III/data/indexes/invertedIndex"
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
	case "5", "criarIndiceInvertido":

	case "6", "lerMultiplosIndiceInvertido":
		var strs []string
		var line string
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Printf("nome\n>")
		scanner.Scan()
		line = scanner.Text()
		strs = strings.Fields(line)
		nome := invertedIndex.Read(binManager.FILES_PATH, "nome", strs...)

		fmt.Printf("nomeJap\n>")
		scanner.Scan()
		line = scanner.Text()
		strs = strings.Fields(line)

		nomeJap := invertedIndex.Read(binManager.FILES_PATH, "nomeJap", strs...)

		fmt.Printf("especie\n>")
		scanner.Scan()
		line = scanner.Text()
		strs = strings.Fields(line)
		especie := invertedIndex.Read(binManager.FILES_PATH, "especie", strs...)

		fmt.Printf("tipo\n>")
		scanner.Scan()
		line = scanner.Text()
		strs = strings.Fields(line)
		tipo := invertedIndex.Read(binManager.FILES_PATH, "tipo", strs...)

		fmt.Printf("descricao\n>")
		scanner.Scan()
		line = scanner.Text()
		strs = strings.Fields(line)
		descricao := invertedIndex.Read(binManager.FILES_PATH, "descricao", strs...)

		scoredDocuments := invertedIndex.Merge(nome, nomeJap, especie, tipo, descricao)

		for _, scoredDocument := range scoredDocuments {
			pokeAddress, _ := hashing.HashRead(scoredDocument.DocumentID, binManager.FILES_PATH, "hashIndex")
			tmpPoke := binManager.ReadTargetPokemon(pokeAddress)
			fmt.Printf("ID = %3d  |  Nome = %30s  |  compatibilidade = %d\n", tmpPoke.Numero, tmpPoke.Nome, scoredDocument.Score)
		}
    case "7":
        controler, _ := binManager.InicializarControleLeitura(binManager.BIN_FILE)
        bplustree.StartBPlusTreeFile(binManager.FILES_PATH, "altura", controler)
    
    case "8":
        tree, _ := bplustree.ReadBPlusTree(binManager.FILES_PATH, "numero")
        tree.FindRange(256, 300)
        tree.PrintFile()
        
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
	http.HandleFunc("/invertedIndex/", middlewares.EnableCORS(handlers.InvertedIndex))

	// Ordenação externa
	http.HandleFunc("/ordenacao/", middlewares.EnableCORS(handlers.Ordenacao))

	// Inicializa o servidor HTTP na porta 8080 e escreve no log eventuais erros
	logger.Fatal(http.ListenAndServe(":8080", nil))
}
