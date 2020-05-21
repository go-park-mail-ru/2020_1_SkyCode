package repository

import (
	"database/sql"
	"fmt"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/orders"
	"github.com/2020_1_Skycode/internal/restaurants"
	"time"
)

type OrdersRepository struct {
	RestRepository restaurants.Repository
	db             *sql.DB
}

func NewOrdersRepository(db *sql.DB, rR restaurants.Repository) orders.Repository {
	return &OrdersRepository{
		RestRepository: rR,
		db:             db,
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

	rows, err := oR.db.Query("select id, userId, restId, address, price, phone, comment, personnum, "+
		"datetime, status from orders where userId = $1"+
		" LIMIT $2 OFFSET $3",
		userID, count, (page-1)*count)
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
		time := time.Time{}

		err = rows.Scan(&order.ID, &order.UserID, &order.RestID, &order.Address, &order.Price, &order.Phone, &order.Comment, &order.PersonNum, &time, &order.Status)

		order.CreatedAt = time.Format("2006/Jan/_2/15:04:05")

		if err != nil {
			return nil, 0, err
		}

		products, err := oR.getOrderProducts(order.ID)

		if err != nil {
			return nil, 0, err
		}

		restaurant, err := oR.RestRepository.GetByID(order.RestID)

		if err != nil {
			return nil, 0, err
		}

		order.RestName = restaurant.Name

		order.Products = products
		ordersList = append(ordersList, order)
	}

	return ordersList, total, nil
}

func (oR *OrdersRepository) GetByID(orderID uint64) (*models.Order, error) {
	order := &models.Order{}
	err := oR.db.QueryRow("SELECT id, userid, address, phone, price, comment, personnum, datetime, status "+
		"FROM orders WHERE id = $1",
		orderID).Scan(&order.ID, &order.UserID, &order.Address, &order.Phone, &order.Price, &order.Comment,
		&order.PersonNum, &order.CreatedAt, &order.Status)

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

func (oR *OrdersRepository) ChangeStatus(orderID uint64, status string) error {
	if _, err := oR.db.Exec("UPDATE orders SET status = $2 WHERE id = $1", orderID, status); err != nil {
		return err
	}

	return nil
}

func (oR *OrdersRepository) getOrderProducts(orderID uint64) ([]*models.Product, error) {
	var ProductsList []*models.Product

	rows, err := oR.db.Query("select p.id, p.rest_id, p.name, p.price, p.image, orderProducts.count "+
		"from products p join (select productId, count from orderproducts where orderId = $1) as orderProducts "+
		"on p.id = orderProducts.productId;",
		orderID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		order := &models.Product{}
		err = rows.Scan(&order.ID, &order.RestId, &order.Name, &order.Price, &order.Image, &order.Count)

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

func (oR *OrdersRepository) GetAllByRestID(restID uint64, count uint64, page uint64) ([]*models.Order, uint64, error) {
	rows, err := oR.db.Query("SELECT id, userId, restId, address, price, phone, comment, "+
		"datetime, status, "+
		"CASE WHEN status = 'Accepted' THEN 1 "+
		"WHEN status = 'Delivering' THEN 2 "+
		"WHEN status = 'Canceled' THEN 3 "+
		"WHEN status = 'Done' THEN 4 "+
		"end as id_status "+
		"FROM orders WHERE restid = $1 ORDER BY id_status, datetime "+
		"LIMIT $2 OFFSET $3", restID, count, count*(page-1))
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	var ordersList []*models.Order
	for rows.Next() {
		o := &models.Order{}

		var statusId uint64
		if err := rows.Scan(&o.ID, &o.UserID, &o.RestID, &o.Address, &o.Price, &o.Phone, &o.Comment, &o.CreatedAt,
			&o.Status, &statusId); err != nil {
			return nil, 0, err
		}

		ordersList = append(ordersList, o)
	}

	var total uint64
	if err := oR.db.QueryRow("SELECT COUNT(*) FROM orders WHERE restid = $1", restID).
		Scan(&total); err != nil {
		return nil, 0, err
	}

	return ordersList, total, nil
}
