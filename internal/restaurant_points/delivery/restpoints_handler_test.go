package delivery

import (
	"encoding/json"
	_middleware "github.com/2020_1_Skycode/internal/middlewares"
	"github.com/2020_1_Skycode/internal/models"
	mock_restpoints "github.com/2020_1_Skycode/internal/restaurant_points/mocks"
	mock_sessions "github.com/2020_1_Skycode/internal/sessions/mocks"
	"github.com/2020_1_Skycode/internal/tools"
	_csrfManager "github.com/2020_1_Skycode/internal/tools/CSRFManager"
	mock_users "github.com/2020_1_Skycode/internal/users/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRestPointsHandler_GetAllPoints(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRestPointsUCase := mock_restpoints.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	userID := uint64(1)

	sessRes := &models.Session{ID: 1, UserId: userID}
	userRes := &models.User{ID: sessRes.UserId, Role: "Admin"}

	expectRP := []*models.RestaurantPoint{
		{
			ID:      1,
			Address: "Pushkina dom Kolotushkina",
			MapPoint: &models.GeoPos{
				Longitude: 55.753227,
				Latitude:  37.619030,
			},
			ServiceRadius: 5,
			RestID:        1,
		},
	}

	expectRes := tools.Body{"points": expectRP}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(userID).Return(userRes, nil)
	mockRestPointsUCase.EXPECT().GetAllPoints().Return(expectRP, nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")

	_ = NewRestPointsHandler(privateGroup, publicGroup, mockRestPointsUCase, mwareC)

	target := "/api/v1/points"
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
	expect, err := json.Marshal(expectRes)
	require.NoError(t, err)

	require.EqualValues(t, expect, result)
}

func TestRestPointsHandler_GetPoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRestPointsUCase := mock_restpoints.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	userID := uint64(1)

	sessRes := &models.Session{ID: 1, UserId: userID}
	userRes := &models.User{ID: sessRes.UserId, Role: "Admin"}

	rpID := uint64(1)

	expectRP := &models.RestaurantPoint{
		ID:      rpID,
		Address: "Pushkina dom Kolotushkina",
		MapPoint: &models.GeoPos{
			Longitude: 55.753227,
			Latitude:  37.619030,
		},
		ServiceRadius: 5,
		RestID:        1,
	}

	expectRes := tools.Body{"point": expectRP}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(userID).Return(userRes, nil)
	mockRestPointsUCase.EXPECT().GetPoint(rpID).Return(expectRP, nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")

	_ = NewRestPointsHandler(privateGroup, publicGroup, mockRestPointsUCase, mwareC)

	target := "/api/v1/points/1"
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
	expect, err := json.Marshal(expectRes)
	require.NoError(t, err)

	require.EqualValues(t, expect, result)
}

func TestRestPointsHandler_DeletePoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRestPointsUCase := mock_restpoints.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	userID := uint64(1)

	sessRes := &models.Session{ID: 1, UserId: userID}
	userRes := &models.User{ID: sessRes.UserId, Role: "Admin"}

	rpID := uint64(1)

	expectRes := tools.Message{Message: "Deleted"}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(userID).Return(userRes, nil)
	mockRestPointsUCase.EXPECT().Delete(rpID).Return(nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")

	_ = NewRestPointsHandler(privateGroup, publicGroup, mockRestPointsUCase, mwareC)

	target := "/api/v1/points/1"
	req, err := http.NewRequest("DELETE", target, nil)
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
	expect, err := json.Marshal(expectRes)
	require.NoError(t, err)

	require.EqualValues(t, expect, result)
}
