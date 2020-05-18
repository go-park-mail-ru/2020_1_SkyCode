package product_tags

import "github.com/2020_1_Skycode/internal/models"

type Repository interface {
	InsertInto(tag *models.ProductTag) error
	GetByID(ID uint64) (*models.ProductTag, error)
	GetByRestID(restID uint64) ([]*models.ProductTag, error)
	Delete(ID uint64) error
}
