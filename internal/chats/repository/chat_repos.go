package repository

import (
	"database/sql"
	"github.com/2020_1_Skycode/internal/models"
)

type ChatsRepository struct {
	db *sql.DB
}

func NewChatsRepository(db *sql.DB) *ChatsRepository {
	return &ChatsRepository{
		db: db,
	}
}

func (cR *ChatsRepository) InsertChatMessage(message *models.ChatMessage) error {
	if _, err := cR.db.Exec("insert into chat_messages (user_id, username, chat, message) values($1, $2, $3, %4)",
		message.UserID,
		message.UserName,
		message.ChatID,
		message.Message); err != nil {
		return err
	}

	return nil
}

func (cR *ChatsRepository) SelectMessagesByChatID(chatID string) ([]*models.ChatMessage, error) {
	var messages []*models.ChatMessage
	rows, err := cR.db.Query("select * from chat_messages where chat = $1",
		chatID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		message := &models.ChatMessage{}

		if err := rows.Scan(&message.UserID, &message.UserName, &message.ChatID, &message.Message, &message.Created); err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}
