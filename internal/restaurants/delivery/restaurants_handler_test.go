package delivery

import (
	"encoding/json"
	"github.com/2020_1_Skycode/internal/models"
	mock_restaurants "github.com/2020_1_Skycode/internal/restaurants/mocks"
	"github.com/2020_1_Skycode/internal/tools"
	_rValidator "github.com/2020_1_Skycode/internal/tools/requestValidator"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestRestaurantHandler_GetRestaurants(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_restaurants.NewMockUseCase(ctrl)

	resultRests := []*models.Restaurant{
		{ID: 1, Name: "test1", Rating: 4.2, Image: "./default.jpg"},
	}

	mockUC.EXPECT().GetRestaurants().Return(resultRests, nil)

	g := gin.New()
	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewRestaurantHandler(privateGroup, publicGroup, reqValidator, mockUC)

	req := httptest.NewRequest("GET", "/api/v1/restaurants", nil)
	w := httptest.NewRecorder()

	g.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("Status is not ok")
		return
	}

	type Result struct {
		Restaurants []*models.Restaurant
	}

	var result Result
	_ = json.NewDecoder(w.Result().Body).Decode(&result)

	require.EqualValues(t, resultRests, result.Restaurants)
}

func TestRestaurantHandler_GetRestaurantByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_restaurants.NewMockUseCase(ctrl)

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
	mockUC.EXPECT().GetRestaurantByID(restID).Return(resultRests, nil)

	g := gin.New()
	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewRestaurantHandler(privateGroup, publicGroup, reqValidator, mockUC)

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

	mockUC := mock_restaurants.NewMockUseCase(ctrl)

	reqRests := &models.Restaurant{
		Name:        "test1",
		Description: "asdasdqweasdqwe",
	}

	type restaurantRequest struct {
		Name        string `json:"name, omitempty" binding:"required" validate:"min=3"`
		Description string `json:"description, omitempty" binding:"required" validate:"min=10"`
	}

	reqRest := &restaurantRequest{
		Name:        reqRests.Name,
		Description: reqRests.Description,
	}

	j, err := json.Marshal(reqRest)
	require.NoError(t, err)

	expectResult := &tools.Message{Message: "success"}

	mockUC.EXPECT().CreateRestaurant(reqRests).Return(nil)

	g := gin.New()
	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewRestaurantHandler(privateGroup, publicGroup, reqValidator, mockUC)

	target := "/api/v1/restaurants"
	req, err := http.NewRequest("POST", target, strings.NewReader(string(j)))
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

func TestRestaurantHandler_DeleteRestaurant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_restaurants.NewMockUseCase(ctrl)

	restID := uint64(1)

	resultRests := &models.Restaurant{
		ID:       1,
		Name:     "test1",
		Rating:   4.2,
		Image:    "",
		Products: nil,
	}

	mockUC.EXPECT().GetRestaurantByID(restID).Return(resultRests, nil)
	mockUC.EXPECT().Delete(restID).Return(nil)

	expectResult := &tools.Message{Message: "success"}

	g := gin.New()
	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewRestaurantHandler(privateGroup, publicGroup, reqValidator, mockUC)

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

// Будет работать только на тачке

//func TestRestaurantHandler_UpdateImage(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockUC := mock_restaurants.NewMockUseCase(ctrl)
//
//	restID := uint64(1)
//
//	resultRests := &models.Restaurant{
//		ID:     1,
//		Name:   "test1",
//		Rating: 4.2,
//		Image:  "",
//		Products: nil,
//	}
//
//	if err := os.MkdirAll(tools.RestaurantImagesPath, 0777); err != nil {
//		t.Errorf("Error on image work: %s", err)
//		return
//	}
//
//	expectResult := &tools.Message{ Message:"success"}
//
//	mockUC.EXPECT().GetRestaurantByID(restID).Return(resultRests, nil)
//	mockUC.EXPECT().UpdateImage(restID, gomock.Any).Return(nil)
//
//	g := gin.New()
//	publicGroup := g.Group("/api/v1")
//	privateGroup := g.Group("/api/v1")
//	reqValidator := _rValidator.NewRequestValidator()
//
//	_ = NewRestaurantHandler(privateGroup, publicGroup, reqValidator, mockUC)
//
//	bodyUpdate := &bytes.Buffer{}
//	writer := multipart.NewWriter(bodyUpdate)
//	part, _ := writer.CreateFormFile("image", "testfile")
//
//	part.Write([]byte("SOME FILE CONTENT"))
//
//	writer.Close()
//
//	target := "/api/v1/restaurants/" + strconv.Itoa(int(restID)) + "/image"
//	req, err:= http.NewRequest("PUT", target, bodyUpdate)
//	require.NoError(t, err)
//	req.Header.Set("Content-Type", writer.FormDataContentType())
//	w := httptest.NewRecorder()
//
//	g.ServeHTTP(w, req)
//
//	if w.Code != http.StatusOK {
//		t.Error("Status is not ok")
//		return
//	}
//
//	if err = os.RemoveAll(tools.RestaurantImagesPath); err != nil {
//		t.Errorf("Error on image work: %s", err)
//		return
//	}
//
//	var result *tools.Message
//	_ = json.NewDecoder(w.Result().Body).Decode(&result)
//
//	require.EqualValues(t, expectResult, result)
//}

func TestRestaurantHandler_UpdateRestaurant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_restaurants.NewMockUseCase(ctrl)

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

	mockUC.EXPECT().UpdateRestaurant(reqRests).Return(nil)

	g := gin.New()
	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewRestaurantHandler(privateGroup, publicGroup, reqValidator, mockUC)

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
