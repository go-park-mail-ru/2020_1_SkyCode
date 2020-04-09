package repository

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRepository_InsertInto(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewSessionRepository(db)

	testSess := &models.Session{
		UserId: 1,
		Token:  "test-token",
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectQuery("INSERT INTO sessions").
		WithArgs(testSess.UserId, testSess.Token).
		WillReturnRows(rows)

	err = repo.InsertInto(testSess)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}

	require.EqualValues(t, 1, testSess.ID)
}

func TestRepository_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewSessionRepository(db)

	testSess := &models.Session{
		ID:     1,
		UserId: 1,
		Token:  "test-token",
	}

	rows := sqlmock.NewRows([]string{"id", "userId", "token"}).
		AddRow(testSess.ID, testSess.UserId, testSess.Token)

	mock.ExpectQuery("SELECT").
		WithArgs(testSess.Token).
		WillReturnRows(rows)

	resultSess := &models.Session{
		Token: testSess.Token,
	}

	err = repo.Get(resultSess)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}

	require.EqualValues(t, testSess, resultSess)
}

func TestRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewSessionRepository(db)

	testSess := &models.Session{
		ID:     1,
		UserId: 1,
		Token:  "test-token",
	}

	rows := sqlmock.NewRows([]string{"token", "userId"}).
		AddRow(testSess.Token, testSess.UserId)

	mock.ExpectQuery("DELETE").
		WithArgs(testSess.ID).
		WillReturnRows(rows)

	resultSess := &models.Session{
		ID: testSess.ID,
	}

	err = repo.Delete(resultSess)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}

	require.EqualValues(t, testSess, resultSess)
}
