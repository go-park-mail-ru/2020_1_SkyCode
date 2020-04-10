package orders

import "github.com/2020_1_Skycode/internal/models"

type UseCase interface {
	CheckoutOrder(order *models.Order, ordProducts []*models.OrderProduct) error
	GetAllUserOrders(userID uint64) ([]*models.Order, error)
	GetOrderByID(orderID uint64, userID uint64) (*models.Order, error)
	DeleteOrder(orderID uint64, userID uint64) error
}
