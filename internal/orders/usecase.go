package orders

import "github.com/2020_1_Skycode/internal/models"

type UseCase interface {
	CheckoutOrder(order *models.Order) error
}
