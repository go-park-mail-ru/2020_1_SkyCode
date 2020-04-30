package usecase

import (
	"github.com/2020_1_Skycode/internal/geodata"
	"github.com/2020_1_Skycode/internal/models"
)

type GeoDataUseCase struct {
	GeoDataRepo geodata.Repository
}

func NewGeoDataUseCase(gdr geodata.Repository) geodata.UseCase {
	return &GeoDataUseCase{
		GeoDataRepo: gdr,
	}
}

func (gdUC *GeoDataUseCase) CheckGeoPos(address string) (*models.GeoPos, error) {
	pos, err := gdUC.GeoDataRepo.GetGeoPosByAddress(address)
	if err != nil {
		return nil, err
	}

	return pos, nil
}
