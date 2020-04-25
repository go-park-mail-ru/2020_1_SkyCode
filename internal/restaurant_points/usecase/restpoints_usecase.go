package usecase

import (
	"database/sql"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/restaurant_points"
	"github.com/2020_1_Skycode/internal/tools"
)

type RestPointsUseCase struct {
	RestPointsRepo restaurant_points.Repository
}

func NewRestPointsUseCase(rpr restaurant_points.Repository) restaurant_points.UseCase {
	return &RestPointsUseCase{
		RestPointsRepo: rpr,
	}
}

func (rpUC *RestPointsUseCase) GetPoint(id uint64) (*models.RestaurantPoint, error) {
	returnPoint, err := rpUC.RestPointsRepo.GetPointByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, tools.RestPointNotFound
		}

		return nil, err
	}

	return returnPoint, nil
}

func (rpUC *RestPointsUseCase) GetAllPoints() ([]*models.RestaurantPoint, error) {
	returnPoints, err := rpUC.RestPointsRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return returnPoints, err
}

func (rpUC *RestPointsUseCase) Delete(id uint64) error {

	if _, err := rpUC.RestPointsRepo.GetPointByID(id); err != nil {
		if err == sql.ErrNoRows {
			return tools.RestPointNotFound
		}

		return err
	}

	if err := rpUC.RestPointsRepo.Delete(id); err != nil {
		return err
	}

	return nil
}
