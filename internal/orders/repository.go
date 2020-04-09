package orders

import "github.com/2020_1_Skycode/internal/models"

type Repository interface {
	GetAllByUserID(userID uint64) ([]*models.Order, error)
	GetByID(orderID uint64, userID uint64) (*models.Order, error)
	InsertOrder(order *models.Order) error
	DeleteOrder(orderID uint64, userID uint64) error
}
