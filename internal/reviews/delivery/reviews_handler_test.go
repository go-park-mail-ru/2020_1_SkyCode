package delivery

import (
	"encoding/json"
	_middleware "github.com/2020_1_Skycode/internal/middlewares"
	"github.com/2020_1_Skycode/internal/models"
	mock_reviews "github.com/2020_1_Skycode/internal/reviews/mocks"
	mock_sessions "github.com/2020_1_Skycode/internal/sessions/mocks"
	"github.com/2020_1_Skycode/internal/tools"
	_csrfManager "github.com/2020_1_Skycode/internal/tools/CSRFManager"
	_rValidator "github.com/2020_1_Skycode/internal/tools/requestValidator"
	mock_users "github.com/2020_1_Skycode/internal/users/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestReviewsHandler_DeleteReview(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewsUCase := mock_reviews.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	userID := uint64(1)
	testID := uint64(1)

	expectRes := &tools.Message{Message: "Deleted"}

	sessRes := &models.Session{ID: 1, UserId: userID}
	userRes := &models.User{ID: sessRes.UserId, Role: "User"}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(userID).Return(userRes, nil)
	mockReviewsUCase.EXPECT().DeleteReview(testID, userRes).Return(nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)
	reqValidator := _rValidator.NewRequestValidator()

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")

	_ = NewReviewsHandler(privateGroup, publicGroup, mockReviewsUCase, reqValidator, mwareC)
	target := "/api/v1/reviews/" + strconv.Itoa(int(testID))
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

func TestReviewsHandler_UpdateReview(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewsUCase := mock_reviews.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	userID := uint64(1)
	expectRes := &tools.Message{Message: "Updated"}

	sessRes := &models.Session{ID: 1, UserId: userID}
	userRes := &models.User{ID: sessRes.UserId, Role: "User"}

	testReview := &models.Review{
		ID:     1,
		Text:   "Good",
		Author: &models.User{ID: userRes.ID},
		Rate:   5,
	}
	type reviewUpdateRequest struct {
		Text string  `json:"text, omitempty" binding:"required" validate:"min=2"`
		Rate float64 `json:"rate" binding:"required" validate:"min=0,max=5"`
	}

	testReq := &reviewUpdateRequest{
		Text: testReview.Text,
		Rate: testReview.Rate,
	}

	j, err := json.Marshal(testReq)
	require.NoError(t, err)

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(userID).Return(userRes, nil)
	mockReviewsUCase.EXPECT().UpdateReview(testReview, userRes).Return(nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)
	reqValidator := _rValidator.NewRequestValidator()

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")

	_ = NewReviewsHandler(privateGroup, publicGroup, mockReviewsUCase, reqValidator, mwareC)
	target := "/api/v1/reviews/" + strconv.Itoa(int(testReview.ID))
	req, err := http.NewRequest("PUT", target, strings.NewReader(string(j)))
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

	result, err := ioutil.ReadAll(w.Result().Body)
	require.NoError(t, err)
	expect, err := json.Marshal(expectRes)
	require.NoError(t, err)

	require.EqualValues(t, expect, result)
}

func TestReviewsHandler_GetUserReviews(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewsUCase := mock_reviews.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	userID := uint64(1)

	sessRes := &models.Session{ID: 1, UserId: userID}
	userRes := &models.User{ID: sessRes.UserId, Role: "User"}

	testReviews := []*models.Review{
		{
			ID:           1,
			RestID:       1,
			Text:         "Good",
			Author:       userRes,
			CreationDate: time.Now(),
			Rate:         5,
		},
	}

	expectRes := &tools.Body{"reviews": testReviews, "total": uint64(1)}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(userID).Return(userRes, nil)
	mockReviewsUCase.EXPECT().GetUserReviews(userRes.ID, uint64(1), uint64(1)).Return(testReviews, uint64(1), nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)
	reqValidator := _rValidator.NewRequestValidator()

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")

	_ = NewReviewsHandler(privateGroup, publicGroup, mockReviewsUCase, reqValidator, mwareC)
	target := "/api/v1/reviews"
	req, err := http.NewRequest("GET", target, nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "SkyDelivery",
		Value: "1234",
	})
	q := req.URL.Query()
	q.Add("page", "1")
	q.Add("count", "1")
	req.URL.RawQuery = q.Encode()

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

func TestReviewsHandler_GetReview(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewsUCase := mock_reviews.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	userID := uint64(1)

	sessRes := &models.Session{ID: 1, UserId: userID}
	userRes := &models.User{ID: sessRes.UserId, Role: "User"}

	testReview := &models.Review{
		ID:           1,
		RestID:       1,
		Text:         "Good",
		Author:       userRes,
		CreationDate: time.Now(),
		Rate:         5,
	}

	expectRes := &tools.Body{"review": testReview}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(userID).Return(userRes, nil)
	mockReviewsUCase.EXPECT().GetReview(testReview.ID).Return(testReview, nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)
	reqValidator := _rValidator.NewRequestValidator()

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")

	_ = NewReviewsHandler(privateGroup, publicGroup, mockReviewsUCase, reqValidator, mwareC)
	target := "/api/v1/reviews/" + strconv.Itoa(int(testReview.ID))
	req, err := http.NewRequest("GET", target, nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "SkyDelivery",
		Value: "1234",
	})
	q := req.URL.Query()
	q.Add("page", "1")
	q.Add("count", "1")
	req.URL.RawQuery = q.Encode()

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
