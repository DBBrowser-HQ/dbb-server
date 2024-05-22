package model

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Payload any    `json:"payload"`
}
