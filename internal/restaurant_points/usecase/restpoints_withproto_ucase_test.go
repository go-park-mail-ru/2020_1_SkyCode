package usecase

import (
	"github.com/2020_1_Skycode/internal/models"
	mock_restpoints "github.com/2020_1_Skycode/internal/restaurant_points/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"testing"
)

func TestRestPointsWithProtoUseCase_GetAllPoints(t *testing.T) {
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

	c := &grpc.ClientConn{}
	restPointsUCase := NewRestPointsWithProtoUseCase(restPointsRepo, c)

	result, err := restPointsUCase.GetAllPoints()
	require.NoError(t, err)

	require.EqualValues(t, expect, result)
}

func TestRestPointsWithProtoUseCase_GetPoint(t *testing.T) {
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

	c := &grpc.ClientConn{}
	restPointsUCase := NewRestPointsWithProtoUseCase(restPointsRepo, c)

	result, err := restPointsUCase.GetPoint(testID)
	require.NoError(t, err)

	require.EqualValues(t, expect, result)
}
