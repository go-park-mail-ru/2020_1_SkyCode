package notifications

import "github.com/2020_1_Skycode/internal/models"

type UseCase interface {
	GetAllByUser(userID uint64) ([]*models.Notification, error)
	ChangeReadStatus(ID uint64, userID uint64) error
	Delete(ID uint64, userID uint64) error
}
