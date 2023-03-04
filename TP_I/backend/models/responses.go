package models

import "net/http"

type Error struct {
	Codigo  int    `json:"codigo"`
	Message string `json:"error"`
}

func ErrorResponse(codigo int) Error {
	var msg string

	switch codigo {
	case 1:
		msg = "Sucesso"
	case 2:
		msg = "Erro ao buscar o Pokemon"
	case 3:
		msg = "Erro ao criar o Pokemon"
	case 4:
		msg = "Erro ao atualizar o Pokemon"
	case 5:
		msg = "Erro ao deletar o Pokemon"
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

	return Error{Codigo: codigo, Message: msg}
}

type Success struct {
	Codigo  int    `json:"codigo"`
	Message string `json:"message"`
}

func SuccessResponse(codigo int) Success {
	var msg string

	switch codigo {
	case 1:
		msg = "Sucesso"
	case 2:
		msg = "Pokemon Atualizado com sucesso"
	case 3:
		msg = "Pokemon Deletado com sucesso"
	case 4:
		msg = "Banco de dados carregado com sucesso!"
	default:
		msg = "Mensagem de sucesso desconhecida"
	}

	return Success{Codigo: codigo, Message: msg}
}
