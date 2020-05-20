package supportChat

import (
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

type chatMember struct {
	UserID uint64
	UserName string
	ws     *websocket.Conn
}

func (cm *chatMember) CloseConn() error {
	if cm.ws != nil {
		err := cm.ws.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

type SupportChat struct {
	User    *chatMember
	Support *chatMember
}

func (sc *SupportChat) NotifyMembers(message interface{}) error {
	if user := sc.User; user != nil {
		if err := user.ws.WriteJSON(message); err != nil {
			return err
		}
	}

	if support := sc.Support; support != nil {
		if err := support.ws.WriteJSON(message); err != nil {
			return err
		}
	}

	return nil
}

func (sc *SupportChat) Dead() bool {
	if sc.Support == nil && sc.User == nil {
		return true
	}

	return false
}

type InputMessage struct {
	ChatID   string `json:"chat_id"`
	UserID uint64 `json:"user_id"`
	UserName string `json:"user_name"`
	Message  string `json:"message"`
}

type JoinStatus struct {
	ChatID   string `json:"chat_id,omitempty"`
	UserID uint64 `json:"user_id"`
	UserName string `json:"user_name"`
	Joined   bool   `json:"joined"`
}

type LeaveStatus struct {
	ChatID   string `json:"chat_id"`
	UserID uint64 `json:"user_id"`
	UserName string `json:"user_name"`
	Leaved   bool   `json:"leaved"`
}

type ChatServer struct {
	supportChats map[string]*SupportChat
	inputCh      chan InputMessage
	joinCh       chan JoinStatus
	leaveCh      chan LeaveStatus
	upd          *websocket.Upgrader
}

func NewChatServer() *ChatServer {
	cS := &ChatServer{
		supportChats: make(map[string]*SupportChat),
		inputCh:      make(chan InputMessage),
		joinCh:       make(chan JoinStatus),
		leaveCh:      make(chan LeaveStatus),
		upd: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	go cS.handleMessages()

	return cS
}

func (cs *ChatServer) GetChat(chatID string) *SupportChat {
	return cs.supportChats[chatID]
}

func (cs *ChatServer) WriteInputCh(message InputMessage) {
	cs.inputCh <- message
}

func (cs *ChatServer) WriteJoinCh(message JoinStatus) {
	cs.joinCh <- message
}

func (cs *ChatServer) WriteLeaveCh(message LeaveStatus) {
	cs.leaveCh <- message
}

func (cs *ChatServer) handleMessages() {
	for {
		select {
		case msgIN := <-cs.inputCh:
			logrus.Info(msgIN)
			if chat := cs.supportChats[msgIN.ChatID]; chat != nil {
				if err := cs.supportChats[msgIN.ChatID].NotifyMembers(msgIN); err != nil {
					logrus.Error(err)
					break
				}
			}

		case joinIN := <-cs.joinCh:
			logrus.Info(joinIN)
			if chat := cs.supportChats[joinIN.ChatID]; chat != nil {
				if err := cs.supportChats[joinIN.ChatID].NotifyMembers(joinIN); err != nil {
					logrus.Error(err)
					break
				}
			}


		case leaveIN := <-cs.leaveCh:
			logrus.Info(leaveIN)
			if chat := cs.supportChats[leaveIN.ChatID]; chat != nil {
				if err := cs.supportChats[leaveIN.ChatID].NotifyMembers(leaveIN); err != nil {
					logrus.Error(err)
					break
				}
			}
		}
	}
}

func (cs *ChatServer) CreateChat(w http.ResponseWriter, r *http.Request) (*websocket.Conn, *JoinStatus, error) {
	ws, err := cs.upd.Upgrade(w, r, nil)

	if err != nil {
		return nil, nil, err
	}

	joinMessage := &JoinStatus{}

	if err := ws.ReadJSON(&joinMessage); err != nil {
		return ws, nil, err
	}

	if joinMessage.ChatID != "" {
		if cs.supportChats[joinMessage.ChatID] != nil {
			return ws, joinMessage, nil
		}
	}

	chatID := uuid.New().String()
	cs.supportChats[chatID] = &SupportChat{}

	joinMessage.ChatID = chatID
	return ws, joinMessage, nil
}

func (cs *ChatServer) SearchChat(w http.ResponseWriter, r *http.Request, chatID string) (*websocket.Conn, *JoinStatus, error) {
	ws, err := cs.upd.Upgrade(w, r, nil)

	if err != nil {
		return nil, nil, err
	}

	joinMessage := &JoinStatus{}

	if err := ws.ReadJSON(&joinMessage); err != nil {
		return ws, nil, err
	}

	if chatID == "" {
		return ws, nil, errors.New("chat id not presented")
	}

	if cs.supportChats[chatID] == nil {
		return ws, joinMessage, errors.New("chat not found")
	}

	joinMessage.ChatID = chatID

	return ws, joinMessage, nil
}

func (cs *ChatServer) JoinUser(conn *websocket.Conn, jM *JoinStatus) error {
	var chat *SupportChat

	if chat = cs.supportChats[jM.ChatID]; chat == nil {
		return errors.New("chat not found")
	}

	cs.supportChats[jM.ChatID].User = &chatMember{
		UserID: jM.UserID,
		UserName: jM.UserName,
		ws:       conn,
	}

	if err := cs.supportChats[jM.ChatID].NotifyMembers(&JoinStatus{
		ChatID:   jM.ChatID,
		UserID: jM.UserID,
		UserName: jM.UserName,
		Joined:   true,
	}); err != nil {
		return err
	}

	return nil
}

func (cs *ChatServer) LeaveUser(chatID string) error {
	var chat *SupportChat

	if chat = cs.supportChats[chatID]; chat == nil {
		return errors.New("chat not found")
	}

	if err := chat.User.CloseConn(); err != nil {
		return err
	}

	chat.User = nil

	if chat.Dead() {
		delete(cs.supportChats, chatID)
	}

	return nil
}

func (cs *ChatServer) JoinSupport(conn *websocket.Conn, jM *JoinStatus) error {
	var chat *SupportChat

	if chat = cs.supportChats[jM.ChatID]; chat == nil {
		return errors.New("chat not found")
	}

	cs.supportChats[jM.ChatID].Support = &chatMember{
		UserID: jM.UserID,
		UserName: jM.UserName,
		ws:       conn,
	}

	if err := cs.supportChats[jM.ChatID].NotifyMembers(&JoinStatus{
		ChatID:   jM.ChatID,
		UserID: jM.UserID,
		UserName: jM.UserName,
		Joined:   true,
	}); err != nil {
		return err
	}

	return nil
}

func (cs *ChatServer) LeaveSupport(chatID string) error {
	var chat *SupportChat

	if chat = cs.supportChats[chatID]; chat == nil {
		return errors.New("chat not found")
	}

	if err := chat.Support.CloseConn(); err != nil {
		return err
	}

	chat.Support = nil

	if chat.Dead() {
		delete(cs.supportChats, chatID)
	}

	return nil
}

func (cs *ChatServer) GetSupportChats() map[string]*SupportChat {
	return cs.supportChats
}
