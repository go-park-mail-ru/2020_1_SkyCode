package restaurants_tags

import "github.com/2020_1_Skycode/internal/models"

type Repository interface {
	InsertInto(tag *models.RestTag) error
	GetByID(id uint64) (*models.RestTag, error)
	GetAll() ([]*models.RestTag, error)
	Update(tag *models.RestTag) error
	Delete(id uint64) error
}
