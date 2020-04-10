package restaurants

import "github.com/2020_1_Skycode/internal/models"

type UseCase interface {
	GetRestaurants(count uint64, page uint64) ([]*models.Restaurant, uint64, error)
	GetRestaurantByID(id uint64) (*models.Restaurant, error)
	CreateRestaurant(rest *models.Restaurant) error
	UpdateRestaurant(rest *models.Restaurant) error
	UpdateImage(restID uint64, filename string) error
	Delete(restID uint64) error
}
