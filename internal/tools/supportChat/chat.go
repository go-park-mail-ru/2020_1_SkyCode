package supportChat

import (
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

type chatMember struct {
	FullName string
	ws       *websocket.Conn
}

func (cm *chatMember) CloseConn() error {
	return cm.ws.Close()
}

type supportChat struct {
	User    *chatMember
	Support *chatMember
}

func (sc *supportChat) NotifyMembers(message interface{}) error {
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

func (sc *supportChat) Dead() bool {
	if sc.Support == nil && sc.User == nil {
		return true
	}

	return false
}

type InputMessage struct {
	ChatID   string `json:"chat_id"`
	FullName string `json:"full_name"`
	Message  string `json:"message"`
}

type JoinStatus struct {
	ChatID   string `json:"chat_id, omitempty"`
	FullName string `json:"full_name"`
	Joined   bool   `json:"joined"`
}

type LeaveStatus struct {
	ChatID   string `json:"chat_id"`
	FullName string `json:"full_name"`
	Leaved   bool   `json:"leaved"`
}

type ChatServer struct {
	supportChats map[string]*supportChat
	inputCh      chan InputMessage
	joinCh       chan JoinStatus
	leaveCh      chan LeaveStatus
	upd          *websocket.Upgrader
}

func NewChatServer() *ChatServer {
	cS := &ChatServer{
		supportChats: make(map[string]*supportChat),
		inputCh: make(chan InputMessage),
		joinCh: make(chan JoinStatus),
		leaveCh: make(chan LeaveStatus),
		upd: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	go cS.handleMessages()

	return cS
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
	
	logrus.Error(joinMessage)
	if joinMessage.ChatID != "" {
		if cs.supportChats[joinMessage.ChatID] == nil {
			return ws, joinMessage, errors.New("chat not found")
		}

		return ws, joinMessage, nil
	}

	chatID := uuid.New().String()
	cs.supportChats[chatID] = &supportChat{}

	joinMessage.ChatID = chatID

	logrus.Error(joinMessage)

	return ws, joinMessage, nil
}

func (cs *ChatServer) SearchChat(w http.ResponseWriter, r *http.Request) (*websocket.Conn, *JoinStatus, error) {
	ws, err := cs.upd.Upgrade(w, r, nil)

	if err != nil {
		return nil, nil, err
	}

	joinMessage := &JoinStatus{}

	if err := ws.ReadJSON(&joinMessage); err != nil {
		return ws, nil, err
	}

	if joinMessage.ChatID != "" {
		if cs.supportChats[joinMessage.ChatID] == nil {
			return ws, joinMessage, errors.New("chat not found")
		}

		return ws, joinMessage, nil
	}

	return ws, joinMessage, errors.New("chat id not found")
}

func (cs *ChatServer) JoinUser(conn *websocket.Conn, jM *JoinStatus) error {
	if chat := cs.supportChats[jM.ChatID]; chat == nil {
		return errors.New("chat not found")
	}

	cs.supportChats[jM.ChatID].User = &chatMember{
		FullName: jM.FullName,
		ws:       conn,
	}

	if err := cs.supportChats[jM.ChatID].NotifyMembers(&JoinStatus{
		ChatID:   jM.ChatID,
		FullName: jM.FullName,
		Joined:   true,
	}); err != nil {
		return err
	}

	return nil
}

func (cs *ChatServer) LeaveUser(chatID string) error {
	var chat supportChat

	if chat := cs.supportChats[chatID]; chat == nil {
		return errors.New("chat not found")
	}

	if err := chat.User.CloseConn(); err != nil {
		return err
	}

	if chat.Dead() {
		delete(cs.supportChats, chatID)
	}

	return nil
}

func (cs *ChatServer) JoinSupport(conn *websocket.Conn, jM *JoinStatus) error {
	if chat := cs.supportChats[jM.ChatID]; chat == nil {
		return errors.New("chat not found")
	}

	if cs.supportChats[jM.ChatID].Support != nil {
		return errors.New("support already joined")
	}

	cs.supportChats[jM.ChatID].Support = &chatMember{
		FullName: jM.FullName,
		ws:       conn,
	}

	if err := cs.supportChats[jM.ChatID].NotifyMembers(&JoinStatus{
		ChatID:   jM.ChatID,
		FullName: jM.FullName,
		Joined:   true,
	}); err != nil {
		return err
	}

	return nil
}

func (cs *ChatServer) LeaveSupport(chatID string) error {
	var chat supportChat

	if chat := cs.supportChats[chatID]; chat == nil {
		return errors.New("chat not found")
	}

	if err := chat.Support.CloseConn(); err != nil {
		return err
	}

	if chat.Dead() {
		delete(cs.supportChats, chatID)
	}

	return nil
}

func (cs *ChatServer) GetSupportChats() map[string]*supportChat {
	return cs.supportChats
}
