package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Bernardo46-2/AEDS-III/dataManager"
	"github.com/Bernardo46-2/AEDS-III/handlers"
	"github.com/Bernardo46-2/AEDS-III/logger"
	"github.com/Bernardo46-2/AEDS-III/middlewares"
)

func main() {

	var command int

	fmt.Printf("0 = realizar comando\n1 = ligar server\n>")
	fmt.Scanln(&command)
	if command == 0 {
		dataManager.IntercalacaoBalanceadaComum()
	} else if command == 1 {
        dataManager.IntercalacaoTamanhoVariavel()
    } else if command == 2 {
		logger.LigarServidor()

		// define os handlers para GET e POST
		http.HandleFunc("/getAll/", middlewares.EnableCORS(handlers.GetAllPokemon))
		http.HandleFunc("/get/", middlewares.EnableCORS(handlers.GetPokemon))
		http.HandleFunc("/post/", middlewares.EnableCORS(handlers.PostPokemon))
		http.HandleFunc("/put/", middlewares.EnableCORS(handlers.PutPokemon))
		http.HandleFunc("/delete/", middlewares.EnableCORS(handlers.DeletePokemon))
		http.HandleFunc("/loadDatabase", middlewares.EnableCORS(handlers.LoadDatabase))
		http.HandleFunc("/toKatakana/", middlewares.EnableCORS(handlers.ToKatakana))

		log.Fatal(http.ListenAndServe(":8080", nil))
	}
}
