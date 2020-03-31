package usecase

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/products"
	"github.com/2020_1_Skycode/internal/restaurants"
)

type RestarauntUseCase struct {
	restarauntRepo restaurants.Repository
	productRepo    products.Repository
}

func NewRestaurantsUseCase(rr restaurants.Repository, pr products.Repository) *RestarauntUseCase {
	return &RestarauntUseCase{
		restarauntRepo: rr,
		productRepo:    pr,
	}
}

func (rUC *RestarauntUseCase) GetRestaurants() ([]*models.Restaurant, error) {
	restList, err := rUC.restarauntRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return restList, nil
}

func (rUC *RestarauntUseCase) GetRestaurantByID(id uint64) (*models.Restaurant, error) {
	rest, err := rUC.restarauntRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	productList, err := rUC.productRepo.GetRestaurantProducts(id)
	if err != nil {
		return nil, err
	}

	rest.Products = productList

	return rest, nil
}
