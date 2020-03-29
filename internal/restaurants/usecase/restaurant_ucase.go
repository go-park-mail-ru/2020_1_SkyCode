package usecase

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/restaurants"
)

type RestarauntUseCase struct {
	restarauntRepo restaurants.Repository
}

func NewRestaurantsUseCase(rr restaurants.Repository) *RestarauntUseCase {
	return &RestarauntUseCase{
		restarauntRepo: rr,
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

	return rest, nil
}
