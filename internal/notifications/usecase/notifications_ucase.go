package usecase

import (
	"database/sql"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/notifications"
	"github.com/2020_1_Skycode/internal/tools"
)

type NotificationsUseCase struct {
	notificationsRepo notifications.Repository
}

func NewNotificationsUseCase(nr notifications.Repository) notifications.UseCase {
	return &NotificationsUseCase{
		notificationsRepo: nr,
	}
}

func (nUC *NotificationsUseCase) GetAllByUser(userID uint64) ([]*models.Notification, error) {
	notes, err := nUC.notificationsRepo.GetAllUserNotifications(userID)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (nUC *NotificationsUseCase) ChangeReadStatus(ID uint64, userID uint64) error {
	note, err := nUC.notificationsRepo.GetByID(ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return tools.NotificationNotFound
		}

		return err
	}

	if note.UserID != userID {
		return tools.PermissionError
	}

	if !note.UnreadStatus {
		return nil
	}

	if err := nUC.notificationsRepo.ChangeReadStatus(ID); err != nil {
		return err
	}

	return nil
}

func (nUC *NotificationsUseCase) Delete(ID uint64, userID uint64) error {
	note, err := nUC.notificationsRepo.GetByID(ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return tools.NotificationNotFound
		}

		return err
	}

	if note.UserID != userID {
		return tools.PermissionError
	}

	if err := nUC.notificationsRepo.Delete(ID); err != nil {
		return err
	}

	return nil
}
