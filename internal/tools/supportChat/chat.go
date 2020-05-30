package supportChat

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	pongWait   = 30 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type chatMember struct {
	UserID   uint64
	UserName string
	wsList   []*websocket.Conn
}

func (cm *chatMember) CheckWS() bool {
	i := 0
	for _, ws := range cm.wsList {
		if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
			ws.Close()
			logrus.Debug("Closing ws")
			continue
		}
		cm.wsList[i] = ws
		i++
	}
	cm.wsList = cm.wsList[:i]
	if len(cm.wsList) == 0 {
		return false
	}

	return true
}

func (cm *chatMember) CheckConnection() bool {
	if len(cm.wsList) == 0 {
		return false
	}

	return true
}

func (cm *chatMember) AddConnection(ws *websocket.Conn) error {
	if err := ws.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		logrus.Error(err)
		return err
	}

	ws.SetPongHandler(func(string) error {
		if err := ws.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			return err
		}
		return nil
	})

	cm.wsList = append(cm.wsList, ws)

	return nil
}

type SupportChat struct {
	inputCh chan *InputMessage
	User    *chatMember
	Support *chatMember
}

func (sc *SupportChat) NotifyMembers(message interface{}) {
	if user := sc.User; user != nil {
		for _, ws := range user.wsList {
			if err := ws.WriteJSON(message); err != nil {
				logrus.Error(err)
			}

		}
	}

	if support := sc.Support; support != nil {
		for _, ws := range support.wsList {
			if err := ws.WriteJSON(message); err != nil {
				logrus.Error(err)
			}
		}
	}
}

func (sc *SupportChat) handleMessages() {
	logrus.Debug("Starting handle chat")
	ticker := time.NewTicker(pingPeriod)

	for {
		select {
		case <-ticker.C:
			sc.User.CheckWS()
			sc.Support.CheckWS()
		}
	}
}

func (cs *SupportChat) JoinUser(conn *websocket.Conn, user *models.User) error {
	if cs.User.UserID != user.ID {
		return tools.PermissionError
	}

	if err := cs.User.AddConnection(conn); err != nil {
		return err
	}

	cs.NotifyMembers(&JoinStatus{
		ChatID:   cs.User.UserID,
		UserID:   user.ID,
		UserName: user.FirstName,
		Joined:   true,
	})

	return nil
}

func (cs *SupportChat) JoinSupport(conn *websocket.Conn, user *models.User) error {
	if cs.Support.UserID != user.ID {
		if len(cs.Support.wsList) != 0 {
			return tools.NewSupportJoinError
		}
		cs.Support.UserID = user.ID
		cs.Support.UserName = user.FirstName
	}

	if err := cs.Support.AddConnection(conn); err != nil {
		return err
	}

	cs.NotifyMembers(&JoinStatus{
		ChatID:   cs.User.UserID,
		UserID:   user.ID,
		UserName: user.FirstName,
		Joined:   true,
	})

	return nil
}

type InputMessage struct {
	ChatID   uint64 `json:"chat_id"`
	UserID   uint64 `json:"user_id"`
	UserName string `json:"user_name"`
	Message  string `json:"message"`
}

type JoinStatus struct {
	ChatID   uint64 `json:"chat_id"`
	UserID   uint64 `json:"user_id"`
	UserName string `json:"user_name"`
	Joined   bool   `json:"joined"`
}

type ChatServer struct {
	supportChats map[uint64]*SupportChat
	chatControl  chan uint64
	upd          *websocket.Upgrader
}

func NewChatServer() *ChatServer {
	cS := &ChatServer{
		supportChats: make(map[uint64]*SupportChat),
		chatControl:  make(chan uint64),
		upd: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	return cS
}

func (cs *ChatServer) GetChat(chatID uint64) *SupportChat {
	return cs.supportChats[chatID]
}

func (cs *ChatServer) JoinUserToChat(user *models.User, w http.ResponseWriter, r *http.Request) (
	*websocket.Conn, *SupportChat, error) {
	conn, err := cs.upd.Upgrade(w, r, nil)
	if err != nil {
		return nil, nil, err
	}
	if chat := cs.supportChats[user.ID]; chat != nil {
		if err := chat.JoinUser(conn, user); err != nil {
			return nil, nil, err
		}
		return conn, chat, nil
	}

	chat := &SupportChat{
		inputCh: make(chan *InputMessage),
		User: &chatMember{
			UserID:   user.ID,
			UserName: user.FirstName,
		},
		Support: &chatMember{},
	}
	cs.supportChats[user.ID] = chat
	if err := chat.JoinUser(conn, user); err != nil {
		return nil, nil, err
	}

	go chat.handleMessages()

	return conn, chat, nil
}

func (cs *ChatServer) JoinSupportToChat(chatID uint64, sup *models.User, w http.ResponseWriter,
	r *http.Request) (*websocket.Conn, *SupportChat, error) {
	conn, err := cs.upd.Upgrade(w, r, nil)
	if err != nil {
		return nil, nil, err
	}
	if chat := cs.supportChats[chatID]; chat != nil {
		if err := chat.JoinSupport(conn, sup); err != nil {
			return nil, nil, err
		}
		return conn, chat, nil
	}

	return nil, nil, tools.ChatNotFound
}

func (cs *ChatServer) GetSupportChats() []*models.Chat {
	chats := []*models.Chat{}

	logrus.Debug(cs.supportChats)
	if len(cs.supportChats) == 0 {
		return chats
	}

	for ind, val := range cs.supportChats {
		logrus.Error(ind, val)

		chat := &models.Chat{
			ChatID:        val.User.UserID,
			UserName:      val.User.UserName,
			UserID:        val.User.UserID,
			UserConnected: val.User.CheckConnection(),
			SupConnected:  val.Support.CheckConnection(),
		}

		chats = append(chats, chat)
	}

	return chats
}

func (cs *ChatServer) DeleteChat(chatID uint64) error {
	if chat := cs.supportChats[chatID]; chat != nil {
		if chat.Support.CheckConnection() || chat.User.CheckConnection() {
			return tools.ChatInUse
		}

		close(chat.inputCh)
		delete(cs.supportChats, chatID)

		return nil
	}

	return tools.ChatNotFound
}
