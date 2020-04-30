package restaurant_points

import "github.com/2020_1_Skycode/internal/models"

type UseCase interface {
	GetPoint(id uint64) (*models.RestaurantPoint, error)
	GetAllPoints() ([]*models.RestaurantPoint, error)
	Delete(id uint64) error
}
