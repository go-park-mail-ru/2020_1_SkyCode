package products

import "github.com/2020_1_Skycode/internal/models"

type UseCase interface {
	CreateProduct(product *models.Product) error
	GetProductByID(id uint64) (*models.Product, error)
	UpdateProduct(product *models.Product) error
	UpdateProductImage(id uint64, path string) error
	DeleteProduct(id uint64) error
}
