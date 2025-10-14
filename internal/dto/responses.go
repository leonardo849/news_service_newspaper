package dto

type MessageDTO struct {
	Message string `json:"string"`
}

type ErrorDTO struct {
	Error string `json:"error"`
}

type CreatedMessage struct {
	Message string `json:"message"`
	Id string `json:"id"`
}