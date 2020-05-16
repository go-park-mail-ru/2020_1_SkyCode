package products

import "github.com/2020_1_Skycode/internal/models"

type Repository interface {
	GetProductsByRestID(restID uint64) ([]*models.Product, error)
	GetProductByID(prodID uint64) (*models.Product, error)
	InsertInto(product *models.Product) error
	Update(product *models.Product) error
	UpdateImage(product *models.Product) error
	Delete(prodID uint64) error
}
