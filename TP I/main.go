package main

import (
    "fmt"
)

func main() {
    var csvFile CSV
    menu := "1 - Create\n2 - Read\n3 - Update\n4 - Delete\n8 - Convert CSV to Bin\n0 - Exit\n"

    for quit := false; !quit; {
        switch lerInt(menu) {
        case 0:
            fmt.Printf("\nSaindo do programa...\n\n")
            quit = true
        case 1:
            create()
        case 2:
            read()
        case 3:
            update()
        case 4:
            delete()
        case 8:
            csvFile = importCSV()
            csvFile.CsvToBin()
        default:
            fmt.Println("Opção inválida")
        }
        pause()
    }
}
