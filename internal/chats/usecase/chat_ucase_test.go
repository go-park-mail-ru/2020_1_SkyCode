package usecase

import (
	mock_chats "github.com/2020_1_Skycode/internal/chats/mocks"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestChatUseCase_GetChatMessages(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockChatsRepo := mock_chats.NewMockRepository(ctrl)

	chatID := "testid"

	resultMessages := []*models.ChatMessage{
		{UserID: 1, UserName: "test", ChatID: chatID, Message: "Che s den'gami?", Created: time.Now()},
	}

	mockChatsRepo.EXPECT().SelectMessagesByChatID(chatID).Return(resultMessages, nil)
	chatUCase := NewChatUseCase(mockChatsRepo)

	messages, err := chatUCase.GetChatMessages(chatID)
	require.NoError(t, err)

	require.EqualValues(t, resultMessages, messages)
}

func TestChatUseCase_StoreMessage(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockChatsRepo := mock_chats.NewMockRepository(ctrl)

	testMessage := &models.ChatMessage{
		UserID:   1,
		UserName: "test",
		ChatID:   "testid",
		Message:  "Che s den'gami?",
	}

	mockChatsRepo.EXPECT().InsertChatMessage(testMessage).Return(nil)
	chatUCase := NewChatUseCase(mockChatsRepo)

	err := chatUCase.StoreMessage(testMessage)
	require.NoError(t, err)
}
