package repository

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/restaurants/repository"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestOrdersRepository_InsertOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	restRepo := repository.NewRestaurantRepository(db)
	repo := NewOrdersRepository(db, restRepo)

	testOrder := &models.Order{
		UserID:    1,
		RestID:    1,
		Address:   "test address",
		Phone:     "79986543321",
		Comment:   "faster",
		PersonNum: 3,
		Price:     555,
	}

	testProducts := []*models.OrderProduct{
		&models.OrderProduct{
			ProductID: 1,
			Count:     2,
		},
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectQuery("INSERT INTO orders").
		WithArgs(testOrder.UserID, testOrder.RestID, testOrder.Address, testOrder.Comment,
			testOrder.PersonNum, testOrder.Phone, testOrder.Price).
		WillReturnRows(rows)

	mock.ExpectExec("INSERT INTO orderProducts").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.InsertOrder(testOrder, testProducts)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}

	require.EqualValues(t, 1, testOrder.ID)
}

func TestOrdersRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	restRepo := repository.NewRestaurantRepository(db)
	repo := NewOrdersRepository(db, restRepo)

	testOrderProduct := []*models.Product{
		&models.Product{
			ID:     1,
			Name:   "test product",
			Price:  322,
			Image:  "default.jpg",
			RestId: 1,
			Count:  2,
		},
	}

	testOrder := &models.Order{
		ID:        1,
		Address:   "test address",
		Phone:     "79986543321",
		Comment:   "faster",
		Products:  testOrderProduct,
		CreatedAt: time.Now().String(),
		PersonNum: 3,
		Price:     555,
	}

	userID := uint64(1)

	rows := sqlmock.NewRows([]string{"id", "address", "phone", "price", "comment", "personNum", "datetime"})
	rows.AddRow(testOrder.ID, testOrder.Address, testOrder.Phone,
		testOrder.Price, testOrder.Comment, testOrder.PersonNum, testOrder.CreatedAt)

	rowsProd := sqlmock.NewRows([]string{"id", "restid", "name", "price", "image", "count"})
	rowsProd.AddRow(testOrderProduct[0].ID, testOrderProduct[0].RestId, testOrderProduct[0].Name,
		testOrderProduct[0].Price, testOrderProduct[0].Image, testOrderProduct[0].Count)

	mock.ExpectQuery("SELECT id, address").
		WithArgs(testOrder.ID, userID).
		WillReturnRows(rows)

	mock.ExpectQuery("select p.id, p.rest_id, p.name, p.price, p.image, orderProducts.count").
		WithArgs(1).
		WillReturnRows(rowsProd)

	resultOrder, err := repo.GetByID(testOrder.ID, userID)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}

	require.EqualValues(t, testOrder, resultOrder)
}

func TestOrdersRepository_GetAllByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	restRepo := repository.NewRestaurantRepository(db)
	repo := NewOrdersRepository(db, restRepo)

	testOrderProduct := []*models.Product{
		&models.Product{
			ID:     1,
			Name:   "test product",
			Price:  322,
			Image:  "default.jpg",
			RestId: 1,
			Count:  2,
		},
	}

	tN := time.Now()

	testOrder := []*models.Order{
		{
			ID:        1,
			UserID:    1,
			RestID:    1,
			RestName:  "test",
			Phone:     "89765433221",
			Address:   "test address",
			Comment:   "faster",
			PersonNum: 3,
			Price:     555.0,
			Products:  testOrderProduct,
			CreatedAt: tN.Format("2006/Jan/_2/15:04:05"),
			Status:    "Accepted",
		},
	}

	rows := sqlmock.NewRows([]string{"id", "userid", "restid", "address", "price", "phone", "comment", "personNum",
		"datetime", "status"})

	for _, item := range testOrder {
		rows = rows.AddRow(item.ID, item.UserID, item.RestID, item.Address, item.Price, item.Phone,
			item.Comment, item.PersonNum, tN, item.Status)
	}

	rowsProd := sqlmock.NewRows([]string{"id", "restid", "name", "price", "image", "count"})
	rowsProd.AddRow(testOrderProduct[0].ID, testOrderProduct[0].RestId, testOrderProduct[0].Name,
		testOrderProduct[0].Price, testOrderProduct[0].Image, testOrderProduct[0].Count)

	rowsCount := sqlmock.NewRows([]string{"count"}).AddRow(1)

	expectRest := &models.Restaurant{
		ID:          uint64(1),
		ManagerID:   uint64(1),
		Name:        "test",
		Description: "some description",
		Rating:      5.0,
		Image:       "./default.jpg",
		Products:    nil,
	}

	rowRest := sqlmock.
		NewRows([]string{"id", "moderId", "name", "description", "rating", "image"}).
		AddRow(expectRest.ID, expectRest.ManagerID, expectRest.Name, expectRest.Description,
			expectRest.Rating, expectRest.Image)

	mock.ExpectQuery("select id, userId").
		WithArgs(1, 1, 0).
		WillReturnRows(rows)

	mock.ExpectQuery("SELECT COUNT").
		WithArgs(uint64(1)).
		WillReturnRows(rowsCount)

	mock.ExpectQuery("select p.id, p.rest_id, p.name, p.price, p.image, orderProducts.count").
		WithArgs(1).
		WillReturnRows(rowsProd)

	mock.
		ExpectQuery("SELECT id, moderId").
		WithArgs(1).
		WillReturnRows(rowRest)

	resultOrder, total, err := repo.GetAllByUserID(1, 1, 1)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}

	require.EqualValues(t, testOrder, resultOrder)
	require.EqualValues(t, 1, total)
}

func TestOrdersRepository_DeleteOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	restRepo := repository.NewRestaurantRepository(db)
	repo := NewOrdersRepository(db, restRepo)

	orderID := uint64(1)
	userID := uint64(1)

	rows := sqlmock.NewRows([]string{"id"})
	rows.AddRow(1)

	mock.ExpectQuery("DELETE").
		WithArgs(orderID, userID).WillReturnRows(rows)

	err = repo.DeleteOrder(orderID, userID)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There was unfulfilled expectations %s", err)
		return
	}
}
