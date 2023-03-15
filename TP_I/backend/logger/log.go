package logger

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
)

func Println(logType string, message string) {
	file, err := os.OpenFile("logger/log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Erro ao abrir arquivo de log:", err)
		fmt.Println("Fechando servidor . . .")
		os.Exit(0)
		return
	}
	defer file.Close()

	log.SetOutput(file)
	log.SetPrefix(strings.ToUpper(logType) + ": ")
	log.Println(message)
}

func LigarServidor() {
	Println("STATUS", "Servidor iniciado")
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)
		<-sigChan
		Println("STATUS", "Servidor desligado")
		os.Exit(0)
	}()
}

func Fatal(v ...any) {
	Println(fmt.Sprint(v...), "ERROR")
	os.Exit(1)
}
