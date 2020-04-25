package usecase

import "github.com/2020_1_Skycode/internal/geodata"

type GeoDataUseCase struct {
	GeoDataRepo geodata.Repository
}

func NewGeoDataUseCase(gdr geodata.Repository) geodata.UseCase {
	return &GeoDataUseCase{
		GeoDataRepo: gdr,
	}
}

func (gdUC *GeoDataUseCase) CheckGeoPos(address string) (bool, error) {
	if _, err := gdUC.GeoDataRepo.GetGeoPosByAddress(address); err != nil {
		return false, err
	}

	return true, nil
}
