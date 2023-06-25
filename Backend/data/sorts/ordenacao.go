// o arquivo sorts permite a realizacao de ordenacao externa na base de dados binaria
// atraves do algoritmo de intercalacao e suas variantes
//
// # Ordenacoes possiveis:
//
//	Intercalacao Comum
//	Intercalacao Por Bloco De Tamanho Variavel
//	Intercalacao Por Substituicao (Heap)
package sorts

import "github.com/Bernardo46-2/AEDS-III/data/binManager"

// Constante para gerenciamento do banco de dados
const (
	FILE     string = binManager.CSV_PATH
	BIN_FILE string = binManager.BIN_FILE
)

// Interface para criacao de hash de funcoes de ordenacao
type SortFunc func()

// Hash de funcoes para direcionamento do respectivo metodo
var SortingFunctions = []SortFunc{
	IntercalacaoBalanceadaComum,
	IntercalacaoBalanceadaVariavel,
	IntercalacaoPorSubstituicao,
}
