package restaurant_points

import "github.com/2020_1_Skycode/internal/models"

type Repository interface {
	InsertInto(point *models.RestaurantPoint) error
	Delete(id uint64) error
	GetAll() ([]*models.RestaurantPoint, error)
	GetPointsByRestID(restID, count, page uint64) ([]*models.RestaurantPoint, uint64, error)
	GetCloserPointByRestID(restID uint64, pos *models.GeoPos) (*models.RestaurantPoint, error)
	GetPointByID(id uint64) (*models.RestaurantPoint, error)
}
