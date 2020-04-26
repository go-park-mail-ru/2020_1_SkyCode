package usecase

import (
	"github.com/2020_1_Skycode/internal/models"
	mock_products "github.com/2020_1_Skycode/internal/products/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProductUseCase_CreateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockProductRepo := mock_products.NewMockRepository(ctrl)

	testProd := &models.Product{
		Name:   "test",
		Price:  3.22,
		Image:  "./default.jpg",
		RestId: 1,
	}

	mockProductRepo.EXPECT().InsertInto(testProd).Return(nil)
	prodUCase := NewProductUseCase(mockProductRepo)

	err := prodUCase.CreateProduct(testProd)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
}

func TestProductUseCase_DeleteProduct(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockProductRepo := mock_products.NewMockRepository(ctrl)
	prodID := uint64(1)

	mockProductRepo.EXPECT().Delete(prodID).Return(nil)
	prodUCase := NewProductUseCase(mockProductRepo)

	err := prodUCase.DeleteProduct(prodID)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
}

func TestProductUseCase_GetProductByID(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockProductRepo := mock_products.NewMockRepository(ctrl)

	testProd := &models.Product{
		ID:     1,
		Name:   "test",
		Price:  3.22,
		Image:  "./default.jpg",
		RestId: 1,
	}

	mockProductRepo.EXPECT().GetProductByID(testProd.ID).Return(testProd, nil)
	prodUCase := NewProductUseCase(mockProductRepo)

	resutlProd, err := prodUCase.GetProductByID(testProd.ID)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}

	require.EqualValues(t, testProd, resutlProd)
}

func TestProductUseCase_GetProductsByRestaurantID(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockProductRepo := mock_products.NewMockRepository(ctrl)

	testProdList := []*models.Product{
		{ID: 1, Name: "test", Price: 3.22, Image: "./default.jpg"},
		{ID: 2, Name: "test2", Price: 2.50, Image: "./nonDefault.jpg"},
	}
	restID := uint64(1)

	mockProductRepo.EXPECT().GetProductsByRestID(restID, uint64(1), uint64(1)).Return(testProdList, uint64(1), nil)
	prodUCase := NewProductUseCase(mockProductRepo)

	resultList, total, err := prodUCase.GetProductsByRestaurantID(restID, uint64(1), uint64(1))
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}

	require.EqualValues(t, testProdList, resultList)
	require.EqualValues(t, uint64(1), total)
}

func TestProductUseCase_UpdateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockProductRepo := mock_products.NewMockRepository(ctrl)

	testProd := &models.Product{
		Name:  "test",
		Price: 3.22,
	}

	mockProductRepo.EXPECT().Update(testProd).Return(nil)
	prodUCase := NewProductUseCase(mockProductRepo)

	err := prodUCase.UpdateProduct(testProd)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
}

func TestProductUseCase_UpdateProductImage(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockProductRepo := mock_products.NewMockRepository(ctrl)

	testProd := &models.Product{
		ID:    uint64(1),
		Image: "./default.jpg",
	}

	mockProductRepo.EXPECT().UpdateImage(testProd).Return(nil)
	prodUCase := NewProductUseCase(mockProductRepo)

	err := prodUCase.UpdateProductImage(testProd.ID, testProd.Image)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
}
