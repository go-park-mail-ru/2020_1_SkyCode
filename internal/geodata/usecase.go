package geodata

import "github.com/2020_1_Skycode/internal/models"

type UseCase interface {
	CheckGeoPos(address string) (*models.GeoPos, error)
}
