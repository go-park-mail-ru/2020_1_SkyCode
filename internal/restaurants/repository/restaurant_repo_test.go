package repository

import (
	"testing"

	"github.com/2020_1_Skycode/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestRestaurantRepository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewRestaurantRepository(db)

	rows := sqlmock.
		NewRows([]string{"id", "name", "rating", "image"})
	expect := []*models.Restaurant{
		{uint64(1), 0, "rest1", "", 5.0, "./default.jpg", nil},
		{uint64(2), 0, "rest2", "", 4.4, "./not_default.jpg", nil},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.Name, item.Rating, item.Image)
	}

	mock.
		ExpectQuery("SELECT").
		WillReturnRows(rows)

	restList, err := repo.GetAll()
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}
	require.EqualValues(t, expect, restList)
}

func TestRestaurantRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewRestaurantRepository(db)

	rows := sqlmock.
		NewRows([]string{"id", "moderId", "name", "description", "rating", "image"})
	expect := &models.Restaurant{
		ID:          uint64(1),
		ManagerID:   uint64(1),
		Name:        "rest1",
		Description: "some description",
		Rating:      5.0,
		Image:       "./default.jpg",
		Products:    nil,
	}

	rows = rows.AddRow(expect.ID, expect.ManagerID, expect.Name, expect.Description, expect.Rating, expect.Image)

	var elemID uint64 = 1

	mock.
		ExpectQuery("SELECT").
		WithArgs(elemID).
		WillReturnRows(rows)

	rest, err := repo.GetByID(elemID)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}

	require.EqualValues(t, expect, rest)
}

func TestRestaurantRepository_InsertInto(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewRestaurantRepository(db)

	testRest := &models.Restaurant{
		ManagerID:   1,
		Name:        "test",
		Description: "test restaurant",
		Rating:      5.0,
		Image:       "default.jpg",
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectQuery("INSERT INTO restaurants").
		WithArgs(testRest.ManagerID, testRest.Name, testRest.Description, testRest.Rating, testRest.Image).
		WillReturnRows(rows)

	err = repo.InsertInto(testRest)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}

	require.EqualValues(t, 1, testRest.ID)
}

func TestRestaurantRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewRestaurantRepository(db)

	testRest := &models.Restaurant{
		ID:          1,
		ManagerID:   1,
		Name:        "test",
		Description: "test restaurant",
	}

	mock.ExpectExec("UPDATE restaurants SET").
		WithArgs(testRest.ID, testRest.Name, testRest.Description).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(testRest)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}
}

func TestRestaurantRepository_UpdateImage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewRestaurantRepository(db)

	testRest := &models.Restaurant{
		ID:    1,
		Image: "default.jpg",
	}

	mock.ExpectExec("UPDATE restaurants SET").
		WithArgs(testRest.ID, testRest.Image).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateImage(testRest)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}
}

func TestRestaurantRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewRestaurantRepository(db)

	prodID := uint64(1)

	mock.ExpectExec("DELETE FROM restaurants WHERE").
		WithArgs(prodID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(prodID)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}
}
