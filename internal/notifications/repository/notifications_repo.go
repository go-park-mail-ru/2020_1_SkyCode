package repository

import (
	"database/sql"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/notifications"
)

type NotificationsRepository struct {
	db *sql.DB
}

func NewNotificationsRepository(db *sql.DB) notifications.Repository {
	return &NotificationsRepository{
		db: db,
	}
}

func (nr *NotificationsRepository) InsertInto(note *models.Notification) error {
	if err := nr.db.QueryRow("INSERT INTO order_notifications (user_id, order_id, order_status) "+
		"VALUES ($1, $2, $3) RETURNING id", note.UserID, note.OrderID, note.Status).Scan(&note.ID); err != nil {
		return err
	}

	return nil
}

func (nr *NotificationsRepository) GetByID(ID uint64) (*models.Notification, error) {
	note := &models.Notification{}
	if err := nr.db.QueryRow("SELECT id, user_id, order_id, unread, order_status, getting_time "+
		"FROM order_notifications WHERE id = $1", ID).Scan(&note.ID, &note.UserID, &note.OrderID,
		&note.UnreadStatus, &note.Status, &note.DateTime); err != nil {
		return nil, err
	}

	return note, nil
}

func (nr *NotificationsRepository) ChangeReadStatus(ID uint64) error {
	if _, err := nr.db.Exec("UPDATE order_notifications SET unread = FALSE "+
		"WHERE id = $1", ID); err != nil {
		return err
	}

	return nil
}

func (nr *NotificationsRepository) GetAllUserNotifications(userID uint64) ([]*models.Notification, error) {
	rows, err := nr.db.Query("SELECT id, user_id, order_id, unread, order_status, getting_time "+
		"FROM order_notifications WHERE user_id = $1 ORDER BY getting_time DESC", userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var notes []*models.Notification
	for rows.Next() {
		note := &models.Notification{}

		if err := rows.Scan(&note.ID, &note.UserID, &note.OrderID, &note.UnreadStatus, &note.Status,
			&note.DateTime); err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}

	return notes, nil
}

func (nr *NotificationsRepository) Delete(ID uint64) error {
	if _, err := nr.db.Exec("DELETE FROM order_notifications WHERE id = $1", ID); err != nil {
		return err
	}

	return nil
}
