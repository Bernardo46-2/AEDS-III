// O pacote logger fornece funções para criar logs da aplicação, permitindo
// registrar as informações relevantes sobre o comportamento do software em
// diferentes níveis de gravidade.
// Ele também oferece a possibilidade de personalizar o formato do log
package logger

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
)

// Println recebe um cabeçalho de resposta (STATUS, INFO, CRUD, ...) e uma mensagem e
// formata para escrita em arquivo de log
func Println(logType string, message string) {
	// Abre o arquivo para append, caso nao exista cria
	file, err := os.OpenFile("logger/log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Erro ao abrir arquivo de log:", err)
		fmt.Println("Fechando servidor . . .")
		os.Exit(0)
		return
	}
	defer file.Close()

	// Formatação e escrita
	log.SetOutput(file)
	log.SetPrefix(strings.ToUpper(logType) + ": ")
	log.Println(message)
}

// LigarServidor registra no logger a inicialização do servidor e intercepta mensagens
// do sistema operacional sobre comandos de erro ou fechamento atraves de uma rotina
// e salva no log
func LigarServidor() {
	// Oficializa a ligação no log
	Println("STATUS", "Servidor iniciado")

	// Função de rotina para interceptação do desligamento do servidor
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)
		<-sigChan
		Println("STATUS", "Servidor desligado")
		os.Exit(0)
	}()
}

// Fatal intercepta erros fatais de qualquer especie e formaliza no log
func Fatal(v ...any) {
	Println(fmt.Sprint(v...), "ERROR")
	os.Exit(1)
}
