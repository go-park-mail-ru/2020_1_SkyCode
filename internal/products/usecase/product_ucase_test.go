package usecase

import (
	"github.com/2020_1_Skycode/internal/models"
	mock_prodtags "github.com/2020_1_Skycode/internal/product_tags/mocks"
	mock_products "github.com/2020_1_Skycode/internal/products/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProductUseCase_CreateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockProductRepo := mock_products.NewMockRepository(ctrl)
	mockProdTagsRepo := mock_prodtags.NewMockRepository(ctrl)

	testTag := &models.ProductTag{
		ID:     1,
		Name:   "Tag",
		RestID: 1,
	}

	testProd := &models.Product{
		Name:   "test",
		Price:  3.22,
		Image:  "./default.jpg",
		RestId: 1,
		Tag:    testTag.ID,
	}

	mockProdTagsRepo.EXPECT().GetByID(testTag.ID).Return(testTag, nil)
	mockProductRepo.EXPECT().InsertInto(testProd).Return(nil)
	prodUCase := NewProductUseCase(mockProductRepo, mockProdTagsRepo)

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
	mockProdTagsRepo := mock_prodtags.NewMockRepository(ctrl)
	prodID := uint64(1)

	mockProductRepo.EXPECT().Delete(prodID).Return(nil)
	prodUCase := NewProductUseCase(mockProductRepo, mockProdTagsRepo)

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
	mockProdTagsRepo := mock_prodtags.NewMockRepository(ctrl)

	testProd := &models.Product{
		ID:     1,
		Name:   "test",
		Price:  3.22,
		Image:  "./default.jpg",
		RestId: 1,
	}

	mockProductRepo.EXPECT().GetProductByID(testProd.ID).Return(testProd, nil)
	prodUCase := NewProductUseCase(mockProductRepo, mockProdTagsRepo)

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
	mockProdTagsRepo := mock_prodtags.NewMockRepository(ctrl)

	testProdList := []*models.Product{
		{ID: 1, Name: "test", Price: 3.22, Image: "./default.jpg"},
		{ID: 2, Name: "test2", Price: 2.50, Image: "./nonDefault.jpg"},
	}
	restID := uint64(1)

	mockProductRepo.EXPECT().GetProductsByRestID(restID).Return(testProdList, nil)
	prodUCase := NewProductUseCase(mockProductRepo, mockProdTagsRepo)

	resultList, err := prodUCase.GetProductsByRestaurantID(restID)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}

	require.EqualValues(t, testProdList, resultList)
}

func TestProductUseCase_UpdateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockProductRepo := mock_products.NewMockRepository(ctrl)
	mockProdTagsRepo := mock_prodtags.NewMockRepository(ctrl)

	testTag := &models.ProductTag{
		ID:     1,
		Name:   "Tag",
		RestID: 1,
	}

	testProd := &models.Product{
		Name:  "test",
		Price: 3.22,
		Tag:   testTag.ID,
	}

	mockProdTagsRepo.EXPECT().GetByID(testTag.ID).Return(testTag, nil)
	mockProductRepo.EXPECT().Update(testProd).Return(nil)
	prodUCase := NewProductUseCase(mockProductRepo, mockProdTagsRepo)

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
	mockProdTagsRepo := mock_prodtags.NewMockRepository(ctrl)

	testProd := &models.Product{
		ID:    uint64(1),
		Image: "./default.jpg",
	}

	mockProductRepo.EXPECT().UpdateImage(testProd).Return(nil)
	prodUCase := NewProductUseCase(mockProductRepo, mockProdTagsRepo)

	err := prodUCase.UpdateProductImage(testProd.ID, testProd.Image)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
}
