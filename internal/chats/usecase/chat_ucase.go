package usecase

import (
	"github.com/2020_1_Skycode/internal/chats"
	"github.com/2020_1_Skycode/internal/models"
)

type ChatUseCase struct {
	sR chats.Repository
}

func NewChatUseCase(sR chats.Repository) chats.UseCase {
	return &ChatUseCase{
		sR: sR,
	}
}

func (cU *ChatUseCase) StoreMessage(message *models.ChatMessage) error {
	if err := cU.sR.InsertChatMessage(message); err != nil {
		return err
	}

	return nil
}

func (cU *ChatUseCase) GetChatMessages(chatID uint64) ([]*models.ChatMessage, error) {
	var messages []*models.ChatMessage
	messages, err := cU.sR.SelectMessagesByChatID(chatID)

	if err != nil {
		return messages, err
	}

	return messages, nil
}
