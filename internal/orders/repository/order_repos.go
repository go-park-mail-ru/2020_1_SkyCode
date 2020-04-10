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

func (oR *OrdersRepository) InsertOrder(order *models.Order, ordProducts []*models.OrderProduct) error {
	if err := oR.db.QueryRow("INSERT INTO orders(userId, restId, address, comment, personNum, phone, price) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		order.UserID,
		order.RestID,
		order.Address,
		order.Comment,
		order.PersonNum,
		order.Phone,
		order.Price).Scan(&order.ID); err != nil {
		return err
	}

	if err := oR.insertOrderProducts(order.ID, ordProducts); err != nil {
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

	rows, err := oR.db.Query("select id, userId, restId, address, price, phone, comment, personnum, datetime from orders where userId = $1" +
		" LIMIT $2 OFFSET $3",
		userID, count, (page - 1) * count)
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
		err = rows.Scan(&order.ID, &order.UserID, &order.RestID, &order.Address, &order.Price, &order.Phone, &order.Comment, &order.PersonNum, &order.CreatedAt)

		if err != nil {
			return nil, 0, err
		}

		products, err := oR.getOrderProducts(order.ID)

		if err != nil {
			return nil, 0, err
		}

		order.Products = products
		ordersList = append(ordersList, order)
	}

	return ordersList, total, nil
}

func (oR *OrdersRepository) GetByID(orderID uint64, userID uint64) (*models.Order, error) {
	order := &models.Order{}
	err := oR.db.QueryRow("SELECT id, address, phone, price, comment, personnum, datetime FROM orders WHERE id = $1 AND userid = $2",
		orderID,
		userID).Scan(&order.ID, &order.Address, &order.Phone, &order.Price, &order.Comment, &order.PersonNum, &order.CreatedAt)

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

func (oR *OrdersRepository) getOrderProducts(orderID uint64) ([]*models.Product, error) {
	var ProductsList []*models.Product

	rows, err := oR.db.Query("select id, rest_id, name, price, image from products where id in (select id from orderproducts where orderId = $1)",
		orderID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		order := &models.Product{}
		err = rows.Scan(&order.ID, &order.RestId, &order.Name, &order.Price, &order.Image)

		if err != nil {
			return nil, err
		}

		ProductsList = append(ProductsList, order)
	}

	return ProductsList, nil
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
