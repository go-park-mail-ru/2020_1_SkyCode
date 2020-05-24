package notificationsWS

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type NotificationClient struct {
	UserID uint64
	server *NotificationServer
	noteCh chan *models.Notification
	ws     []*websocket.Conn
}

type StatusMessage struct {
	UserID uint64
	Status uint64
}

type NotificationServer struct {
	noteChats map[uint64]*NotificationClient
	client    chan *NotificationClient
	upd       *websocket.Upgrader
}

func NewNotificationServer() *NotificationServer {
	logrus.SetLevel(logrus.DebugLevel)
	ns := &NotificationServer{
		noteChats: make(map[uint64]*NotificationClient),
		client:    make(chan *NotificationClient),
		upd: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	go ns.run()

	return ns
}

func (ns *NotificationServer) run() {
	logrus.Debug("Server started")
	for c := range ns.client {
		if curClient := ns.noteChats[c.UserID]; curClient != nil {
			delete(ns.noteChats, c.UserID)
			close(c.noteCh)
		}
	}
}

func (ns *NotificationServer) CreateClient(w http.ResponseWriter, r *http.Request,
	userID uint64) (*websocket.Conn, error) {
	ws, err := ns.upd.Upgrade(w, r, nil)

	if ns.noteChats[userID] != nil {
		logrus.Debug("Adding ws")
		if err := ns.noteChats[userID].AddConnection(ws); err != nil {
			return nil, err
		}

		return ws, nil
	}

	if err != nil {
		return nil, err
	}

	nc := &NotificationClient{
		UserID: userID,
		server: ns,
		noteCh: make(chan *models.Notification),
		ws:     []*websocket.Conn{},
	}

	if err := nc.AddConnection(ws); err != nil {
		return nil, err
	}

	ns.noteChats[userID] = nc

	go ns.noteChats[userID].handleMessage()

	return ws, nil
}

func (ns *NotificationServer) SendNotification(note *models.Notification) {
	if client := ns.noteChats[note.UserID]; client != nil {
		client.noteCh <- note
	}
}

func (nc *NotificationClient) handleMessage() {
	logrus.Debug("Starting handle client")
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		logrus.Debug("Close Client")
		nc.server.client <- nc
	}()
	for {
		select {
		case noteMes := <-nc.noteCh:
			nc.WSSendNotification(noteMes)
		case <-ticker.C:
			i := 0
			for _, ws := range nc.ws {
				if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
					ws.Close()
					logrus.Debug("Closing ws")
					continue
				}
				nc.ws[i] = ws
				i++
			}
			nc.ws = nc.ws[:i]
			if len(nc.ws) == 0 {
				return
			}
		}
	}
}

func (nc *NotificationClient) WSSendNotification(note *models.Notification) {
	for _, ws := range nc.ws {
		if err := ws.WriteJSON(note); err != nil {
			logrus.Error(err)
		}
	}
}

func (nc *NotificationClient) AddConnection(ws *websocket.Conn) error {
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

	nc.ws = append(nc.ws, ws)

	return nil
}
