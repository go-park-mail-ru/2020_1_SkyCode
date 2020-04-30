package usecase

import (
	"context"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/products"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/2020_1_Skycode/tools/protobuf/adminwork"
	"google.golang.org/grpc"
)

type ProductWithProtoUseCase struct {
	productRepo  products.Repository
	adminManager adminwork.RestaurantAdminWorkerClient
}

func NewProductWithProtoUseCase(pr products.Repository, conn *grpc.ClientConn) products.UseCase {
	return &ProductWithProtoUseCase{
		productRepo:  pr,
		adminManager: adminwork.NewRestaurantAdminWorkerClient(conn),
	}
}

func (pUC *ProductWithProtoUseCase) CreateProduct(product *models.Product) error {
	answ, err := pUC.adminManager.CreateProduct(
		context.Background(),
		&adminwork.ProtoProduct{
			Name:      product.Name,
			Price:     product.Price,
			ImagePath: product.Image,
			RestID:    product.RestId,
		})

	if answ.ID != tools.OK {
		if answ.ID == tools.DoesntExist {
			return tools.RestaurantNotFoundError
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (pUC *ProductWithProtoUseCase) GetProductByID(id uint64) (*models.Product, error) {
	product, err := pUC.productRepo.GetProductByID(id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (pUC *ProductWithProtoUseCase) GetProductsByRestaurantID(
	id uint64, count uint64, page uint64) ([]*models.Product, uint64, error) {
	productList, total, err := pUC.productRepo.GetProductsByRestID(id, count, page)
	if err != nil {
		return nil, 0, err
	}

	return productList, total, nil
}

func (pUC *ProductWithProtoUseCase) UpdateProduct(product *models.Product) error {
	answ, err := pUC.adminManager.UpdateProduct(
		context.Background(),
		&adminwork.ProtoProduct{
			ID:    product.ID,
			Name:  product.Name,
			Price: product.Price,
		})

	if answ.ID != tools.OK {
		if answ.ID == tools.DoesntExist {
			return tools.ProductNotFoundError
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (pUC *ProductWithProtoUseCase) UpdateProductImage(id uint64, path string) error {
	answ, err := pUC.adminManager.UpdateProductImage(
		context.Background(),
		&adminwork.ProtoImage{
			ID:        id,
			ImagePath: path,
		})

	if answ.ID != tools.OK {
		if answ.ID == tools.DoesntExist {
			return tools.ProductNotFoundError
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (pUC *ProductWithProtoUseCase) DeleteProduct(id uint64) error {
	answ, err := pUC.adminManager.DeleteProduct(
		context.Background(),
		&adminwork.ProtoID{
			ID: id,
		})

	if answ.ID != tools.OK {
		if answ.ID == tools.DoesntExist {
			return tools.ProductNotFoundError
		}
		if err != nil {
			return err
		}
	}

	return nil
}
