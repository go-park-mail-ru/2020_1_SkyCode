package delivery

import (
	"encoding/json"
	"errors"
	_middleware "github.com/2020_1_Skycode/internal/middlewares"
	"github.com/2020_1_Skycode/internal/models"
	mock_sessions "github.com/2020_1_Skycode/internal/sessions/mocks"
	"github.com/2020_1_Skycode/internal/tools"
	_csrfManager "github.com/2020_1_Skycode/internal/tools/CSRFManager"
	_rValidator "github.com/2020_1_Skycode/internal/tools/requestValidator"
	mock_users "github.com/2020_1_Skycode/internal/users/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSessionHandler_LogOut(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	sessRes := &models.Session{ID: 1, UserId: 1}
	userRes := &models.User{Role: "Admin"}

	resResp := &tools.Message{Message: "success"}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(sessRes.UserId).Return(userRes, nil)
	mockSessUC.EXPECT().DeleteSession(sessRes.ID).Return(nil)

	g := gin.New()

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewSessionHandler(privateGroup, publicGroup, mockSessUC, mockUserUC, reqValidator, csrfManager, mwareC)

	target := "/api/v1/logout"
	req, err := http.NewRequest("POST", target, nil)
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{
		Name:  "SkyDelivery",
		Value: "1234",
	})
	w := httptest.NewRecorder()

	g.ServeHTTP(w, req)

	cookie := w.Result().Cookies()
	if len(cookie) == 0 {
		t.Error("Cookie is not seted")
		return
	}

	if w.Code != http.StatusOK {
		t.Error("Status is not ok")
		return
	}

	var result *tools.Message
	_ = json.NewDecoder(w.Result().Body).Decode(&result)

	require.EqualValues(t, resResp, result)
}

func TestSessionHandler_SignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	type signInRequest struct {
		Phone    string `json:"phone" binding:"required" validate:"min=11,max=15"`
		Password string `json:"password" binding:"required" validate:"passwd"`
	}

	userResp := &models.User{
		ID:       1,
		Password: "123456789",
		Phone:    "89765433221",
	}

	signInReq := &signInRequest{
		Phone:    userResp.Phone,
		Password: userResp.Password,
	}

	sigInJson, err := json.Marshal(signInReq)
	require.NoError(t, err)

	resResp := &tools.UserMessage{User: &models.User{
		ID:    userResp.ID,
		Phone: userResp.Phone,
	}}

	mockSessUC.EXPECT().GetSession(gomock.Any()).Return(nil, errors.New("bad"))
	mockUserUC.EXPECT().GetUserByPhone(signInReq.Phone).Return(userResp, nil)
	mockSessUC.EXPECT().StoreSession(gomock.Any()).Return(nil)

	g := gin.New()

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewSessionHandler(privateGroup, publicGroup, mockSessUC, mockUserUC, reqValidator, csrfManager, mwareC)

	target := "/api/v1/signin"
	req, err := http.NewRequest("POST", target, strings.NewReader(string(sigInJson)))
	req.Header.Set("Content-Type", "application/json")
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{
		Name:  "SkyDelivery",
		Value: "1234",
	})
	w := httptest.NewRecorder()

	g.ServeHTTP(w, req)

	cookie := w.Result().Cookies()
	if len(cookie) == 0 {
		t.Error("Cookie is not seted")
		return
	}

	if w.Code != http.StatusOK {
		t.Error("Status is not ok")
		return
	}

	var result *tools.UserMessage
	_ = json.NewDecoder(w.Result().Body).Decode(&result)

	require.EqualValues(t, resResp, result)
}
