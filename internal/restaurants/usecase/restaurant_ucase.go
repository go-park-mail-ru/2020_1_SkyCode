package usecase

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/products"
	"github.com/2020_1_Skycode/internal/restaurants"
)

type RestaurantUseCase struct {
	restaurantRepo restaurants.Repository
	productRepo    products.Repository
}

func NewRestaurantsUseCase(rr restaurants.Repository, pr products.Repository) *RestaurantUseCase {
	return &RestaurantUseCase{
		restaurantRepo: rr,
		productRepo:    pr,
	}
}

func (rUC *RestaurantUseCase) GetRestaurants() ([]*models.Restaurant, error) {
	restList, err := rUC.restaurantRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return restList, nil
}

func (rUC *RestaurantUseCase) GetRestaurantByID(id uint64) (*models.Restaurant, error) {
	rest, err := rUC.restaurantRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	productList, err := rUC.productRepo.GetProductsByRestID(id)
	if err != nil {
		return nil, err
	}

	rest.Products = productList

	return rest, nil
}
