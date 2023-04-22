// O arquivo Responses do pacote Models permite a criação de uma estrutura de
// erro ou sucesso para comunicação agil e efetiva com o servidor.
// Uma mesma estrutura é utilizada para as duas respostas
// E diferentes funções de parsing de erro para mensagem são feitas
// de acordo com a necessidade de cada situação.
//
// Além disso, também permite a definição de códigos de status personalizados,
// bem como a adição de dados extras à resposta.
package models

import "net/http"

// Response gera um padrão de resposta para o frontend
// contendo um bool para dizer se é sucesso ou erro
// codigo da mensagem e a mensagem
type Response struct {
	Success bool   `json:"sucesso"`
	Code    int    `json:"codigo"`
	Message string `json:"mensagem"`
}

// ErrorResponse faz o parse do codigo do erro em sua respectiva mensagem
func ErrorResponse(codigo int) Response {
	var msg string

	switch codigo {
	case 0:
		msg = "Funcao nao implementada"
	case 2:
		msg = "Pokemon nao encontrado"
	case 3:
		msg = "Erro ao criar o Pokemon"
	case 4:
		msg = "Erro ao atualizar o Pokemon"
	case 5:
		msg = "Pokemon nao encontrado"
	case 6:
		msg = "Erro ao importar CSV"
	case http.StatusBadRequest:
		msg = "Erro ao converter o JSON para Pokemon"
	case http.StatusInternalServerError:
		msg = "Erro interno do servidor"
	case http.StatusUnauthorized:
		msg = "Não autorizado"
	case http.StatusNotFound:
		msg = "Recurso não encontrado"
	case http.StatusForbidden:
		msg = "Acesso proibido"
	case http.StatusMethodNotAllowed:
		msg = "Método não permitido"
	case http.StatusRequestTimeout:
		msg = "Tempo limite da solicitação esgotado"
	case http.StatusConflict:
		msg = "Conflito na solicitação"
	case http.StatusUnsupportedMediaType:
		msg = "Tipo de mídia não suportado"
	default:
		msg = "Erro desconhecido"
	}

	return Response{Success: false, Code: codigo, Message: msg}
}

// ErrorResponse faz o parse do codigo de sucesso em sua respectiva mensagem
func SuccessResponse(codigo int) Response {
	var msg string

	switch codigo {
	case 0:
		msg = "Funcao nao implementada"
	case 1:
		msg = "Sucesso!"
	case 2:
		msg = "Pokemon encontrado com sucesso!"
	case 3:
		msg = "Pokemon criado com sucesso!"
	case 4:
		msg = "Pokemon atualizado com sucesso!"
	case 5:
		msg = "Pokemon deletado com sucesso!"
	case 6:
		msg = "CSV importado! <br>Ìndice Hash Criado! <br>Árvore B Criada! <br>Árvore B+ Criada! <br>Ìndice Invertido Criado!"
	case 7:
		msg = "Ordenacao realizada com sucesso!"
	default:
		msg = "Mensagem de sucesso desconhecida!"
	}

	return Response{Success: true, Code: codigo, Message: msg}
}
