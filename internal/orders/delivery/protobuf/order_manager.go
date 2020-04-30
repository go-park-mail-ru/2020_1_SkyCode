package protobuf_order

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/orders"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/2020_1_Skycode/tools/protobuf/orderswork"
	"golang.org/x/net/context"
)

type OrderManager struct {
	orderRepo orders.Repository
}

func NewOrderProtoManager(orderRepo orders.Repository) *OrderManager {
	return &OrderManager{
		orderRepo: orderRepo,
	}
}

func (oU *OrderManager) CheckOutOrder(ctx context.Context, c *orderswork.Checkout) (*orderswork.Error, error) {
	order := &models.Order{
		ID:        c.Order.ID,
		UserID:    c.Order.UserID,
		RestID:    c.Order.RestID,
		RestName:  c.Order.RestName,
		Address:   c.Order.Address,
		Phone:     c.Order.Phone,
		Comment:   c.Order.Comment,
		PersonNum: c.Order.PersonNum,
		Price:     c.Order.Price,
	}

	products := []*models.OrderProduct{}

	for _, val := range c.Products {
		product := &models.OrderProduct{
			ID:        val.ID,
			OrderID:   val.OrderID,
			ProductID: val.ProductID,
			Count:     val.Count,
		}

		products = append(products, product)
	}

	if err := oU.orderRepo.InsertOrder(order, products); err != nil {
		return &orderswork.Error{Err: tools.CheckoutOrderError.Error()}, err
	}

	return &orderswork.Error{}, nil
}

func (oU *OrderManager) GetAllUserOrders(ctx context.Context, u *orderswork.UserOrders) (*orderswork.GetAllResponse, error) {
	userOrders, total, err := oU.orderRepo.GetAllByUserID(u.UserID, u.Count, u.Page)

	orders := []*orderswork.Order{}

	for _, val := range userOrders {
		products := []*orderswork.Product{}

		for _, val := range val.Products {
			products = append(products, &orderswork.Product{
				ID:     val.ID,
				Name:   val.Name,
				Price:  val.Price,
				Image:  val.Image,
				RestID: val.RestId,
				Count:  val.Count,
			})
		}

		orders = append(orders, &orderswork.Order{
			ID:        val.ID,
			UserID:    val.UserID,
			RestID:    val.RestID,
			RestName:  val.RestName,
			Address:   val.Address,
			Phone:     val.Phone,
			Comment:   val.Comment,
			PersonNum: val.PersonNum,
			Products:  products,
			Price:     val.Price,
			CreatedAt: val.CreatedAt,
			Status:    val.Status,
		})
	}

	res := &orderswork.GetAllResponse{
		Order: orders,
		Total: total,
	}

	if err != nil {
		return res, err
	}

	return res, nil
}

func (oU *OrderManager) GetOrderByID(ctx context.Context, u *orderswork.GetByID) (*orderswork.GetByIDResponse, error) {
	order, err := oU.orderRepo.GetByID(u.OrderID, u.UserID)

	if err != nil {
		return &orderswork.GetByIDResponse{}, err
	}

	products := []*orderswork.Product{}

	for _, val := range order.Products {
		products = append(products, &orderswork.Product{
			ID:     val.ID,
			Name:   val.Name,
			Price:  val.Price,
			Image:  val.Image,
			RestID: val.RestId,
			Count:  val.Count,
		})
	}

	resOrder := &orderswork.Order{
		ID:        order.ID,
		UserID:    order.UserID,
		RestID:    order.RestID,
		RestName:  order.RestName,
		Address:   order.Address,
		Phone:     order.Comment,
		Comment:   order.Comment,
		PersonNum: order.PersonNum,
		Products:  products,
		Price:     order.Price,
		CreatedAt: order.CreatedAt,
		Status:    order.Status,
	}

	res := &orderswork.GetByIDResponse{
		Order: resOrder,
	}

	return res, err
}

func (oU *OrderManager) DeleteOrder(ctx context.Context, d *orderswork.DelOrder) (*orderswork.Error, error) {
	if err := oU.orderRepo.DeleteOrder(d.OrderID, d.UserID); err != nil {
		return &orderswork.Error{Err: err.Error()}, err
	}

	return &orderswork.Error{}, nil
}
