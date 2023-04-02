package ordenacao

import "github.com/Bernardo46-2/AEDS-III/data/binManager"

const FILE string = binManager.FILE
const BIN_FILE string = binManager.BIN_FILE

type SortFunc func()

var SortingFunctions = []SortFunc{
	IntercalacaoBalanceadaComum,
	IntercalacaoBalanceadaVariavel,
	IntercalacaoBalanceadaVariavel,
}
