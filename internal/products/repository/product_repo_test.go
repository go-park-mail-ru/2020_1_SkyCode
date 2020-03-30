package repository

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetRestaurantsProducts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewProductRepository(db)

	rows := sqlmock.
		NewRows([]string{"id", "name", "price", "image"})
	expect := []*models.Product{
		{ID: 1, Name: "test", Price: 2.50, Image: "./default_img.jpg"},
	}

	var restId = uint64(1)

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.Name, item.Price, item.Image)
	}

	mock.
		ExpectQuery("SELECT id, name, price, image FROM products WHERE").
		WithArgs(restId).
		WillReturnRows(rows)

	prodList, err := repo.GetRestaurantProducts(restId)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}
	require.EqualValues(t, expect, prodList)
}
