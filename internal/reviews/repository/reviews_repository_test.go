package repository

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestReviewsRepository_CreateReview(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	testReview := &models.Review{
		RestID:       1,
		Text:         "Good",
		Author:       &models.User{ID: 1},
		CreationDate: time.Now(),
		Rate:         5,
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(uint64(1))

	mock.ExpectQuery("INSERT INTO").
		WithArgs(testReview.RestID, testReview.Author.ID, testReview.Text,
			testReview.CreationDate, testReview.Rate).WillReturnRows(rows)

	mockRepo := NewReviewsRepository(db)
	err = mockRepo.CreateReview(testReview)
	require.NoError(t, err)
}

func TestReviewsRepository_DeleteReview(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	testID := uint64(1)

	mock.ExpectExec("DELETE").
		WithArgs(testID).WillReturnResult(sqlmock.NewResult(0, 1))

	mockRepo := NewReviewsRepository(db)
	err = mockRepo.DeleteReview(testID)
	require.NoError(t, err)
}

func TestReviewsRepository_GetRatingByRestID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	testID := uint64(1)

	rows := sqlmock.NewRows([]string{"avg"}).AddRow(float64(5))

	mock.ExpectQuery("SELECT avg").
		WithArgs(testID).WillReturnRows(rows)

	mockRepo := NewReviewsRepository(db)
	rate, err := mockRepo.GetRatingByRestID(testID)
	require.NoError(t, err)
	require.EqualValues(t, float64(5), rate)
}

func TestReviewsRepository_GetRestaurantReviewByUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	testID := uint64(1)

	testReview := &models.Review{
		ID:     testID,
		RestID: 1,
		Text:   "Good",
		Author: &models.User{
			ID:        1,
			FirstName: "abc",
			LastName:  "asd",
		},
		CreationDate: time.Now(),
		Rate:         5,
	}

	rows := sqlmock.NewRows([]string{"id", "restid", "id", "firstname",
		"lastname", "text", "creationdate", "rate"}).
		AddRow(testReview.ID, testReview.RestID, testReview.Author.ID, testReview.Author.FirstName,
			testReview.Author.LastName, testReview.Text, testReview.CreationDate, testReview.Rate)

	mock.ExpectQuery("SELECT").WithArgs(testReview.ID, testReview.Author.ID).WillReturnRows(rows)
	mockRepo := NewReviewsRepository(db)
	result, err := mockRepo.GetRestaurantReviewByUser(testReview.ID, testReview.Author.ID)
	require.NoError(t, err)
	require.EqualValues(t, testReview, result)
}

func TestReviewsRepository_GetReviewByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	testID := uint64(1)

	testReview := &models.Review{
		ID:     1,
		RestID: 1,
		Text:   "Good",
		Author: &models.User{
			ID:        1,
			FirstName: "abc",
			LastName:  "asd",
		},
		CreationDate: time.Now(),
		Rate:         5,
	}

	rows := sqlmock.NewRows([]string{"id", "restid", "id", "firstname",
		"lastname", "text", "creationdate", "rate"}).
		AddRow(testReview.ID, testReview.RestID, testReview.Author.ID, testReview.Author.FirstName,
			testReview.Author.LastName, testReview.Text, testReview.CreationDate, testReview.Rate)

	mock.ExpectQuery("SELECT").WithArgs(testID).WillReturnRows(rows)
	mockRepo := NewReviewsRepository(db)
	result, err := mockRepo.GetReviewByID(testID)
	require.NoError(t, err)
	require.EqualValues(t, testReview, result)
}

func TestReviewsRepository_GetReviewsByRestID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	testID := uint64(1)

	testReview := []*models.Review{
		{
			ID:           1,
			RestID:       1,
			Text:         "Good",
			Author:       &models.User{ID: testID},
			CreationDate: time.Now(),
			Rate:         5,
		},
	}

	rows := sqlmock.NewRows([]string{"id", "restid", "id", "firstname",
		"lastname", "text", "creationdate", "rate"}).
		AddRow(testReview[0].ID, testReview[0].RestID, testReview[0].Author.ID, testReview[0].Author.FirstName,
			testReview[0].Author.LastName, testReview[0].Text, testReview[0].CreationDate, testReview[0].Rate)

	mock.ExpectQuery("SELECT").
		WithArgs(testID, 1, 0).
		WillReturnRows(rows)

	mockRepo := NewReviewsRepository(db)
	result, err := mockRepo.GetReviewsByUserID(testID, 1, 1)
	require.NoError(t, err)
	require.EqualValues(t, testReview, result)
}

func TestReviewsRepository_GetReviewsByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	testID := uint64(1)

	testReview := []*models.Review{
		{
			ID:           1,
			RestID:       testID,
			Text:         "Good",
			Author:       &models.User{ID: 1},
			CreationDate: time.Now(),
			Rate:         5,
		},
	}

	rows := sqlmock.NewRows([]string{"id", "restid", "id", "firstname",
		"lastname", "text", "creationdate", "rate"}).
		AddRow(testReview[0].ID, testReview[0].RestID, testReview[0].Author.ID, testReview[0].Author.FirstName,
			testReview[0].Author.LastName, testReview[0].Text, testReview[0].CreationDate, testReview[0].Rate)

	mock.ExpectQuery("SELECT").
		WithArgs(testID, 1, 0).
		WillReturnRows(rows)

	mockRepo := NewReviewsRepository(db)
	result, err := mockRepo.GetReviewsByRestID(testID, 1, 1)
	require.NoError(t, err)
	require.EqualValues(t, testReview, result)
}

func TestReviewsRepository_GetReviewsCountByRestID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	testID := uint64(1)

	rows := sqlmock.NewRows([]string{"count"}).AddRow(uint64(5))

	mock.ExpectQuery("SELECT COUNT").
		WithArgs(testID).WillReturnRows(rows)

	mockRepo := NewReviewsRepository(db)
	rate, err := mockRepo.GetReviewsCountByUserID(testID)
	require.NoError(t, err)
	require.EqualValues(t, uint64(5), rate)
}

func TestReviewsRepository_GetReviewsCountByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	testID := uint64(1)

	rows := sqlmock.NewRows([]string{"count"}).AddRow(uint64(5))

	mock.ExpectQuery("SELECT COUNT").
		WithArgs(testID).WillReturnRows(rows)

	mockRepo := NewReviewsRepository(db)
	rate, err := mockRepo.GetReviewsCountByUserID(testID)
	require.NoError(t, err)
	require.EqualValues(t, uint64(5), rate)
}

func TestReviewsRepository_UpdateReview(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	testReview := &models.Review{
		ID:           1,
		RestID:       1,
		Text:         "Good",
		Author:       &models.User{ID: 1},
		CreationDate: time.Now(),
		Rate:         5,
	}

	mock.ExpectExec("UPDATE").
		WithArgs(testReview.ID, testReview.Text, testReview.Rate).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mockRepo := NewReviewsRepository(db)
	err = mockRepo.UpdateReview(testReview)
	require.NoError(t, err)
}
