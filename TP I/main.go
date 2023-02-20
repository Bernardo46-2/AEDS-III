package main

import (
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

func clearScreen() {
    cmd := exec.Command("cmd", "/c", "cls")
    cmd.Stdout = os.Stdout
    cmd.Run()
}

func pause() {
    var input string
    fmt.Printf("\nPressione Enter para continuar...\n")
    fmt.Scanf("%s\n", &input)
}

func removeAfterSpace(str string) string {
    parts := strings.Split(str, " ")
    return parts[0]
}

func main() {
    quit := false
    var csvFile CSV
    
    for !quit {
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
                    break
                case 1:
                    fmt.Printf("Create\n")
                    break
                case 2:
                    fmt.Printf("Read\n")
                    break
                case 3:
                    fmt.Printf("Update\n")
                    break
                case 4:
                    fmt.Printf("Delete\n")
                    break
                case 8:
                    csvFile.CsvToBin()
                    break
                case 9:
                    csvFile = importCSV()
                    break
                default:
                    fmt.Println("Opção inválida")
                    break
            }
        } else {
            panic(fmt.Errorf("Erro ao ler opção: %v", err))
        }
        pause()
    }
}
