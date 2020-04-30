package repository

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRestPointsRepository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	restPointsRepo := NewRestPosintsRepository(db)

	expect := []*models.RestaurantPoint{
		{
			ID:      1,
			Address: "Pushkina dom Kolotushkina",
			MapPoint: &models.GeoPos{
				Longitude: 55.753227,
				Latitude:  37.619030,
			},
			ServiceRadius: 5,
			RestID:        1,
		},
	}

	rows := sqlmock.NewRows([]string{"id", "address", "latitude", "longitude", "restid", "radius"}).
		AddRow(expect[0].ID, expect[0].Address, expect[0].MapPoint.Latitude,
			expect[0].MapPoint.Longitude, expect[0].RestID, expect[0].ServiceRadius)

	mock.ExpectQuery("SELECT").WithArgs().WillReturnRows(rows)

	result, err := restPointsRepo.GetAll()
	require.NoError(t, err)

	require.EqualValues(t, expect, result)
}

func TestRestPointsRepository_GetCloserPointByRestID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	restPointsRepo := NewRestPosintsRepository(db)

	point := &models.GeoPos{
		Longitude: 55.753227,
		Latitude:  37.619030,
	}

	testRestID := uint64(1)

	expect := &models.RestaurantPoint{
		ID:            1,
		Address:       "Pushkina dom Kolotushkina",
		MapPoint:      point,
		ServiceRadius: 5,
		RestID:        testRestID,
	}

	rows := sqlmock.NewRows([]string{"id", "restid", "latitude", "longitude", "address", "radius"}).
		AddRow(expect.ID, expect.RestID, expect.MapPoint.Latitude,
			expect.MapPoint.Longitude, expect.Address, expect.ServiceRadius)

	mock.ExpectQuery("SELECT").
		WithArgs(testRestID, point.Latitude, point.Longitude).
		WillReturnRows(rows)

	result, err := restPointsRepo.GetCloserPointByRestID(testRestID, point)
	require.NoError(t, err)

	require.EqualValues(t, expect, result)
}

func TestRestPointsRepository_GetPointByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	restPointsRepo := NewRestPosintsRepository(db)

	testID := uint64(1)

	expect := &models.RestaurantPoint{
		ID:      testID,
		Address: "Pushkina dom Kolotushkina",
		MapPoint: &models.GeoPos{
			Longitude: 55.753227,
			Latitude:  37.619030,
		},
		ServiceRadius: 5,
		RestID:        1,
	}

	rows := sqlmock.NewRows([]string{"id", "address", "latitude", "longitude", "restid", "radius"}).
		AddRow(expect.ID, expect.Address, expect.MapPoint.Latitude,
			expect.MapPoint.Longitude, expect.RestID, expect.ServiceRadius)

	mock.ExpectQuery("SELECT").WithArgs().WillReturnRows(rows)

	result, err := restPointsRepo.GetPointByID(testID)
	require.NoError(t, err)

	require.EqualValues(t, expect, result)
}

func TestRestPointsRepository_GetPointsByRestID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	restPointsRepo := NewRestPosintsRepository(db)

	testRestID := uint64(1)

	expect := []*models.RestaurantPoint{
		{
			ID:      1,
			Address: "Pushkina dom Kolotushkina",
			MapPoint: &models.GeoPos{
				Longitude: 55.753227,
				Latitude:  37.619030,
			},
			ServiceRadius: 5,
			RestID:        1,
		},
	}

	rows := sqlmock.NewRows([]string{"id", "address", "latitude", "longitude", "restid", "radius"}).
		AddRow(expect[0].ID, expect[0].Address, expect[0].MapPoint.Latitude,
			expect[0].MapPoint.Longitude, expect[0].RestID, expect[0].ServiceRadius)

	rowsCount := sqlmock.NewRows([]string{"count"}).AddRow(1)

	mock.ExpectQuery("SELECT").WithArgs(testRestID, 1, 0).WillReturnRows(rows)
	mock.ExpectQuery("SELECT count").WithArgs(testRestID).WillReturnRows(rowsCount)

	result, total, err := restPointsRepo.GetPointsByRestID(testRestID, 1, 1)
	require.NoError(t, err)

	require.EqualValues(t, expect, result)
	require.EqualValues(t, 1, total)
}

func TestRestPointsRepository_InsertInto(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	restPointsRepo := NewRestPosintsRepository(db)

	expect := &models.RestaurantPoint{
		ID:      1,
		Address: "Pushkina dom Kolotushkina",
		MapPoint: &models.GeoPos{
			Longitude: 55.753227,
			Latitude:  37.619030,
		},
		ServiceRadius: 5,
		RestID:        1,
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(expect.ID)

	mock.ExpectQuery("INSERT INTO").WithArgs(expect.Address, expect.MapPoint.Latitude,
		expect.MapPoint.Longitude, expect.RestID, expect.ServiceRadius).WillReturnRows(rows)

	err = restPointsRepo.InsertInto(expect)
	require.NoError(t, err)
}

func TestRestPointsRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	restPointsRepo := NewRestPosintsRepository(db)

	id := uint64(1)

	mock.ExpectExec("DELETE").WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = restPointsRepo.Delete(id)
	require.NoError(t, err)
}
