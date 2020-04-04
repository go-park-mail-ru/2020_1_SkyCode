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

func (oU *OrderUseCase) CheckoutOrder(order *models.Order) error {
	if err := oU.orderRepo.InsertOrder(order); err != nil {
		return err
	}

	return nil
}
