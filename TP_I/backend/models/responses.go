package models

import "net/http"

type Response struct {
	Success bool   `json:"sucesso"`
	Code    int    `json:"codigo"`
	Message string `json:"mensagem"`
}

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
		msg = "CSV importado com sucesso!"
	case 7:
		msg = "Intercalacao Balanceada Comum realizada com sucesso!"
	case 8:
		msg = "Intercalacao Balanceada Variavel realizada com sucesso!"
	case 9:
		msg = "Intercalacao Por Substituicao realizada com sucesso!"
	default:
		msg = "Mensagem de sucesso desconhecida!"
	}

	return Response{Success: true, Code: codigo, Message: msg}
}
