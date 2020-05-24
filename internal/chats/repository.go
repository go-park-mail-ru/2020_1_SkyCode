package chats

import "github.com/2020_1_Skycode/internal/models"

type Repository interface {
	InsertChatMessage(message *models.ChatMessage) error
	SelectMessagesByChatID(chatID uint64) ([]*models.ChatMessage, error)
}
