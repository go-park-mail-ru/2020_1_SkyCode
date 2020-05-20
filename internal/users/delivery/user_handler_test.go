package delivery

import (
	"bytes"
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
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUserHandler_SignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	type signUpRequest struct {
		FirstName string `json:"firstName,omitempty" binding:"required" validate:"min=2"`
		LastName  string `json:"lastName,omitempty" binding:"required" validate:"min=2"`
		Phone     string `json:"phone,omitempty" binding:"required" validate:"min=11,max=15"`
		Password  string `json:"password,omitempty" binding:"required" validate:"passwd"`
	}

	reqSignUp := signUpRequest{
		FirstName: "AD",
		LastName:  "BS",
		Phone:     "89765433221",
		Password:  "1234567890",
	}

	userReq := &models.User{
		Password:  reqSignUp.Password,
		FirstName: reqSignUp.FirstName,
		LastName:  reqSignUp.LastName,
		Phone:     reqSignUp.Phone,
	}

	reqJson, err := json.Marshal(reqSignUp)
	require.NoError(t, err)

	expectResult := &tools.Message{Message: "User has been registered"}

	mockUserUC.EXPECT().GetUserByPhone(reqSignUp.Phone).Return(nil, errors.New("nothing"))
	mockUserUC.EXPECT().CreateUser(userReq).Return(nil)
	mockSessUC.EXPECT().StoreSession(gomock.Any()).Return(nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewUserHandler(privateGroup, publicGroup, mockUserUC, mockSessUC, reqValidator, csrfManager, mwareC)

	target := "/api/v1/signup"
	req, err := http.NewRequest("POST", target, strings.NewReader(string(reqJson)))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
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

	require.EqualValues(t, expectResult, result)
}

func TestUserHandler_ChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	reqChangePassword := changePasswordRequest{
		NewPassword: "1234567890",
	}

	userReq := &models.User{
		ID:       1,
		Password: reqChangePassword.NewPassword,
	}

	reqJson, err := json.Marshal(reqChangePassword)
	require.NoError(t, err)

	expectResult := &tools.Message{Message: "success"}

	sessRes := &models.Session{ID: 1, UserId: userReq.ID}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(userReq.ID).Return(userReq, nil)
	mockUserUC.EXPECT().UpdatePassword(userReq.ID, userReq.Password).Return(nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewUserHandler(privateGroup, publicGroup, mockUserUC, mockSessUC, reqValidator, csrfManager, mwareC)

	target := "/api/v1/profile/password"
	req, err := http.NewRequest("PUT", target, strings.NewReader(string(reqJson)))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
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

	var result *tools.Message
	_ = json.NewDecoder(w.Result().Body).Decode(&result)

	require.EqualValues(t, expectResult, result)
}

func TestUserHandler_ChangePhoneNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	reqChangePhone := changePhoneNumberRequest{
		NewPhone: "89765433221",
	}

	userReq := &models.User{
		ID:    1,
		Phone: reqChangePhone.NewPhone,
	}

	reqJson, err := json.Marshal(reqChangePhone)
	require.NoError(t, err)

	expectResult := &tools.Message{Message: "success"}

	sessRes := &models.Session{ID: 1, UserId: userReq.ID}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(userReq.ID).Return(userReq, nil)
	mockUserUC.EXPECT().UpdatePhoneNumber(userReq.ID, userReq.Phone).Return(nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewUserHandler(privateGroup, publicGroup, mockUserUC, mockSessUC, reqValidator, csrfManager, mwareC)

	target := "/api/v1/profile/phone"
	req, err := http.NewRequest("PUT", target, strings.NewReader(string(reqJson)))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
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

	var result *tools.Message
	_ = json.NewDecoder(w.Result().Body).Decode(&result)

	require.EqualValues(t, expectResult, result)
}

func TestUserHandler_GetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	userReq := &models.User{
		ID:        1,
		Phone:     "89765432211",
		FirstName: "AB",
		LastName:  "FG",
	}

	sessRes := &models.Session{ID: 1, UserId: userReq.ID}

	expectUser := &tools.UserMessage{User: userReq}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(userReq.ID).Return(userReq, nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewUserHandler(privateGroup, publicGroup, mockUserUC, mockSessUC, reqValidator, csrfManager, mwareC)

	target := "/api/v1/profile"
	req, err := http.NewRequest("GET", target, nil)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
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

	var result *tools.UserMessage
	_ = json.NewDecoder(w.Result().Body).Decode(&result)

	require.EqualValues(t, expectUser, result)
}

func TestUserHandler_EditBio(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	reqChangeBio := editBioRequest{
		FirstName: "ASG",
		LastName:  "ASJS",
		Email:     "newemail@m.ru",
	}

	userReq := &models.User{
		ID:        1,
		FirstName: reqChangeBio.FirstName,
		LastName:  reqChangeBio.LastName,
		Email:     reqChangeBio.Email,
	}

	reqJson, err := json.Marshal(reqChangeBio)
	require.NoError(t, err)

	expectResult := &tools.Message{Message: "success"}

	sessRes := &models.Session{ID: 1, UserId: userReq.ID}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(userReq.ID).Return(userReq, nil)
	mockUserUC.EXPECT().UpdateBio(userReq).Return(nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewUserHandler(privateGroup, publicGroup, mockUserUC, mockSessUC, reqValidator, csrfManager, mwareC)

	target := "/api/v1/profile/bio"
	req, err := http.NewRequest("PUT", target, strings.NewReader(string(reqJson)))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
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

	var result *tools.Message
	_ = json.NewDecoder(w.Result().Body).Decode(&result)

	require.EqualValues(t, expectResult, result)
}

func TestUserHandler_EditAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	if err := os.MkdirAll(tools.RestaurantImagesPath, 0777); err != nil {
		t.Errorf("Error on image work: %s", err)
		return
	}

	imageName := "1234.jpg"

	f, err := os.Create(filepath.Join(tools.ProductImagesPath, filepath.Base(imageName)))
	require.NoError(t, err)
	f.Close()

	userReq := &models.User{
		ID:     1,
		Avatar: imageName,
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("avatar", "testfile")

	part.Write([]byte("SOME FILE CONTENT"))

	writer.Close()

	expectResult := &tools.Message{Message: "success"}

	sessRes := &models.Session{ID: 1, UserId: userReq.ID}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(userReq.ID).Return(userReq, nil)
	mockUserUC.EXPECT().UpdateAvatar(userReq.ID, gomock.Any()).Return(nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewUserHandler(privateGroup, publicGroup, mockUserUC, mockSessUC, reqValidator, csrfManager, mwareC)

	target := "/api/v1/profile/avatar"
	req, err := http.NewRequest("PUT", target, body)
	require.NoError(t, err)

	req.Header.Set("Content-Type", writer.FormDataContentType())
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

	var result *tools.Message
	_ = json.NewDecoder(w.Result().Body).Decode(&result)

	require.EqualValues(t, expectResult, result)

	if err = os.RemoveAll(tools.RestaurantImagesPath); err != nil {
		t.Errorf("Error on image work: %s", err)
		return
	}
}
