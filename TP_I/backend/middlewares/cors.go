package middlewares

import "net/http"

func EnableCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// adiciona o header de Access-Control-Allow-Origin para permitir todos os origens
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// verifica se a requisição é um OPTIONS (pré-voo)
		if r.Method == "OPTIONS" {
			return
		}

		// chama o handler fornecido
		handler(w, r)
	}
}
