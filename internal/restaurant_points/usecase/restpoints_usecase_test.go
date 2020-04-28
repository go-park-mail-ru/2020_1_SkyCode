package usecase

import (
	"github.com/2020_1_Skycode/internal/models"
	mock_restpoints "github.com/2020_1_Skycode/internal/restaurant_points/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRestPointsUseCase_GetAllPoints(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	restPointsRepo := mock_restpoints.NewMockRepository(ctrl)

	expect := []*models.RestaurantPoint{
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

	restPointsRepo.EXPECT().GetAll().Return(expect, nil)

	restPointsUCase := NewRestPointsUseCase(restPointsRepo)
	result, err := restPointsUCase.GetAllPoints()
	require.NoError(t, err)

	require.EqualValues(t, expect, result)
}

func TestRestPointsUseCase_GetPoint(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	restPointsRepo := mock_restpoints.NewMockRepository(ctrl)

	testID := uint64(1)

	expect := &models.RestaurantPoint{
		ID:      testID,
		Address: "Pushkina dom Kolotushkina",
		MapPoint: &models.GeoPos{
			Longitude: 55.753227,
			Latitude:  37.619030,
		},
		ServiceRadius: 5,
		RestID:        1,
	}

	restPointsRepo.EXPECT().GetPointByID(testID).Return(expect, nil)

	restPointsUCase := NewRestPointsUseCase(restPointsRepo)
	result, err := restPointsUCase.GetPoint(testID)
	require.NoError(t, err)

	require.EqualValues(t, expect, result)
}

func TestRestPointsUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	restPointsRepo := mock_restpoints.NewMockRepository(ctrl)

	testID := uint64(1)
	expect := &models.RestaurantPoint{
		ID:      testID,
		Address: "Pushkina dom Kolotushkina",
		MapPoint: &models.GeoPos{
			Longitude: 55.753227,
			Latitude:  37.619030,
		},
		ServiceRadius: 5,
		RestID:        1,
	}

	restPointsRepo.EXPECT().GetPointByID(testID).Return(expect, nil)
	restPointsRepo.EXPECT().Delete(testID).Return(nil)

	restPointsUCase := NewRestPointsUseCase(restPointsRepo)
	err := restPointsUCase.Delete(testID)
	require.NoError(t, err)
}
