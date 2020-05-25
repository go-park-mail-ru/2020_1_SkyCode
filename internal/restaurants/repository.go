package restaurants

import "github.com/2020_1_Skycode/internal/models"

type Repository interface {
	GetAll(count uint64, page uint64, tagID uint64) ([]*models.Restaurant, uint64, error)
	GetRecommendationsByOrder(userID uint64, count uint64) ([]*models.Restaurant, error)
	GetByID(id uint64) (*models.Restaurant, error)
	GetByName(name string) (*models.Restaurant, error)
	GetAllInServiceRadius(pos *models.GeoPos, count, page, tagID uint64) ([]*models.Restaurant, uint64, error)
	GetRecommendationsInRadius(pos *models.GeoPos, userID uint64, count uint64) ([]*models.Restaurant, error)
	InsertInto(rest *models.Restaurant) error
	Update(rest *models.Restaurant) error
	UpdateImage(rest *models.Restaurant) error
	Delete(restID uint64) error
}
