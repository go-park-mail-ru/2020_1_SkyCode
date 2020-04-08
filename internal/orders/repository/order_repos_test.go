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

func TestOrdersRepository_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Can't create mock: %s", err)
	}
	defer db.Close()

	repo := NewOrdersRepository(db)

	testOrder := &models.Order{
		ID:        1,
		UserID:    1,
		Address:   "test address",
		Comment:   "faster",
		PersonNum: 3,
		Price:     555,
	}

	rows := sqlmock.NewRows([]string{"userID", "address", "comment", "personNum", "price"})
	rows.AddRow(testOrder.UserID, testOrder.Address, testOrder.Comment,
		testOrder.PersonNum, testOrder.Price)

	mock.ExpectQuery("SELECT").
		WithArgs(testOrder.ID).
		WillReturnRows(rows)

	resultOrder := &models.Order{ID: 1}

	err = repo.Get(resultOrder)
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
