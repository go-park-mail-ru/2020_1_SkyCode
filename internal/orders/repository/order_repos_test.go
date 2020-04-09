package repository

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOrdersRepository_InsertOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewOrdersRepository(db)

	testOrder := &models.Order{
		UserID:    1,
		Address:   "test address",
		Phone:     "79986543321",
		Comment:   "faster",
		PersonNum: 3,
		Products: []*models.OrderProduct{
			{
				ProductID: 1,
				Count:     2,
			},
		},
		Price: 555,
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectQuery("INSERT INTO orders").
		WithArgs(testOrder.UserID, testOrder.Address, testOrder.Comment,
			testOrder.PersonNum, testOrder.Phone, testOrder.Price).
		WillReturnRows(rows)

	mock.ExpectExec("INSERT INTO orderProducts").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.InsertOrder(testOrder)
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

	repo := NewOrdersRepository(db)

	testOrderProduct := &models.OrderProduct{
		ID:        1,
		OrderID:   1,
		ProductID: 1,
		Count:     2,
	}

	testOrder := &models.Order{
		ID:        1,
		Phone:     "89765433221",
		Address:   "test address",
		Comment:   "faster",
		Products:  []*models.OrderProduct{testOrderProduct},
		PersonNum: 3,
		Price:     555,
	}

	userID := uint64(1)

	rows := sqlmock.NewRows([]string{"id", "address", "phone", "price", "comment", "personNum"})
	rows.AddRow(testOrder.ID, testOrder.Address, testOrder.Phone,
		testOrder.Price, testOrder.Comment, testOrder.PersonNum)

	rowsProd := sqlmock.NewRows([]string{"id", "orderId", "productId", "count"})
	rowsProd.AddRow(testOrderProduct.ID, testOrderProduct.OrderID, testOrderProduct.ProductID, testOrderProduct.Count)

	mock.ExpectQuery("SELECT").
		WithArgs(testOrder.ID, userID).
		WillReturnRows(rows)

	mock.ExpectQuery("SELECT id, orderid, productid, count FROM orderproducts").
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

	repo := NewOrdersRepository(db)

	testOrder := []*models.Order{
		{
			ID:        1,
			Phone:     "89765433221",
			Address:   "test address",
			Comment:   "faster",
			PersonNum: 3,
			Price:     555.0,
		},
	}

	rows := sqlmock.NewRows([]string{"id", "address", "phone", "price", "comment", "personNum"})

	for _, item := range testOrder {
		rows = rows.AddRow(item.ID, item.Address, item.Phone,
			item.Price, item.Comment, item.PersonNum)
	}

	mock.ExpectQuery("SELECT").
		WithArgs(1).
		WillReturnRows(rows)

	resultOrder, err := repo.GetAllByUserID(1)
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

func TestOrdersRepository_DeleteOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewOrdersRepository(db)

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
