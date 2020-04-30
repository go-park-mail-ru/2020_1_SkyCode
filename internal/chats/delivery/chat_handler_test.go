package delivery

import (
	"encoding/json"
	"errors"
	mock_chats "github.com/2020_1_Skycode/internal/chats/mocks"
	_middleware "github.com/2020_1_Skycode/internal/middlewares"
	"github.com/2020_1_Skycode/internal/models"
	mock_sessions "github.com/2020_1_Skycode/internal/sessions/mocks"
	_csrfManager "github.com/2020_1_Skycode/internal/tools/CSRFManager"
	"github.com/2020_1_Skycode/internal/tools/supportChat"
	mock_users "github.com/2020_1_Skycode/internal/users/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestChatHandler_GetChatMessages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)
	mockChatUC := mock_chats.NewMockUseCase(ctrl)

	userID := uint64(1)
	chatID := "testid"

	sessRes := &models.Session{ID: 1, UserId: userID}
	userRes := &models.User{ID: sessRes.UserId, Role: "Support"}

	returnChat := &supportChat.SupportChat{}

	returnMessgs := []*models.ChatMessage{
		{
			UserID:   userID,
			UserName: "test",
			ChatID:   chatID,
			Message:  "Che s den'gami?",
			Created:  time.Now(),
		},
	}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(userID).Return(userRes, nil)
	mockChatUC.EXPECT().GetChat(chatID).Return(returnChat)
	mockChatUC.EXPECT().GetChatMessages(chatID).Return(returnMessgs, nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")

	_ = NewChatsHandler(privateGroup, publicGroup, mockChatUC, mwareC)

	target := "/api/v1/chats/" + chatID + "/details"
	req, err := http.NewRequest("GET", target, nil)
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{
		Name:  "SkyDelivery",
		Value: "1234",
	})

	w := httptest.NewRecorder()

	g.ServeHTTP(w, req)

	result, err := ioutil.ReadAll(w.Result().Body)
	require.NoError(t, err)
	expect, err := json.Marshal(returnMessgs)
	require.NoError(t, err)

	require.EqualValues(t, expect, result)
}

func TestChatHandler_StartUserChat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)
	mockChatUC := mock_chats.NewMockUseCase(ctrl)

	userID := uint64(1)
	chatID := "testid"

	sessRes := &models.Session{ID: 1, UserId: userID}
	userRes := &models.User{ID: sessRes.UserId, FirstName: "test", Role: "User"}

	returnStatus := &supportChat.JoinStatus{
		ChatID: chatID,
		Joined: false,
	}

	conn := &websocket.Conn{}

	mess := supportChat.InputMessage{
		ChatID:   chatID,
		UserID:   userID,
		UserName: userRes.FirstName,
		Message:  "asd",
	}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(userID).Return(userRes, nil)
	mockChatUC.EXPECT().StartChat(gomock.Any(), gomock.Any()).Return(conn, returnStatus, nil)
	mockChatUC.EXPECT().JoinUserToChat(conn, userRes.ID, userRes.FirstName, chatID).Return(nil)
	mockChatUC.EXPECT().ReadMessageFromUSer(conn).Return(mess, nil)
	mockChatUC.EXPECT().StoreMessage(gomock.Any()).Return(nil)
	mockChatUC.EXPECT().WriteFromUserMessage(gomock.Any()).Return()
	mockChatUC.EXPECT().ReadMessageFromUSer(conn).Return(supportChat.InputMessage{}, errors.New("close"))

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")

	_ = NewChatsHandler(privateGroup, publicGroup, mockChatUC, mwareC)

	target := "/api/v1/chat"
	req, err := http.NewRequest("GET", target, nil)
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{
		Name:  "SkyDelivery",
		Value: "1234",
	})

	w := httptest.NewRecorder()

	g.ServeHTTP(w, req)
}

func TestChatHandler_JoinSupport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)
	mockChatUC := mock_chats.NewMockUseCase(ctrl)

	userID := uint64(1)
	chatID := "testid"

	sessRes := &models.Session{ID: 1, UserId: userID}
	userRes := &models.User{ID: sessRes.UserId, FirstName: "test", Role: "Support"}

	returnStatus := &supportChat.JoinStatus{
		ChatID: chatID,
		Joined: false,
	}

	conn := &websocket.Conn{}

	mess := supportChat.InputMessage{
		ChatID:   chatID,
		UserID:   userID,
		UserName: userRes.FirstName,
		Message:  "asd",
	}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(userID).Return(userRes, nil)
	mockChatUC.EXPECT().FindChat(gomock.Any(), gomock.Any(), chatID).Return(conn, returnStatus, nil)
	mockChatUC.EXPECT().JoinSupportToChat(conn, userRes.ID, userRes.FirstName, chatID).Return(nil)
	mockChatUC.EXPECT().ReadMessageFromUSer(conn).Return(mess, nil)
	mockChatUC.EXPECT().StoreMessage(gomock.Any()).Return(nil)
	mockChatUC.EXPECT().WriteFromUserMessage(gomock.Any()).Return()
	mockChatUC.EXPECT().ReadMessageFromUSer(conn).Return(supportChat.InputMessage{}, errors.New("close"))
	mockChatUC.EXPECT().LeaveSupportChat(chatID).Return(nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")

	_ = NewChatsHandler(privateGroup, publicGroup, mockChatUC, mwareC)

	target := "/api/v1/chats/" + chatID + "/join"
	req, err := http.NewRequest("GET", target, nil)
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{
		Name:  "SkyDelivery",
		Value: "1234",
	})

	w := httptest.NewRecorder()

	g.ServeHTTP(w, req)
}

func TestChatHandler_GetSupChatList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)
	mockChatUC := mock_chats.NewMockUseCase(ctrl)

	userID := uint64(1)
	chatID := "testid"

	sessRes := &models.Session{ID: 1, UserId: userID}
	userRes := &models.User{ID: sessRes.UserId, Role: "Support"}

	returnChats := []*models.Chat{
		{
			UserID:   userID,
			UserName: "testuser",
			ChatID:   chatID,
		},
	}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(userID).Return(userRes, nil)
	mockChatUC.EXPECT().GetChats().Return(returnChats)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")

	_ = NewChatsHandler(privateGroup, publicGroup, mockChatUC, mwareC)

	target := "/api/v1/chats"
	req, err := http.NewRequest("GET", target, nil)
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{
		Name:  "SkyDelivery",
		Value: "1234",
	})

	w := httptest.NewRecorder()

	g.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("Status is not ok")
		return
	}

	result, err := ioutil.ReadAll(w.Result().Body)
	require.NoError(t, err)
	expect, err := json.Marshal(returnChats)
	require.NoError(t, err)

	require.EqualValues(t, expect, result)
}
