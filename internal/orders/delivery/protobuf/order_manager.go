package protobuf_order

import (
	"database/sql"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/notifications"
	"github.com/2020_1_Skycode/internal/orders"
	"github.com/2020_1_Skycode/internal/tools"
	"golang.org/x/net/context"
)

type OrderManager struct {
	orderRepo         orders.Repository
	notificationsRepo notifications.Repository
}

func NewOrderProtoManager(orderRepo orders.Repository, nr notifications.Repository) *OrderManager {
	return &OrderManager{
		orderRepo:         orderRepo,
		notificationsRepo: nr,
	}
}

func (oU *OrderManager) CheckOutOrder(ctx context.Context, c *Checkout) (*Error, error) {
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
		return &Error{Err: tools.CheckoutOrderError.Error()}, err
	}

	return &Error{}, nil
}

func (oU *OrderManager) ChangeOrderStatus(ctx context.Context, cs *ChangeStatus) (*ErrorCode, error) {
	o, err := oU.orderRepo.GetByID(cs.OrderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &ErrorCode{ID: tools.DoesntExist}, nil
		}

		return &ErrorCode{ID: tools.InternalError}, err
	}

	if o.Status == cs.Status {
		return &ErrorCode{ID: tools.SameStatus}, nil
	}

	if err := oU.orderRepo.ChangeStatus(cs.OrderID, cs.Status); err != nil {
		return &ErrorCode{ID: tools.InternalError}, err
	}

	note := &models.Notification{
		UserID:  o.UserID,
		OrderID: o.ID,
		Status:  cs.Status,
	}

	if err := oU.notificationsRepo.InsertInto(note); err != nil {
		return &ErrorCode{ID: tools.InternalError}, err
	}

	return &ErrorCode{ID: tools.OK}, nil
}

func (oU *OrderManager) GetAllUserOrders(ctx context.Context, u *UserOrders) (*GetAllResponse, error) {
	userOrders, total, err := oU.orderRepo.GetAllByUserID(u.UserID, u.Count, u.Page)

	orders := []*Order{}

	for _, val := range userOrders {
		products := []*Product{}

		for _, val := range val.Products {
			products = append(products, &Product{
				ID:     val.ID,
				Name:   val.Name,
				Price:  val.Price,
				Image:  val.Image,
				RestID: val.RestId,
				Count:  val.Count,
			})
		}

		orders = append(orders, &Order{
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

	res := &GetAllResponse{
		Order: orders,
		Total: total,
	}

	if err != nil {
		return res, err
	}

	return res, nil
}

func (oU *OrderManager) GetOrderByID(ctx context.Context, u *GetByID) (*GetByIDResponse, error) {
	order, err := oU.orderRepo.GetByID(u.OrderID)

	if err != nil {
		return &GetByIDResponse{}, err
	}

	products := []*Product{}

	for _, val := range order.Products {
		products = append(products, &Product{
			ID:     val.ID,
			Name:   val.Name,
			Price:  val.Price,
			Image:  val.Image,
			RestID: val.RestId,
			Count:  val.Count,
		})
	}

	resOrder := &Order{
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

	res := &GetByIDResponse{
		Order: resOrder,
	}

	return res, err
}

func (oU *OrderManager) DeleteOrder(ctx context.Context, d *DelOrder) (*Error, error) {
	if err := oU.orderRepo.DeleteOrder(d.OrderID, d.UserID); err != nil {
		return &Error{Err: err.Error()}, err
	}

	return &Error{}, nil
}
