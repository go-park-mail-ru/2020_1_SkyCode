package usecase

import (
	"github.com/2020_1_Skycode/internal/models"
	mock_orders "github.com/2020_1_Skycode/internal/orders/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOrderUseCase_CheckoutOrder(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockOrdersRepo := mock_orders.NewMockRepository(ctrl)

	testOrderProduct := []*models.OrderProduct{
		{
			ProductID: 1,
			Count:     1,
		},
	}

	testOrder := &models.Order{
		UserID:    1,
		Address:   "Pushkina, dom Kolotushkina",
		Phone:     "89765432211",
		Comment:   "Faster",
		PersonNum: 2,
		Products:  nil,
		Price:     100.0,
	}

	mockOrdersRepo.EXPECT().InsertOrder(testOrder, testOrderProduct).Return(nil)
	ordersUCase := NewOrderUseCase(mockOrdersRepo)

	err := ordersUCase.CheckoutOrder(testOrder, testOrderProduct)
	require.NoError(t, err)
}

func TestOrderUseCase_GetOrderByID(t *testing.T) {
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

	mockOrdersRepo.EXPECT().GetByID(uint64(1), testOrder.UserID).Return(testOrder, nil)
	ordersUCase := NewOrderUseCase(mockOrdersRepo)

	result, err := ordersUCase.GetOrderByID(uint64(1), testOrder.UserID)

	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}

	require.EqualValues(t, testOrder, result)
}

func TestOrderUseCase_GetAllUserOrders(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockOrdersRepo := mock_orders.NewMockRepository(ctrl)

	testOrders := []*models.Order{
		{
			UserID:    1,
			Address:   "Pushkina, dom Kolotushkina",
			Phone:     "89765432211",
			Comment:   "Faster",
			PersonNum: 2,
			Products:  nil,
			Price:     100.0,
		},
	}

	expectTotal := uint64(1)

	mockOrdersRepo.EXPECT().GetAllByUserID(testOrders[0].UserID, uint64(1), uint64(1)).
		Return(testOrders, expectTotal, nil)

	ordersUCase := NewOrderUseCase(mockOrdersRepo)

	result, total, err := ordersUCase.GetAllUserOrders(testOrders[0].UserID, uint64(1), uint64(1))
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}

	require.EqualValues(t, testOrders, result)
	require.EqualValues(t, expectTotal, total)
}

func TestOrderUseCase_DeleteOrder(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockOrdersRepo := mock_orders.NewMockRepository(ctrl)

	expectOrderID := uint64(1)
	expectUserID := uint64(1)

	mockOrdersRepo.EXPECT().DeleteOrder(expectOrderID, expectUserID).Return(nil)

	ordersUCase := NewOrderUseCase(mockOrdersRepo)
	err := ordersUCase.DeleteOrder(expectOrderID, expectUserID)

	require.NoError(t, err)
}
