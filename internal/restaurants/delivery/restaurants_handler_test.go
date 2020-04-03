package delivery

import (
	"encoding/json"
	"github.com/2020_1_Skycode/internal/models"
	mock_restaurants "github.com/2020_1_Skycode/internal/restaurants/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestGetRestaurants(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_restaurants.NewMockUseCase(ctrl)

	resultRests := []*models.Restaurant{
		{ID: 1, Name: "test1", Rating: 4.2, Image: "./default.jpg"},
	}

	mockUC.EXPECT().GetRestaurants().Return(resultRests, nil)

	g := gin.New()
	_ = NewRestaurantHandler(g, mockUC)

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

func TestGetRestaurantByID(t *testing.T) {
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
	_ = NewRestaurantHandler(g, mockUC)

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
