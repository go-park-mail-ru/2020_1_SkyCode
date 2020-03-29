package repository

import (
	"fmt"
	"testing"

	"github.com/2020_1_Skycode/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestGetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewRestaurantRepository(db)

	rows := sqlmock.
		NewRows([]string{"id", "name", "rating", "image"})
	expect := []*models.Restaurant{
		{uint64(1), "rest1", "", 5.0, "./default.jpg", nil},
		{uint64(2), "rest2", "", 4.4, "./not_default.jpg", nil},
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

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewRestaurantRepository(db)

	rows := sqlmock.
		NewRows([]string{"id", "name", "description", "rating", "image"})
	expect := &models.Restaurant{
		ID:          uint64(1),
		Name:        "rest1",
		Description: "some description",
		Rating:      5.0,
		Image:       "./default.jpg",
		Products:    nil,
	}

	rows = rows.AddRow(expect.ID, expect.Name, expect.Description, expect.Rating, expect.Image)

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
	fmt.Println(rest.ID, rest.Name)
	require.EqualValues(t, expect, rest)
}
