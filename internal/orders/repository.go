package orders

import "github.com/2020_1_Skycode/internal/models"

type Repository interface {
	GetAllByUserID(userID uint64, count uint64, page uint64) ([]*models.Order, uint64, error)
	GetAllByRestID(restID uint64, count uint64, page uint64) ([]*models.Order, uint64, error)
	GetByID(orderID uint64) (*models.Order, error)
	ChangeStatus(orderID uint64, status string) error
	InsertOrder(order *models.Order, ordProducts []*models.OrderProduct) error
	DeleteOrder(orderID uint64, userID uint64) error
}
