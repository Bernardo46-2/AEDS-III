package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func runCmd(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func clearScreen() {
	switch runtime.GOOS {
	case "darwin":
		runCmd("clear")
	case "linux":
		runCmd("clear")
	case "windows":
		runCmd("cmd", "/c", "cls")
	default:
		runCmd("clear")
	}
}

func pause() {
	var input string
	fmt.Printf("\nPressione Enter para continuar...\n")
	fmt.Scanf("%s\n", &input)
}

func RemoveAfterSpace(str string) string {
	parts := strings.Split(str, " ")
	return parts[0]
}

func lerInt(str string) int {
	clearScreen()
	fmt.Printf("%s\n> ", str)

	var tmp string
	var result int
	var err error

	if _, err = fmt.Scanln(&tmp); err != nil {
		fmt.Println("Erro ao ler opção:", err)
		pause()
		result = lerInt(str)
	} else {
		if result, err = strconv.Atoi(tmp); err != nil {
			fmt.Println("Erro ao ler opção:", err)
			pause()
			result = lerInt(str)
		}
	}
	return result
}

func main() {
	var csvFile CSV
	menu := "1 - Create\n2 - Read\n3 - Update\n4 - Delete\n8 - Convert CSV to Bin\n0 - Exit\n"

	for quit := false; !quit; {
		switch lerInt(menu) {
		case 0:
			fmt.Printf("Saindo do programa...\n\n")
			quit = true
		case 1:
			fmt.Printf("Create\n")
			incrementNumRegistros()
		case 2:
			pokemon, _ := readBinToPoke(lerInt("Digite o Id a pesquisar:\n"))
			fmt.Printf(pokemon.ToString())
		case 3:
			fmt.Printf("Update\n")
		case 4:
			fmt.Printf("Delete\n")
		case 8:
			csvFile = importCSV()
			csvFile.CsvToBin()
		default:
			fmt.Println("Opção inválida")
		}
		pause()
	}
}
