package usecase

import (
	"context"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/tools/protobuf/orderswork"
	"google.golang.org/grpc"
)

type ProtoUseCase struct {
	orderManager orderswork.OrderWorkerClient
}

func NewOrderProtoUseCase(conn *grpc.ClientConn) *ProtoUseCase {
	return &ProtoUseCase{
		orderManager: orderswork.NewOrderWorkerClient(conn),
	}
}

func (oU *ProtoUseCase) CheckoutOrder(order *models.Order, ordProducts []*models.OrderProduct) error {
	ord := &orderswork.Order{
		UserID:    order.UserID,
		RestID:    order.RestID,
		RestName:  order.RestName,
		Address:   order.Address,
		Phone:     order.Phone,
		Comment:   order.Comment,
		PersonNum: order.PersonNum,
		Price:     order.Price,
	}

	prods := []*orderswork.OrderProduct{}

	for _, val := range ordProducts {
		prods = append(prods, &orderswork.OrderProduct{
			ID:        val.ID,
			ProductID: val.ProductID,
			Count:     val.Count,
		})
	}

	if _, err := oU.orderManager.CheckOutOrder(context.Background(), &orderswork.Checkout{
		Order:    ord,
		Products: prods,
	}); err != nil {
		return err
	}

	return nil
}

func (oU *ProtoUseCase) GetAllUserOrders(userID uint64, count uint64, page uint64) ([]*models.Order, uint64, error) {
	res, err := oU.orderManager.GetAllUserOrders(context.Background(), &orderswork.UserOrders{
		UserID: userID,
		Count:  count,
		Page:   page,
	})

	if err != nil {
		return nil, 0, err
	}

	orders := []*models.Order{}
	for _, val := range res.Order {
		products := []*models.Product{}

		for _, val := range val.Products {
			products = append(products, &models.Product{
				ID:     val.ID,
				Name:   val.Name,
				Price:  val.Price,
				Image:  val.Image,
				RestId: val.RestID,
				Count:  val.Count,
			})
		}

		orders = append(orders, &models.Order{
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

	return orders, res.Total, nil
}

func (oU *ProtoUseCase) GetOrderByID(orderID uint64, userID uint64) (*models.Order, error) {
	order, err := oU.orderManager.GetOrderByID(context.Background(), &orderswork.GetByID{
		OrderID: orderID,
		UserID:  userID,
	})

	if err != nil {
		return nil, err
	}

	products := []*models.Product{}

	for _, val := range order.Order.Products {
		products = append(products, &models.Product{
			ID:     val.ID,
			Name:   val.Name,
			Price:  val.Price,
			Image:  val.Image,
			RestId: val.RestID,
			Count:  val.Count,
		})
	}

	resOrder := &models.Order{
		ID:        order.Order.ID,
		UserID:    order.Order.UserID,
		RestID:    order.Order.RestID,
		RestName:  order.Order.RestName,
		Address:   order.Order.Address,
		Phone:     order.Order.Comment,
		Comment:   order.Order.Comment,
		PersonNum: order.Order.PersonNum,
		Products:  products,
		Price:     order.Order.Price,
		CreatedAt: order.Order.CreatedAt,
		Status:    order.Order.Status,
	}

	return resOrder, err
}

func (oU *ProtoUseCase) DeleteOrder(orderID uint64, userID uint64) error {
	if _, err := oU.orderManager.DeleteOrder(context.Background(), &orderswork.DelOrder{
		OrderID: orderID,
		UserID:  userID,
	}); err != nil {
		return err
	}

	return nil
}
