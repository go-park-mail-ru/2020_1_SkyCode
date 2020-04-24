package geodata

import "github.com/2020_1_Skycode/internal/models"

type Repository interface {
	GetGeoPosByAddress(string) (*models.GeoPos, error)
}
