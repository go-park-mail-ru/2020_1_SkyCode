package orders

import "github.com/2020_1_Skycode/internal/models"

type UseCase interface {
	CheckoutOrder(order *models.Order, ordProducts []*models.OrderProduct) error
	GetAllUserOrders(userID uint64, count uint64, page uint64) ([]*models.Order, uint64, error)
	GetOrderByID(orderID uint64, userID uint64) (*models.Order, error)
	DeleteOrder(orderID uint64, userID uint64) error
}
