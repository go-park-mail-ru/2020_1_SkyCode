package restaurants

import "github.com/2020_1_Skycode/internal/models"

type UseCase interface {
	GetRestaurants() ([]*models.Restaurant, error)
	GetRestaurantByID(id uint64) (*models.Restaurant, error)
}
