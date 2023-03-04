package main

import (
	"fmt"

	"github.com/Bernardo46-2/AEDS-III/crud"
	"github.com/Bernardo46-2/AEDS-III/dataManager"
	"github.com/Bernardo46-2/AEDS-III/utils"
)

func main() {
	menu := "1 - Create\n2 - Read\n3 - Update\n4 - Delete\n8 - Convert CSV to Bin\n0 - Exit\n"

	for quit := false; !quit; {
		switch utils.LerInt(menu) {
		case 0:
			fmt.Printf("\nSaindo do programa...\n\n")
			quit = true
		case 1:
			crud.Create()
		case 2:
			crud.Read()
		case 3:
			crud.Update()
		case 4:
			crud.Delete()
		case 8:
			dataManager.ImportCSV().CsvToBin()
		default:
			fmt.Println("Opção inválida")
		}
		utils.Pause()
	}
}
