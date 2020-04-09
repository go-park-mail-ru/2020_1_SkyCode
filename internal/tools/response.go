package tools

import "github.com/2020_1_Skycode/internal/models"

type Body map[string] interface{}

type Error struct {
	ErrorMessage string `json:"error"`
}

type Message struct {
	Message interface{} `json:"message"`
}

type UserMessage struct {
	User *models.User
}
