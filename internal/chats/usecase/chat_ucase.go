package usecase

import (
	"github.com/2020_1_Skycode/internal/chats"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/tools/supportChat"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ChatUseCase struct {
	sC *supportChat.ChatServer
	sR chats.Repository
}

func NewChatUseCase(sR chats.Repository) *ChatUseCase {
	return &ChatUseCase{
		sC: supportChat.NewChatServer(),
		sR: sR,
	}
}

func (cU *ChatUseCase) StartChat(w http.ResponseWriter, r *http.Request) (*websocket.Conn, *supportChat.JoinStatus, error) {
	return cU.sC.CreateChat(w, r)
}

func (cU *ChatUseCase) FindChat(w http.ResponseWriter, r *http.Request, chatID string) (*websocket.Conn, *supportChat.JoinStatus, error) {
	return cU.sC.SearchChat(w, r, chatID)
}

func (cU *ChatUseCase) ReadMessageFromUSer(ws *websocket.Conn) (supportChat.InputMessage, error) {
	message := supportChat.InputMessage{}

	if err := ws.ReadJSON(&message); err != nil {
		return message, err
	}

	return message, nil
}

func (cU *ChatUseCase) WriteFromUserMessage(message supportChat.InputMessage) {
	cU.sC.WriteInputCh(message)
}

func (cU *ChatUseCase) JoinUserToChat(conn *websocket.Conn, userID uint64, userName string, chatID string) error {
	err := cU.sC.JoinUser(conn, &supportChat.JoinStatus{
		UserID: userID,
		UserName: userName,
		ChatID: chatID,
	})

	if err != nil {
		return err
	}

	return nil
}

func (cU *ChatUseCase) LeaveUserChat(chatID string) error {
	if err := cU.sC.LeaveUser(chatID); err != nil {
		return err
	}

	return nil
}

func (cU *ChatUseCase) JoinSupportToChat(conn *websocket.Conn, userID uint64, userName string, chatID string) error {
	err := cU.sC.JoinSupport(conn, &supportChat.JoinStatus{
		UserID: userID,
		UserName: userName,
		ChatID: chatID,
	})

	if err != nil {
		return err
	}

	return nil
}

func (cU *ChatUseCase) LeaveSupportChat(chatID string) error {
	if err := cU.sC.LeaveSupport(chatID); err != nil {
		return err
	}

	return nil
}

func (cU *ChatUseCase) GetChats() []*models.Chat {
	chats := []*models.Chat{}

	supChats := cU.sC.GetSupportChats()

	logrus.Error(supChats)
	if len(supChats) == 0 {
		return chats
	}

	for ind, val := range supChats {
		logrus.Error(ind, val)
		chat := &models.Chat{
			UserName: val.User.UserName,
			UserID: val.User.UserID,
			ChatID: ind,
		}
		chats = append(chats, chat)
	}

	return chats
}

func (cU *ChatUseCase) StoreMessage(message *models.ChatMessage) error {
	if err := cU.sR.InsertChatMessage(message); err != nil {
		return err
	}

	return nil
}

func (cU *ChatUseCase) GetChatMessages(chatID string) ([]*models.ChatMessage, error) {
	var messages []*models.ChatMessage
	messages, err := cU.sR.SelectMessagesByChatID(chatID)

	if err != nil {
		return messages, err
	}

	return messages, nil
}

func (cU *ChatUseCase) GetChat(chatID string) *supportChat.SupportChat {
	return cU.sC.GetChat(chatID)
}
