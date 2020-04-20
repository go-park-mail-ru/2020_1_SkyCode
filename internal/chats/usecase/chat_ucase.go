package usecase

import (
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

func (cU *ChatUseCase) StartChat() string {
	return cU.sC.CreateChat()
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

func (cU *ChatUseCase) JoinUserToChat(w http.ResponseWriter, r *http.Request,
	chatID string) (*websocket.Conn, error) {
	ws, err := cU.sC.JoinUser(w, r, chatID)

	if err != nil {
		return nil, err
	}

	return ws, err
}

func (cU *ChatUseCase) LeaveUserChat(chatID string) error {
	if err := cU.sC.LeaveUser(chatID); err != nil {
		return err
	}

	return nil
}

func (cU *ChatUseCase) JoinSupportToChat(w http.ResponseWriter, r *http.Request,
	chatID string) (*websocket.Conn, error) {
	ws, err := cU.sC.JoinSupport(w, r, chatID)

	if err != nil {
		return nil, err
	}

	return ws, err
}

func (cU *ChatUseCase) LeaveSupportChat(chatID string) error {
	if err := cU.sC.LeaveSupport(chatID); err != nil {
		return err
	}

	return nil
}
