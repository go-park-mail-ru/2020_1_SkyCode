package usecase

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/restaurants"
)

type RestaurantUseCase struct {
	restaurantRepo restaurants.Repository
}

func NewRestaurantsUseCase(rr restaurants.Repository) *RestaurantUseCase {
	return &RestaurantUseCase{
		restaurantRepo: rr,
	}
}

func (rUC *RestaurantUseCase) GetRestaurants(count uint64, page uint64) ([]*models.Restaurant, uint64, error) {
	restList, total, err := rUC.restaurantRepo.GetAll(count, page)
	if err != nil {
		return nil, total, err
	}

	return restList, total, nil
}

func (rUC *RestaurantUseCase) GetRestaurantByID(id uint64) (*models.Restaurant, error) {
	rest, err := rUC.restaurantRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return rest, nil
}

func (rUC *RestaurantUseCase) CreateRestaurant(rest *models.Restaurant) error {
	if err := rUC.restaurantRepo.InsertInto(rest); err != nil {
		return err
	}

	return nil
}

func (rUC *RestaurantUseCase) UpdateRestaurant(rest *models.Restaurant) error {
	updRestaurant := &models.Restaurant{
		ID:          rest.ID,
		Name:        rest.Name,
		Description: rest.Description,
	}

	if err := rUC.restaurantRepo.Update(updRestaurant); err != nil {
		return err
	}

	return nil
}

func (rUC *RestaurantUseCase) UpdateImage(restID uint64, filename string) error {
	updRestaurant := &models.Restaurant{
		ID:    restID,
		Image: filename,
	}

	if err := rUC.restaurantRepo.UpdateImage(updRestaurant); err != nil {
		return err
	}

	return nil
}

func (rUC *RestaurantUseCase) Delete(restID uint64) error {
	if err := rUC.restaurantRepo.Delete(restID); err != nil {
		return err
	}

	return nil
}
