package usecase

import (
	mock_geodata "github.com/2020_1_Skycode/internal/geodata/mocks"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGeoDataUseCase_CheckGeoPos(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockGeoDataRepo := mock_geodata.NewMockRepository(ctrl)

	testAddress := "Pushkina dom Kolotushkina"

	resultGeoData := &models.GeoPos{
		Latitude:  34.675,
		Longitude: 55.567,
	}

	mockGeoDataRepo.EXPECT().GetGeoPosByAddress(testAddress).Return(resultGeoData, nil)
	geodataUCase := NewGeoDataUseCase(mockGeoDataRepo)

	gd, err := geodataUCase.CheckGeoPos(testAddress)
	require.NoError(t, err)

	require.EqualValues(t, resultGeoData, gd)
}
