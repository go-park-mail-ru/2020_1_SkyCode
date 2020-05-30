package usecase

import (
	"database/sql"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/product_tags"
	"github.com/2020_1_Skycode/internal/products"
	"github.com/2020_1_Skycode/internal/tools"
)

type ProductUseCase struct {
	productRepo     products.Repository
	productTagsRepo product_tags.Repository
}

func NewProductUseCase(pr products.Repository, ptr product_tags.Repository) *ProductUseCase {
	return &ProductUseCase{
		productRepo:     pr,
		productTagsRepo: ptr,
	}
}

func (pUC *ProductUseCase) CreateProduct(product *models.Product) error {
	if _, err := pUC.productTagsRepo.GetByID(product.Tag); err != nil {
		if err == sql.ErrNoRows {
			return tools.ProductTagNotFound
		}

		return err
	}

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

func (pUC *ProductUseCase) GetProductsByRestaurantID(
	id uint64, count uint64, page uint64) ([]*models.Product, uint64, error) {
	productList, total, err := pUC.productRepo.GetProductsByRestID(id, count, page)
	if err != nil {
		return nil, 0, err
	}

	return productList, total, nil
}

func (pUC *ProductUseCase) UpdateProduct(product *models.Product) error {
	if _, err := pUC.productTagsRepo.GetByID(product.Tag); err != nil {
		if err == sql.ErrNoRows {
			return tools.ProductTagNotFound
		}

		return err
	}

	updProduct := &models.Product{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
		Tag:   product.Tag,
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
