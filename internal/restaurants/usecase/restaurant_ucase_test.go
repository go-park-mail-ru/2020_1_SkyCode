package usecase

import (
	"github.com/2020_1_Skycode/internal/models"
	mock_products "github.com/2020_1_Skycode/internal/products/mocks"
	mock_restaurants "github.com/2020_1_Skycode/internal/restaurants/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetRestaraunts(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRestaurantsRepo := mock_restaurants.NewMockRepository(ctrl)

	resultRests := []*models.Restaurant{
		{ID: 1, Name: "test1", Rating: 4.2, Image: "./default.jpg"},
	}

	mockRestaurantsRepo.EXPECT().GetAll().Return(resultRests, nil)
	restUCase := NewRestaurantsUseCase(mockRestaurantsRepo, nil)

	rests, err := restUCase.GetRestaurants()

	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}

	require.EqualValues(t, resultRests, rests)
}

func TestGetRestaurantByID(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRestaurantsRepo := mock_restaurants.NewMockRepository(ctrl)
	mockProductsRepo := mock_products.NewMockRepository(ctrl)

	resultProds := []*models.Product{
		{ID: 1, Name: "testProd", Price: 2.50, Image: "./default_prod.jpg"},
	}
	resultRest := models.Restaurant{
		ID: uint64(1), Name: "test1", Description: "smthng",
		Rating: 4.2, Image: "./default.jpg",
	}
	var elemID uint64 = 1

	expectRest := models.Restaurant{
		ID: uint64(1), Name: "test1", Description: "smthng",
		Rating: 4.2, Image: "./default.jpg", Products: resultProds,
	}

	mockRestaurantsRepo.EXPECT().GetByID(elemID).Return(&resultRest, nil)
	mockProductsRepo.EXPECT().GetRestaurantProducts(elemID).Return(resultProds, nil)
	restUCase := NewRestaurantsUseCase(mockRestaurantsRepo, mockProductsRepo)

	rest, err := restUCase.GetRestaurantByID(elemID)

	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}

	require.EqualValues(t, &expectRest, rest)
}
