package restaurants

import "github.com/2020_1_Skycode/internal/models"

type Repository interface {
	GetAll() ([]*models.Restaurant, error)
	GetByID(id uint64) (*models.Restaurant, error)
	InsertInto(rest *models.Restaurant) error
	Update(rest *models.Restaurant) error
	UpdateImage(rest *models.Restaurant) error
	Delete(restID uint64) error
}
