package usecase

import (
	"github.com/2020_1_Skycode/internal/models"
	mock_orders "github.com/2020_1_Skycode/internal/orders/mocks"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestOrderUseCase_CheckoutOrder(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockOrdersRepo := mock_orders.NewMockRepository(ctrl)

	testOrder := &models.Order{
		UserID:    1,
		Address:   "Pushkina, dom Kolotushkina",
		Phone:     "89765432211",
		Comment:   "Faster",
		PersonNum: 2,
		Products:  nil,
		Price:     100.0,
	}

	mockOrdersRepo.EXPECT().InsertOrder(testOrder).Return(nil)
	ordersUCase := NewOrderUseCase(mockOrdersRepo)

	if err := ordersUCase.CheckoutOrder(testOrder); err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
}
