package repository

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/2020_1_Skycode/internal/models"
)

type OrdersRepository struct {
	db *pgx.Conn
}

func NewOrdersRepository(db *pgx.Conn) *OrdersRepository{
	return &OrdersRepository{
		db: db,
	}
}

func (oR *OrdersRepository) InsertOrder(order *models.Order) error {
	if err := oR.db.QueryRow("INSERT INTO orders(userId, address, comment, personNum, price) "+
		"VALUES ($1, $2, $3, $4, $5) RETURNING id",
		order.UserID,
		order.Address,
		order.Comment,
		order.PersonNum,
		order.Price).Scan(&order.ID); err != nil {
		return err
	}

	if err := oR.insertOrderProducts(order.ID, order.Products); err != nil {
		return err
	}

	return nil
}

func (oR *OrdersRepository) insertOrderProducts(orderID uint64, products []*models.OrderProduct) error {
	var values string
	sqlInsert := "INSERT INTO orders(orderId, productId, count) VALUES"

	for _, v := range products {
		values += fmt.Sprintf(" (%d, %d, %d)", v.ID, v.ProductID, v.Count)
	}

	if _, err := oR.db.Exec(sqlInsert + values); err != nil {
		return err
	}

	return nil
}

func (oR *OrdersRepository) Get(order *models.Order) error {
	if err := oR.db.QueryRow("SELECT userId, address, comment, personNum, price FROM orders WHERE id = $1",
		order.ID).Scan(&order.UserID, &order.Address, &order.Comment, &order.PersonNum, &order.Price); err != nil {
		return err
	}

	return nil
}
