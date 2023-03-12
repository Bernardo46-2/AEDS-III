package main

import (
	"log"
	"net/http"

	"github.com/Bernardo46-2/AEDS-III/handlers"
	"github.com/Bernardo46-2/AEDS-III/logger"
	"github.com/Bernardo46-2/AEDS-III/middlewares"
)

func main() {

	logger.LigarServidor()

	// define os handlers para GET e POST
	http.HandleFunc("/getAll/", middlewares.EnableCORS(handlers.GetAllPokemon))
	http.HandleFunc("/get/", middlewares.EnableCORS(handlers.GetPokemon))
	http.HandleFunc("/post/", middlewares.EnableCORS(handlers.PostPokemon))
	http.HandleFunc("/put/", middlewares.EnableCORS(handlers.PutPokemon))
	http.HandleFunc("/delete/", middlewares.EnableCORS(handlers.DeletePokemon))
	http.HandleFunc("/loadDatabase", middlewares.EnableCORS(handlers.LoadDatabase))
	http.HandleFunc("/toKatakana/", middlewares.EnableCORS(handlers.ToKatakana))
	http.HandleFunc("/intercalacaoComum/", middlewares.EnableCORS(handlers.IntercalacaoComum))
	http.HandleFunc("/intercalacaoVariavel/", middlewares.EnableCORS(handlers.IntercalacaoVariavel))
	http.HandleFunc("/selecaoPorSubstituicao/", middlewares.EnableCORS(handlers.SelecaoPorSubstituicao))

	log.Fatal(http.ListenAndServe(":8080", nil))

}
