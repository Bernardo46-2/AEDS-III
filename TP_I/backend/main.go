package main

import (
	"log"
	"net/http"

	"github.com/Bernardo46-2/AEDS-III/handlers"
	"github.com/Bernardo46-2/AEDS-III/middlewares"
)

func main() {
	// define os handlers para GET e POST
	http.HandleFunc("/get/", middlewares.EnableCORS(handlers.GetPokemon))
	http.HandleFunc("/post/", middlewares.EnableCORS(handlers.PostPokemon))
	http.HandleFunc("/put/", middlewares.EnableCORS(handlers.PutPokemon))
	http.HandleFunc("/delete/", middlewares.EnableCORS(handlers.DeletePokemon))
	http.HandleFunc("/load-database", middlewares.EnableCORS(handlers.LoadDatabase))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
