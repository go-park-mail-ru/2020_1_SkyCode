package products

import "github.com/2020_1_Skycode/internal/models"

type Repository interface {
	GetRestaurantProducts(restID uint64) ([]*models.Product, error)
}
