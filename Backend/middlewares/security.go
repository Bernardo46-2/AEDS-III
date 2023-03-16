// O pacote middlewares contém funções que são utilizadas como intermediárias
// no processamento de requisições HTTP em um servidor web. Essas funções são
// usadas para realizar autenticação, autorização, manipulação de cookies e
// cabeçalhos, entre outras. Ao utilizar as funções deste pacote.
package middlewares

import (
	"net/http"
)

// EnableCORS é uma função intermediaria para regularização do sistema de
// Cross-Origin Resource Sharing
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
