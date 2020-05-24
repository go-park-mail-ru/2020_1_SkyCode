package chats

import (
	"github.com/2020_1_Skycode/internal/models"
)

type UseCase interface {
	StoreMessage(message *models.ChatMessage) error
	GetChatMessages(chatID uint64) ([]*models.ChatMessage, error)
}
