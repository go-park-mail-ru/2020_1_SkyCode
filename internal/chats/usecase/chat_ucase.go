package usecase

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/tools/supportChat"
	"github.com/gorilla/websocket"
	"net/http"
)

type ChatUseCase struct {
	sC *supportChat.ChatServer
}

func NewChatUseCase() *ChatUseCase {
	return &ChatUseCase{
		sC: supportChat.NewChatServer(),
	}
}

func (cU *ChatUseCase) StartChat(w http.ResponseWriter, r *http.Request) (*websocket.Conn, *supportChat.JoinStatus, error) {
	return cU.sC.CreateChat(w, r)
}

func (cU *ChatUseCase) FindChat(w http.ResponseWriter, r *http.Request) (*websocket.Conn, *supportChat.JoinStatus, error) {
	return cU.sC.SearchChat(w, r)
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

func (cU *ChatUseCase) JoinUserToChat(conn *websocket.Conn, fullName string, chatID string) error {
	err := cU.sC.JoinUser(conn, &supportChat.JoinStatus{
		FullName: fullName,
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

func (cU *ChatUseCase) JoinSupportToChat(conn *websocket.Conn, fullName string, chatID string) error {
	err := cU.sC.JoinUser(conn, &supportChat.JoinStatus{
		FullName: fullName,
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

	for ind, val := range cU.sC.GetSupportChats() {
		chat := &models.Chat{
			UserName: val.User.FullName,
			ChatID: ind,
		}
		chats = append(chats, chat)
	}

	return chats
}
