package repository

import (
	"database/sql"
	"fmt"
	"github.com/2020_1_Skycode/internal/models"
)

type OrdersRepository struct {
	db *sql.DB
}

func NewOrdersRepository(db *sql.DB) *OrdersRepository {
	return &OrdersRepository{
		db: db,
	}
}

func (oR *OrdersRepository) InsertOrder(order *models.Order) error {
	if err := oR.db.QueryRow("INSERT INTO orders(userId, address, comment, personNum, phone, price) "+
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		order.UserID,
		order.Address,
		order.Comment,
		order.PersonNum,
		order.Phone,
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
	sqlInsert := "INSERT INTO orderProducts (orderId, productId, count) VALUES"

	for _, v := range products {
		values += fmt.Sprintf(" (%d, %d, %d),", orderID, v.ProductID, v.Count)
	}

	array := []rune(values)
	array[len(array)-1] = ';'
	values = string(array)

	if _, err := oR.db.Exec(sqlInsert + values); err != nil {
		return err
	}

	return nil
}

func (oR *OrdersRepository) GetAllByUserID(userID uint64, count uint64, page uint64) ([]*models.Order, uint64, error) {
	var ordersList []*models.Order

	rows, err := oR.db.Query("SELECT id, address, phone, price, comment, personnum FROM orders WHERE userId = $1 "+
		"LIMIT $2 OFFSET $3", userID, count, (page-1)*count)
	if err != nil {
		return nil, 0, err
	}

	var total uint64
	if err = oR.db.QueryRow("SELECT COUNT(*) FROM orders WHERE userId = $1", userID).
		Scan(&total); err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		order := &models.Order{}
		err = rows.Scan(&order.ID, &order.Address, &order.Phone, &order.Price, &order.Comment, &order.PersonNum)

		if err != nil {
			return nil, 0, err
		}

		ordersList = append(ordersList, order)
	}

	return ordersList, total, nil
}

func (oR *OrdersRepository) GetByID(orderID uint64, userID uint64) (*models.Order, error) {
	order := &models.Order{}
	err := oR.db.QueryRow("SELECT id, address, phone, price, comment, personnum FROM orders WHERE id = $1 AND userid = $2",
		orderID,
		userID).Scan(&order.ID, &order.Address, &order.Phone, &order.Price, &order.Comment, &order.PersonNum)

	if err != nil {
		return nil, err
	}

	products, err := oR.getOrderProducts(order.ID)

	if err != nil {
		return nil, err
	}

	order.Products = products

	return order, nil
}

func (oR *OrdersRepository) getOrderProducts(orderID uint64) ([]*models.OrderProduct, error) {
	var ordersProductList []*models.OrderProduct

	rows, err := oR.db.Query("SELECT id, orderid, productid, count FROM orderproducts WHERE orderid = $1",
		orderID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		order := &models.OrderProduct{}
		err = rows.Scan(&order.ID, &order.OrderID, &order.ProductID, &order.Count)

		if err != nil {
			return nil, err
		}

		ordersProductList = append(ordersProductList, order)
	}

	return ordersProductList, nil
}

func (oR *OrdersRepository) DeleteOrder(orderID uint64, userID uint64) error {
	var id uint64
	if err := oR.db.QueryRow("DELETE FROM orders CASCADE WHERE id = $1 AND userId = $2 RETURNING id",
		orderID,
		userID).Scan(&id); err != nil {
		return err
	}

	return nil
}
