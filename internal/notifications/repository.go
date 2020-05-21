package notifications

import "github.com/2020_1_Skycode/internal/models"

type Repository interface {
	InsertInto(note *models.Notification) error
	GetByID(ID uint64) (*models.Notification, error)
	GetAllUserNotifications(userID uint64) ([]*models.Notification, error)
	ChangeReadStatus(ID uint64) error
	Delete(noteID uint64) error
}
