package repository

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestChatsRepository_InsertChatMessage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewChatsRepository(db)

	testMessage := &models.ChatMessage{
		UserID:   1,
		UserName: "testuser",
		ChatID:   "testidofchat",
		Message:  "Che s den'gami?",
	}

	mock.
		ExpectExec("insert into chat_messages").
		WithArgs(testMessage.UserID, testMessage.UserName, testMessage.ChatID, testMessage.Message).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.InsertChatMessage(testMessage)
	require.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}
}

func TestChatsRepository_SelectMessagesByChatID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewChatsRepository(db)

	chatID := "testidofchat"

	testMessage := &models.ChatMessage{
		UserID:   1,
		UserName: "testuser",
		ChatID:   chatID,
		Message:  "Che s den'gami?",
		Created:  time.Now(),
	}

	rows := sqlmock.NewRows([]string{"userid", "username", "chatid", "message", "created"}).
		AddRow(testMessage.UserID, testMessage.UserName, testMessage.ChatID, testMessage.Message, testMessage.Created)

	mock.
		ExpectQuery("select").
		WithArgs(chatID).
		WillReturnRows(rows)

	messages, err := repo.SelectMessagesByChatID(chatID)
	require.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}

	require.EqualValues(t, []*models.ChatMessage{testMessage}, messages)
}
