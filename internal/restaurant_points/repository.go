package restaurant_points

import "github.com/2020_1_Skycode/internal/models"

type Repository interface {
	InsertInto(point *models.RestaurantPoint) error
	Delete(id uint64) error
	GetPointsByRestID(restID, count, page uint64) ([]*models.RestaurantPoint, uint64, error)
	GetPointByID(id uint64) (*models.RestaurantPoint, error)
}
