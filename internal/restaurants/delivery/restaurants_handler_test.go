package delivery

import (
	"bytes"
	"encoding/json"
	_middleware "github.com/2020_1_Skycode/internal/middlewares"
	"github.com/2020_1_Skycode/internal/models"
	mock_restaurants "github.com/2020_1_Skycode/internal/restaurants/mocks"
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
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestRestaurantHandler_GetRestaurants(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRestUC := mock_restaurants.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	resultRests := []*models.Restaurant{
		{ID: 1, Name: "test1", Rating: 4.2, Image: "./default.jpg"},
	}

	expectResp := &tools.Body{
		"restaurants": resultRests,
		"total":       uint64(1)}

	mockRestUC.EXPECT().GetRestaurants(uint64(1), uint64(1)).Return(resultRests, uint64(1), nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewRestaurantHandler(privateGroup, publicGroup, reqValidator, mockRestUC, mwareC)

	req := httptest.NewRequest("GET", "/api/v1/restaurants", nil)
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
	expect, err := json.Marshal(expectResp)
	require.NoError(t, err)

	require.EqualValues(t, expect, result)
}

func TestRestaurantHandler_GetRestaurantByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRestUC := mock_restaurants.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	resultRests := &models.Restaurant{
		ID:     1,
		Name:   "test1",
		Rating: 4.2,
		Image:  "./default.jpg",
		Products: []*models.Product{
			{ID: 1, Name: "testProd", Price: 2.50, Image: "./default_prod.jpg"},
		},
	}

	restID := uint64(1)

	mockRestUC.EXPECT().GetRestaurantByID(restID).Return(resultRests, nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewRestaurantHandler(privateGroup, publicGroup, reqValidator, mockRestUC, mwareC)

	target := "/api/v1/restaurants/" + strconv.Itoa(int(restID))
	req := httptest.NewRequest("GET", target, nil)
	w := httptest.NewRecorder()

	g.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("Status is not ok")
		return
	}

	var result *models.Restaurant
	_ = json.NewDecoder(w.Result().Body).Decode(&result)

	require.EqualValues(t, resultRests, result)
}

func TestRestaurantHandler_CreateRestaurant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRestUC := mock_restaurants.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	reqRests := &models.Restaurant{
		Name:        "test1",
		Description: "asdasdqweasdqwe",
	}

	if err := os.MkdirAll(tools.RestaurantImagesPath, 0777); err != nil {
		t.Errorf("Error on image work: %s", err)
		return
	}

	expectResult := &tools.Message{Message: "Restaurant has been created"}

	userID := uint64(1)

	sessRes := &models.Session{ID: 1, UserId: userID}
	userRes := &models.User{ID: userID, Role: "Admin"}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(userID).Return(userRes, nil)
	mockRestUC.EXPECT().CreateRestaurant(gomock.Any()).Return(nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.TraceLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	privateGroup.Use(mwareC.CSRFControl())
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewRestaurantHandler(privateGroup, publicGroup, reqValidator, mockRestUC, mwareC)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	err := writer.WriteField("Name", reqRests.Name)
	require.NoError(t, err)
	err = writer.WriteField("Description", reqRests.Description)
	require.NoError(t, err)
	part, _ := writer.CreateFormFile("image", "testfile")

	part.Write([]byte("SOME FILE CONTENT"))

	writer.Close()

	target := "/api/v1/restaurants"
	req, err := http.NewRequest("POST", target, body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	token, err := csrfManager.GenerateCSRF(userID, sessRes.Token)
	require.NoError(t, err)
	req.Header.Set("X-Csrf-Token", token)
	req.AddCookie(&http.Cookie{
		Name:  "SkyDelivery",
		Value: "1234",
	})

	w := httptest.NewRecorder()

	g.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status is not ok %v, %v", w.Code, w.Body)
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

func TestRestaurantHandler_CreateRestaurant_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRestUC := mock_restaurants.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	reqRests := &models.Restaurant{
		Name:        "t",
		Description: "asdasdqweasdqwe",
	}

	userID := uint64(1)

	sessRes := &models.Session{ID: 1, UserId: userID}
	userRes := &models.User{ID: userID, Role: "Admin"}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(userID).Return(userRes, nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	privateGroup.Use(mwareC.CSRFControl())
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewRestaurantHandler(privateGroup, publicGroup, reqValidator, mockRestUC, mwareC)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	err := writer.WriteField("Name", reqRests.Name)
	require.NoError(t, err)
	err = writer.WriteField("Description", reqRests.Description)
	require.NoError(t, err)
	part, _ := writer.CreateFormFile("image", "testfile")

	part.Write([]byte("SOME FILE CONTENT"))

	writer.Close()

	target := "/api/v1/restaurants"
	req, err := http.NewRequest("POST", target, body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	token, err := csrfManager.GenerateCSRF(userID, sessRes.Token)
	require.NoError(t, err)
	req.Header.Set("X-Csrf-Token", token)
	req.AddCookie(&http.Cookie{
		Name:  "SkyDelivery",
		Value: "1234",
	})

	w := httptest.NewRecorder()

	g.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Error("Status is not bad request")
		return
	}
}

func TestRestaurantHandler_DeleteRestaurant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRestUC := mock_restaurants.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	restID := uint64(1)

	resultRests := &models.Restaurant{
		ID:       1,
		Name:     "test1",
		Rating:   4.2,
		Image:    "",
		Products: nil,
	}

	mockRestUC.EXPECT().GetRestaurantByID(restID).Return(resultRests, nil)
	mockRestUC.EXPECT().Delete(restID).Return(nil)

	expectResult := &tools.Message{Message: "success"}

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewRestaurantHandler(privateGroup, publicGroup, reqValidator, mockRestUC, mwareC)

	target := "/api/v1/restaurants/" + strconv.Itoa(int(restID))
	req, err := http.NewRequest("DELETE", target, nil)
	require.NoError(t, err)
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

func TestRestaurantHandler_UpdateImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRestUC := mock_restaurants.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	restID := uint64(1)

	resultRests := &models.Restaurant{
		ID:       1,
		Name:     "test1",
		Rating:   4.2,
		Image:    "1234.jpg",
		Products: nil,
	}

	if err := os.MkdirAll(tools.RestaurantImagesPath, 0777); err != nil {
		t.Errorf("Error on image work: %s", err)
		return
	}
	f, err := os.Create(filepath.Join(tools.RestaurantImagesPath, filepath.Base(resultRests.Image)))
	require.NoError(t, err)
	f.Close()

	expectResult := &tools.Message{Message: "success"}

	mockRestUC.EXPECT().GetRestaurantByID(restID).Return(resultRests, nil)
	mockRestUC.EXPECT().UpdateImage(restID, gomock.Any()).Return(nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewRestaurantHandler(privateGroup, publicGroup, reqValidator, mockRestUC, mwareC)

	bodyUpdate := &bytes.Buffer{}
	writer := multipart.NewWriter(bodyUpdate)
	part, _ := writer.CreateFormFile("image", "testfile")

	part.Write([]byte("SOME FILE CONTENT"))

	writer.Close()

	target := "/api/v1/restaurants/" + strconv.Itoa(int(restID)) + "/image"
	req, err := http.NewRequest("PUT", target, bodyUpdate)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	g.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("Status is not ok")
		return
	}

	if err = os.RemoveAll(tools.RestaurantImagesPath); err != nil {
		t.Errorf("Error on image work: %s", err)
		return
	}

	var result *tools.Message
	_ = json.NewDecoder(w.Result().Body).Decode(&result)

	require.EqualValues(t, expectResult, result)
}

func TestRestaurantHandler_UpdateRestaurant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRestUC := mock_restaurants.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	reqRests := &models.Restaurant{
		ID:          uint64(1),
		Name:        "test1",
		Description: "asdasdqweasdqwe",
	}

	type restaurantRequest struct {
		Name        string `json:"name, omitempty" binding:"required" validate:"min=3"`
		Description string `json:"description, omitempty" binding:"required" validate:"min=10"`
	}

	restID := uint64(1)

	reqRest := &restaurantRequest{
		Name:        reqRests.Name,
		Description: reqRests.Description,
	}

	j, err := json.Marshal(reqRest)
	require.NoError(t, err)

	expectResult := &tools.Message{Message: "Restaurant has been updated"}

	mockRestUC.EXPECT().UpdateRestaurant(reqRests).Return(nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewRestaurantHandler(privateGroup, publicGroup, reqValidator, mockRestUC, mwareC)

	target := "/api/v1/restaurants/" + strconv.Itoa(int(restID)) + "/update"
	req, err := http.NewRequest("PUT", target, strings.NewReader(string(j)))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
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

func TestRestaurantHandler_AddPoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRestUC := mock_restaurants.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	expectResult := &tools.Message{Message: "Created"}

	userID := uint64(1)
	restID := uint64(1)

	type pointRequest struct {
		Address string  `json:"address" binding:"required"`
		Radius  float64 `json:"radius" binding:"required"`
	}

	r := pointRequest{
		Address: "Pushkina dom Kolotushkina",
		Radius:  5,
	}

	reqPoint := &models.RestaurantPoint{
		Address:       r.Address,
		ServiceRadius: r.Radius,
		RestID:        restID,
	}

	sessRes := &models.Session{ID: 1, UserId: userID}
	userRes := &models.User{ID: userID, Role: "Admin"}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(userID).Return(userRes, nil)
	mockRestUC.EXPECT().AddPoint(reqPoint).Return(nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewRestaurantHandler(privateGroup, publicGroup, reqValidator, mockRestUC, mwareC)

	j, err := json.Marshal(r)
	require.NoError(t, err)

	target := "/api/v1/restaurants/" + strconv.Itoa(int(restID)) + "/points"
	req, err := http.NewRequest("POST", target, strings.NewReader(string(j)))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "SkyDelivery",
		Value: "1234",
	})

	w := httptest.NewRecorder()

	g.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("Status is not ok ", w.Code)
		return
	}

	var result *tools.Message
	_ = json.NewDecoder(w.Result().Body).Decode(&result)

	require.EqualValues(t, expectResult, result)
}

func TestRestaurantHandler_AddReview(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRestUC := mock_restaurants.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	expectResult := &tools.Message{Message: "Created"}

	userID := uint64(1)
	restID := uint64(1)

	type reviewRequest struct {
		Text string   `json:"text, omitempty" binding:"required"`
		Rate *float64 `json:"rate" binding:"required" validate:"min=0,max=5"`
	}

	rate := float64(5)

	r := &reviewRequest{
		Text: "Good test",
		Rate: &rate,
	}

	sessRes := &models.Session{ID: 1, UserId: userID}
	userRes := &models.User{ID: userID, Role: "User"}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(userID).Return(userRes, nil)
	mockRestUC.EXPECT().AddReview(gomock.Any()).Return(nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewRestaurantHandler(privateGroup, publicGroup, reqValidator, mockRestUC, mwareC)

	j, err := json.Marshal(r)
	require.NoError(t, err)

	target := "/api/v1/restaurants/" + strconv.Itoa(int(restID)) + "/reviews"
	req, err := http.NewRequest("POST", target, strings.NewReader(string(j)))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "SkyDelivery",
		Value: "1234",
	})

	w := httptest.NewRecorder()

	g.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("Status is not ok ", w.Code)
		return
	}

	var result *tools.Message
	_ = json.NewDecoder(w.Result().Body).Decode(&result)

	require.EqualValues(t, expectResult, result)
}

func TestRestaurantHandler_AddReview_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRestUC := mock_restaurants.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	userID := uint64(1)
	restID := uint64(1)

	type reviewRequest struct {
		Text string   `json:"text, omitempty" binding:"required"`
		Rate *float64 `json:"rate" binding:"required" validate:"min=0,max=5"`
	}

	rate := float64(-1)

	r := &reviewRequest{
		Text: "Good test",
		Rate: &rate,
	}

	sessRes := &models.Session{ID: 1, UserId: userID}
	userRes := &models.User{ID: userID, Role: "User"}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(userID).Return(userRes, nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewRestaurantHandler(privateGroup, publicGroup, reqValidator, mockRestUC, mwareC)

	j, err := json.Marshal(r)
	require.NoError(t, err)

	target := "/api/v1/restaurants/" + strconv.Itoa(int(restID)) + "/reviews"
	req, err := http.NewRequest("POST", target, strings.NewReader(string(j)))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "SkyDelivery",
		Value: "1234",
	})

	w := httptest.NewRecorder()

	g.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Error("Status is not bad request ", w.Code)
		return
	}
}

func TestRestaurantHandler_GetPoints(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRestUC := mock_restaurants.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	restID := uint64(1)

	resultPoints := []*models.RestaurantPoint{
		{
			ID:      1,
			RestID:  1,
			Address: "Pushkina dom Kolotushkina",
			MapPoint: &models.GeoPos{
				Latitude:  33.5555,
				Longitude: 55.3333,
			},
			ServiceRadius: 5,
		},
	}

	expectResp := &tools.Body{
		"points": resultPoints,
		"total":  uint64(1)}

	mockRestUC.EXPECT().GetPoints(restID, uint64(1), uint64(1)).
		Return(resultPoints, uint64(1), nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewRestaurantHandler(privateGroup, publicGroup, reqValidator, mockRestUC, mwareC)

	target := "/api/v1/restaurants/" + strconv.Itoa(int(restID)) + "/points"
	req := httptest.NewRequest("GET", target, nil)
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
	expect, err := json.Marshal(expectResp)
	require.NoError(t, err)

	require.EqualValues(t, expect, result)
}

func TestRestaurantHandler_GetRestaurantsWithCloserPoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRestUC := mock_restaurants.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	testAddress := "Pushkina dom Kolotushkina"

	resultRests := []*models.Restaurant{
		{ID: 1, Name: "test1", Rating: 4.2, Image: "./default.jpg"},
	}

	expectResp := &tools.Body{
		"restaurants": resultRests,
		"total":       uint64(1)}

	mockRestUC.EXPECT().GetRestaurantsInServiceRadius(testAddress, uint64(1), uint64(1)).
		Return(resultRests, uint64(1), nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewRestaurantHandler(privateGroup, publicGroup, reqValidator, mockRestUC, mwareC)

	target := "/api/v1/restaurants_point"
	req := httptest.NewRequest("GET", target, nil)
	q := req.URL.Query()
	q.Add("page", "1")
	q.Add("count", "1")
	q.Add("address", testAddress)
	req.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()

	g.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("Status is not ok")
		return
	}

	result, err := ioutil.ReadAll(w.Result().Body)
	require.NoError(t, err)
	expect, err := json.Marshal(expectResp)
	require.NoError(t, err)

	require.EqualValues(t, expect, result)
}

func TestRestaurantHandler_GetReviews(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRestUC := mock_restaurants.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	restID := uint64(1)
	userID := uint64(1)

	sessRes := &models.Session{ID: 1, UserId: userID}
	userRes := &models.User{ID: userID, Role: "User"}

	resultReviews := []*models.Review{
		{
			ID:     1,
			RestID: 1,
			Text:   "Good test",
			Rate:   5,
			Author: userRes,
		},
	}

	expectResp := &tools.Body{
		"reviews": resultReviews,
		"current": resultReviews[0],
		"total":   uint64(1),
	}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(userID).Return(userRes, nil)
	mockRestUC.EXPECT().GetReviews(restID, userID, uint64(1), uint64(1)).
		Return(resultReviews, resultReviews[0], uint64(1), nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewRestaurantHandler(privateGroup, publicGroup, reqValidator, mockRestUC, mwareC)

	target := "/api/v1/restaurants/" + strconv.Itoa(int(restID)) + "/reviews"
	req := httptest.NewRequest("GET", target, nil)
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
	expect, err := json.Marshal(expectResp)
	require.NoError(t, err)

	require.EqualValues(t, expect, result)
}
