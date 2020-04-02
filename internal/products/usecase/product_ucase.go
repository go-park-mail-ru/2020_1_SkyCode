package usecase

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/products"
)

type ProductUseCase struct {
	productRepo products.Repository
}

func NewProductUseCase(pr products.Repository) *ProductUseCase {
	return &ProductUseCase{
		productRepo: pr,
	}
}

func (pUC *ProductUseCase) CreateProduct(product *models.Product) error {
	if err := pUC.productRepo.InsertInto(product); err != nil {
		return err
	}

	return nil
}

func (pUC *ProductUseCase) GetProductByID(id uint64) (*models.Product, error) {
	product, err := pUC.productRepo.GetProductByID(id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (pUC *ProductUseCase) UpdateProduct(product *models.Product) error {
	updProduct := &models.Product{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
	}
	if err := pUC.productRepo.Update(updProduct); err != nil {
		return err
	}

	return nil
}

func (pUC *ProductUseCase) UpdateProductImage(id uint64, path string) error {
	updProduct := &models.Product{
		ID:    id,
		Image: path,
	}

	if err := pUC.productRepo.UpdateImage(updProduct); err != nil {
		return err
	}

	return nil
}

func (pUC *ProductUseCase) DeleteProduct(id uint64) error {
	if err := pUC.productRepo.Delete(id); err != nil {
		return err
	}

	return nil
}
