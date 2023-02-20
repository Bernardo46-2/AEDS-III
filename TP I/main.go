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

func main() {
    var csvFile CSV

    for quit := false; !quit;{
        clearScreen()
        fmt.Printf("1 - Create\n")
        fmt.Printf("2 - Read\n")
        fmt.Printf("3 - Update\n")
        fmt.Printf("4 - Delete\n")
        fmt.Printf("8 - Convert CSV to Bin\n")
        fmt.Printf("9 - Import CSV\n")
        fmt.Printf("0 - Exit\n")
        fmt.Printf("\n> ")

        var tmp string
        if _, err := fmt.Scanln(&tmp); err != nil {
            fmt.Println("Erro ao ler opção:", err)
        }
        if opcao, err := strconv.Atoi(tmp); err == nil {
            switch opcao {
            case 0:
                fmt.Printf("Saindo do programa...\n\n")
                quit = true
            case 1:
                fmt.Printf("Create\n")
            case 2:
                fmt.Printf("Read\n")
            case 3:
                fmt.Printf("Update\n")
            case 4:
                fmt.Printf("Delete\n")
            case 8:
                csvFile.CsvToBin()
            case 9:
                csvFile = importCSV()
            default:
                fmt.Println("Opção inválida")
            }
        } else {
            panic(fmt.Errorf("Erro ao ler opção: %v", err))
        }
        pause()
    }
}
