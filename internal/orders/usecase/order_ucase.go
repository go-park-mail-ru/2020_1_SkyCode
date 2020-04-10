package usecase

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/orders"
)

type OrderUseCase struct {
	orderRepo orders.Repository
}

func NewOrderUseCase(orderRepo orders.Repository) *OrderUseCase {
	return &OrderUseCase{
		orderRepo: orderRepo,
	}
}

func (oU *OrderUseCase) CheckoutOrder(order *models.Order, ordProducts []*models.OrderProduct) error {
	if err := oU.orderRepo.InsertOrder(order, ordProducts); err != nil {
		return err
	}

	return nil
}

func (oU *OrderUseCase) GetAllUserOrders(userID uint64) ([]*models.Order, error) {
	userOrders, err := oU.orderRepo.GetAllByUserID(userID);

	if err != nil {
		return nil, err
	}

	return userOrders, nil
}

func (oU *OrderUseCase) GetOrderByID(orderID uint64, userID uint64) (*models.Order, error) {
	order, err := oU.orderRepo.GetByID(orderID, userID)

	if err != nil {
		return nil, err
	}

	return order, err
}

func (oU *OrderUseCase) DeleteOrder(orderID uint64, userID uint64) error {
	if err := oU.orderRepo.DeleteOrder(orderID, userID); err != nil {
		return err
	}

	return nil
}
