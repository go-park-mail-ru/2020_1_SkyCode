package usecase

import (
	"github.com/2020_1_Skycode/internal/models"
	mock_restaurants "github.com/2020_1_Skycode/internal/restaurants/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRestaurantUseCase_GetRestaurants(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRestaurantsRepo := mock_restaurants.NewMockRepository(ctrl)

	resultRests := []*models.Restaurant{
		{ID: 1, Name: "test1", Rating: 4.2, Image: "./default.jpg"},
	}

	mockRestaurantsRepo.EXPECT().GetAll().Return(resultRests, nil)
	restUCase := NewRestaurantsUseCase(mockRestaurantsRepo)

	rests, err := restUCase.GetRestaurants()

	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}

	require.EqualValues(t, resultRests, rests)
}

func TestRestaurantUseCase_GetRestaurantByID(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRestaurantsRepo := mock_restaurants.NewMockRepository(ctrl)

	resultRest := &models.Restaurant{
		ID: uint64(1), ManagerID: uint64(1), Name: "test1",
		Description: "smthng", Rating: 4.2, Image: "./default.jpg",
	}
	var elemID uint64 = 1

	mockRestaurantsRepo.EXPECT().GetByID(elemID).Return(resultRest, nil)
	restUCase := NewRestaurantsUseCase(mockRestaurantsRepo)

	rest, err := restUCase.GetRestaurantByID(elemID)

	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}

	require.EqualValues(t, resultRest, rest)
}

func TestRestaurantUseCase_CreateRestaurant(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRestaurantsRepo := mock_restaurants.NewMockRepository(ctrl)

	testRest := &models.Restaurant{
		ManagerID:   1,
		Name:        "test",
		Description: "smthng",
		Rating:      4.2,
		Image:       "./default.jpg",
	}

	mockRestaurantsRepo.EXPECT().InsertInto(testRest).Return(nil)
	restUCase := NewRestaurantsUseCase(mockRestaurantsRepo)

	err := restUCase.CreateRestaurant(testRest)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
}

func TestRestaurantUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	restID := uint64(1)

	mockRestaurantsRepo := mock_restaurants.NewMockRepository(ctrl)
	mockRestaurantsRepo.EXPECT().Delete(restID).Return(nil)
	restUCase := NewRestaurantsUseCase(mockRestaurantsRepo)

	err := restUCase.Delete(restID)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
}

func TestRestaurantUseCase_UpdateImage(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRestaurantsRepo := mock_restaurants.NewMockRepository(ctrl)

	testRest := &models.Restaurant{
		ID:    uint64(1),
		Image: "./default.jpg",
	}

	mockRestaurantsRepo.EXPECT().UpdateImage(testRest).Return(nil)
	restUCase := NewRestaurantsUseCase(mockRestaurantsRepo)

	err := restUCase.UpdateImage(testRest.ID, testRest.Image)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
}

func TestRestaurantUseCase_UpdateRestaurant(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRestaurantsRepo := mock_restaurants.NewMockRepository(ctrl)
	testRest := &models.Restaurant{
		ID:          1,
		Name:        "test",
		Description: "smthng",
	}

	mockRestaurantsRepo.EXPECT().Update(testRest).Return(nil)
	restUCase := NewRestaurantsUseCase(mockRestaurantsRepo)

	err := restUCase.UpdateRestaurant(testRest)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
}
