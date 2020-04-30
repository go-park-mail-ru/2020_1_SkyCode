package supportChat

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func echoUser(w http.ResponseWriter, r *http.Request) {
	serv := NewChatServer()

	ws, joinMsg, err := serv.CreateChat(w, r)
	if err != nil {
		fmt.Errorf("error creating chat %s", err)
		return
	}

	joinMsg.UserID = uint64(1)
	joinMsg.UserName = "test"

	err = serv.JoinUser(ws, &JoinStatus{
		UserID:   joinMsg.UserID,
		UserName: joinMsg.UserName,
		ChatID:   joinMsg.ChatID,
	})

	if err != nil {
		fmt.Errorf("error joining %s", err)
		return
	}

	msg := &InputMessage{}
	err = ws.ReadJSON(msg)
	if err != nil {
		fmt.Errorf("error joining %s", err)
		return
	}

	serv.WriteInputCh(*msg)

	lMsg := &LeaveStatus{}
	err = ws.ReadJSON(lMsg)
	if err != nil {
		fmt.Errorf("error joining %s", err)
		return
	}

	err = serv.LeaveUser(lMsg.ChatID)
	if err != nil {
		fmt.Errorf("error leaving %s", err)
		return
	}

	defer ws.Close()
}

func TestChatServer_JoinUser(t *testing.T) {
	logrus.SetLevel(logrus.PanicLevel)
	s := httptest.NewServer(http.HandlerFunc(echoUser))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")

	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}

	jMsg := &JoinStatus{
		ChatID: "testid",
	}

	msg := &InputMessage{
		ChatID:   "testid",
		Message:  "qq",
		UserID:   1,
		UserName: "test",
	}

	lMsg := &LeaveStatus{
		ChatID:   "testid",
		UserID:   1,
		UserName: "test",
		Leaved:   false,
	}

	err = ws.WriteJSON(jMsg)
	require.NoError(t, err)

	err = ws.WriteJSON(msg)
	require.NoError(t, err)

	err = ws.WriteJSON(lMsg)
	require.NoError(t, err)

	defer ws.Close()
}

func echoSupport(w http.ResponseWriter, r *http.Request) {
	serv := NewChatServer()

	ws, joinMsg, err := serv.CreateChat(w, r)
	if err != nil {
		fmt.Errorf("error creating chat %s", err)
		return
	}

	joinMsg.UserID = uint64(1)
	joinMsg.UserName = "test"

	err = serv.JoinSupport(ws, &JoinStatus{
		UserID:   joinMsg.UserID,
		UserName: joinMsg.UserName,
		ChatID:   joinMsg.ChatID,
	})

	if err != nil {
		fmt.Errorf("error joining %s", err)
		return
	}

	msg := &LeaveStatus{}
	err = ws.ReadJSON(msg)
	if err != nil {
		fmt.Errorf("error leaving %s", err)
		return
	}

	serv.WriteLeaveCh(*msg)

	err = serv.LeaveSupport(msg.ChatID)
	if err != nil {
		fmt.Errorf("error leaving %s", err)
		return
	}

	defer ws.Close()
}

func TestChatServer_JoinSupport(t *testing.T) {
	logrus.SetLevel(logrus.PanicLevel)
	s := httptest.NewServer(http.HandlerFunc(echoSupport))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")

	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}

	jMsg := &JoinStatus{
		ChatID: "testid",
	}

	lMsg := &LeaveStatus{
		ChatID:   "testid",
		UserID:   1,
		UserName: "test",
		Leaved:   false,
	}

	err = ws.WriteJSON(jMsg)
	require.NoError(t, err)

	err = ws.WriteJSON(lMsg)
	require.NoError(t, err)

	defer ws.Close()
}
