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
		{
			ID:     uint64(1),
			Name:   "rest1",
			Rating: 4.2,
			Image:  "./default.jpg",
		},
		{
			ID:     uint64(2),
			Name:   "rest2",
			Rating: 4.4,
			Image:  "./not_default.jpg",
		},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.Name, item.Rating, item.Image)
	}

	rowsCount := sqlmock.NewRows([]string{"count"}).AddRow(2)

	mock.
		ExpectQuery("SELECT id").
		WithArgs(uint64(2), uint64(0)).
		WillReturnRows(rows)

	mock.ExpectQuery("SELECT COUNT").WillReturnRows(rowsCount)

	restList, total, err := repo.GetAll(uint64(2), uint64(1))
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}
	require.EqualValues(t, expect, restList)
	require.EqualValues(t, uint64(2), total)
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

func TestRestaurantRepository_GetAllInServiceRadius(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewRestaurantRepository(db)

	rows := sqlmock.
		NewRows([]string{"id", "name", "description", "rating", "image", "dst"})
	expect := []*models.Restaurant{
		{
			ID:          uint64(1),
			Name:        "rest1",
			Description: "",
			Rating:      4.2,
			Image:       "./default.jpg",
		},
		{
			ID:          uint64(2),
			Name:        "rest2",
			Description: "I saw some shit",
			Rating:      4.4,
			Image:       "./not_default.jpg",
		},
	}

	gp := &models.GeoPos{
		Longitude: 55.753227,
		Latitude:  37.619030,
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.Name, item.Description, item.Rating, item.Image, 5)
	}

	rowsCount := sqlmock.NewRows([]string{"count"}).AddRow(2)

	mock.
		ExpectQuery("SELECT r.id").
		WithArgs(gp.Latitude, gp.Longitude, uint64(2), uint64(0)).
		WillReturnRows(rows)

	mock.ExpectQuery("SELECT COUNT").
		WithArgs(gp.Latitude, gp.Longitude).
		WillReturnRows(rowsCount)

	restList, total, err := repo.GetAllInServiceRadius(gp, uint64(2), uint64(1))
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}
	require.EqualValues(t, expect, restList)
	require.EqualValues(t, uint64(2), total)
}
