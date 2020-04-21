package chats

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/tools/supportChat"
	"github.com/gorilla/websocket"
	"net/http"
)

type UseCase interface {
	StartChat(w http.ResponseWriter, r *http.Request) (*websocket.Conn, *supportChat.JoinStatus, error)
	FindChat(w http.ResponseWriter, r *http.Request, chatID string) (*websocket.Conn, *supportChat.JoinStatus, error)
	JoinUserToChat(conn *websocket.Conn, fullName string, chatID string) error
	LeaveUserChat(chatID string) error
	JoinSupportToChat(conn *websocket.Conn, fullName string, chatID string) error
	LeaveSupportChat(chatID string) error
	ReadMessageFromUSer(ws *websocket.Conn) (supportChat.InputMessage, error)
	WriteFromUserMessage(message supportChat.InputMessage)
	GetChats() []*models.Chat
}
