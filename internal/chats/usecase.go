package chats

import (
	"github.com/2020_1_Skycode/internal/tools/supportChat"
	"github.com/gorilla/websocket"
	"net/http"
)

type UseCase interface {
	StartChat() string
	JoinUserToChat(w http.ResponseWriter, r *http.Request, chatID string) (*websocket.Conn, error)
	LeaveUserChat(chatID string) error
	JoinSupportToChat(w http.ResponseWriter, r *http.Request, chatID string) (*websocket.Conn, error)
	LeaveSupportChat(chatID string) error
	ReadMessageFromUSer(ws *websocket.Conn) (supportChat.InputMessage, error)
	WriteFromUserMessage(message supportChat.InputMessage)
}
