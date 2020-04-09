package orders

import "github.com/2020_1_Skycode/internal/models"

type Repository interface {
	Get(order *models.Order) error
	InsertOrder(order *models.Order) error
}
