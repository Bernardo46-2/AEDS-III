package models

import "net/http"

type Response struct {
	Success bool   `json:"sucesso"`
	Codigo  int    `json:"codigo"`
	Message string `json:"message"`
}

func ErrorResponse(codigo int) Response {
	var msg string

	switch codigo {
	case 1:
		msg = "Erro"
	case 2:
		msg = "Erro ao buscar o Pokemon"
	case 3:
		msg = "Erro ao criar o Pokemon"
	case 4:
		msg = "Erro ao atualizar o Pokemon"
	case 5:
		msg = "Erro ao deletar o Pokemon"
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

	return Response{Success: false, Codigo: codigo, Message: msg}
}

func SuccessResponse(codigo int) Response {
	var msg string

	switch codigo {
	case 1:
		msg = "Sucesso"
	case 2:
		msg = "Pokemon encontrado com sucesso"
	case 3:
		msg = "Pokemon criado com sucesso"
	case 4:
		msg = "Pokemon atualizado com sucesso"
	case 5:
		msg = "Pokemon deletado com sucesso"
	case 6:
		msg = "CSV importado com sucesso"
	default:
		msg = "Mensagem de sucesso desconhecida"
	}

	return Response{Success: true, Codigo: codigo, Message: msg}
}
