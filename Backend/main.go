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
		btree.StartBTreeFile()
	case "5", "criarIndiceInvertido":
		c1, _ := binManager.InicializarControleLeitura(binManager.BIN_FILE)
		invertedIndex.New(c1, "nome", binManager.FILES_PATH)
		c1.Close()
		c2, _ := binManager.InicializarControleLeitura(binManager.BIN_FILE)
		invertedIndex.New(c2, "nome_jap", binManager.FILES_PATH)
		c2.Close()
		c3, _ := binManager.InicializarControleLeitura(binManager.BIN_FILE)
		invertedIndex.New(c3, "especie", binManager.FILES_PATH)
		c3.Close()
		c4, _ := binManager.InicializarControleLeitura(binManager.BIN_FILE)
		invertedIndex.New(c4, "tipo", binManager.FILES_PATH)
		c4.Close()
		c5, _ := binManager.InicializarControleLeitura(binManager.BIN_FILE)
		invertedIndex.New(c5, "descricao", binManager.FILES_PATH)
		c5.Close()
	case "6", "lerMultiplosIndiceInvertido":
		var strs []string
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "exit" {
				break
			}
			words := strings.Fields(line)
			strs = append(strs, words...)
		}

		ids := invertedIndex.Read(binManager.FILES_PATH, "descricao", strs...)
		for i := 0; i < len(ids); i++ {
			pokeAddress, _ := hashing.HashRead(ids[i], binManager.FILES_PATH, "hashIndex")
			tmpPoke := binManager.ReadTargetPokemon(pokeAddress)
			fmt.Printf("ID = %3d | Nome = %s | descricao = %s\n", tmpPoke.Numero, tmpPoke.Nome, tmpPoke.Descricao)
		}
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
	http.HandleFunc("/getAll/", middlewares.EnableCORS(handlers.GetAllPokemon))
	http.HandleFunc("/get/", middlewares.EnableCORS(handlers.GetPokemon))
	http.HandleFunc("/post/", middlewares.EnableCORS(handlers.PostPokemon))
	http.HandleFunc("/put/", middlewares.EnableCORS(handlers.PutPokemon))
	http.HandleFunc("/delete/", middlewares.EnableCORS(handlers.DeletePokemon))
	http.HandleFunc("/loadDatabase", middlewares.EnableCORS(handlers.LoadDatabase))
	http.HandleFunc("/toKatakana/", middlewares.EnableCORS(handlers.ToKatakana))

	// Ordenação externa
	http.HandleFunc("/ordenacao/", middlewares.EnableCORS(handlers.Ordenacao))

	// Inicializa o servidor HTTP na porta 8080 e escreve no log eventuais erros
	logger.Fatal(http.ListenAndServe(":8080", nil))
}
