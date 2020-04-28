package repository

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProductRepository_GetProductsByRestID(t *testing.T) {
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

	rowsCount := sqlmock.NewRows([]string{"count"}).AddRow(1)

	mock.
		ExpectQuery("SELECT id, name, price, image FROM products WHERE").
		WithArgs(restId, uint64(1), uint64(0)).
		WillReturnRows(rows)

	mock.ExpectQuery("SELECT COUNT").
		WithArgs(uint64(1)).
		WillReturnRows(rowsCount)

	prodList, total, err := repo.GetProductsByRestID(restId, uint64(1), uint64(1))
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}
	require.EqualValues(t, expect, prodList)
	require.EqualValues(t, uint64(1), total)
}

func TestProductRepository_GetProductByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewProductRepository(db)

	rows := sqlmock.
		NewRows([]string{"id", "name", "price", "image", "rest_id"})
	expect := &models.Product{ID: 1, Name: "test", Price: 2.50, Image: "./default_img.jpg", RestId: 1}

	var prodID = uint64(1)

	rows = rows.AddRow(expect.ID, expect.Name, expect.Price, expect.Image, expect.RestId)

	mock.
		ExpectQuery("SELECT id, name, price, image, rest_id FROM products WHERE").
		WithArgs(prodID).
		WillReturnRows(rows)

	prod, err := repo.GetProductByID(prodID)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}
	require.EqualValues(t, expect, prod)
}

func TestProductRepository_InsertInto(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewProductRepository(db)

	testProd := &models.Product{
		Name:   "test",
		Price:  2.50,
		Image:  "/default.jpg",
		RestId: 1,
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectQuery("INSERT INTO products").
		WithArgs(testProd.Name, testProd.Price, testProd.Image, testProd.RestId).
		WillReturnRows(rows)

	err = repo.InsertInto(testProd)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}

	require.EqualValues(t, 1, testProd.ID)
}

func TestProductRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewProductRepository(db)

	testProd := &models.Product{
		ID:    1,
		Name:  "test",
		Price: 2.50,
	}

	mock.ExpectExec("UPDATE products SET").
		WithArgs(testProd.ID, testProd.Name, testProd.Price).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(testProd)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}
}

func TestProductRepository_UpdateImage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewProductRepository(db)

	testProd := &models.Product{
		ID:    1,
		Image: "default.jpg",
	}

	mock.ExpectExec("UPDATE products SET").
		WithArgs(testProd.ID, testProd.Image).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateImage(testProd)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}
}

func TestProductRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewProductRepository(db)

	prodID := uint64(1)

	mock.ExpectExec("DELETE FROM products WHERE").
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
