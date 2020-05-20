package delivery

import (
	"fmt"
	"github.com/2020_1_Skycode/internal/middlewares"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/orders"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/2020_1_Skycode/internal/tools/notificationsWS"
	"github.com/2020_1_Skycode/internal/tools/requestValidator"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type OrderHandler struct {
	orderUseCase orders.UseCase
	noteServer   *notificationsWS.NotificationServer
	middlewareC  *middlewares.MWController
	v            *requestValidator.RequestValidator
}

func NewOrderHandler(private *gin.RouterGroup, public *gin.RouterGroup, orderUC orders.UseCase,
	validator *requestValidator.RequestValidator, mw *middlewares.MWController,
	ns *notificationsWS.NotificationServer) *OrderHandler {
	oh := &OrderHandler{
		orderUseCase: orderUC,
		noteServer:   ns,
		middlewareC:  mw,
		v:            validator,
	}
	public.GET("/orders", oh.GetUserOrders())
	public.GET("/orders/:orderID", oh.GetUserOrder())

	private.POST("/orders/:orderID/status", oh.ChangeStatus())
	private.POST("/orders", oh.Checkout())
	private.DELETE("/orders/:orderID", oh.DeleteOrder())

	return oh
}

type orderRequest struct {
	UserID    uint64                 `json:"userId" binding:"required"`
	RestID    uint64                 `json:"restId" binding:"required"`
	Address   string                 `json:"address" binding:"required" validate:"min=5"`
	Comment   string                 `json:"comment"`
	Phone     string                 `json:"phone" binding:"required" validate:"min=11,max=15"`
	PersonNum uint32                 `json:"personNum" binding:"required"`
	Products  []*models.OrderProduct `json:"products" binding:"required" required:"dive,required"`
	Price     float32                `json:"price" binding:"required"`
}

//@Tags Order
//@Summary Create Order Route
//@Description Creating Order
//@Accept json
//@Produce json
//@Param OrderReq body orderRequest true "New order data"
//@Success 200 object tools.Message
//@Failure 400 object tools.Error
//@Failure 401 object tools.Error
//@Failure 500 object tools.Error
//@Security basicAuth
//@Router /orders/checkout [post]
func (oH *OrderHandler) Checkout() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &orderRequest{}

		_, err := oH.middlewareC.GetUser(c)

		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		data, err := c.GetRawData()

		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BindingError.Error(),
			})

			return
		}

		if err := req.UnmarshalJSON(data); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.NotRequiredFields.Error(),
			})

			return
		}

		errorsList := oH.v.ValidateRequest(req)

		if len(*errorsList) > 0 {
			logrus.Info(tools.NotRequiredFields)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.NotRequiredFields.Error(),
			})

			return
		}

		fmt.Println(req)

		order := &models.Order{
			UserID:    req.UserID,
			RestID:    req.RestID,
			Address:   req.Address,
			Comment:   req.Comment,
			PersonNum: req.PersonNum,
			Price:     req.Price,
			Phone:     req.Phone,
		}

		if err := oH.orderUseCase.CheckoutOrder(order, req.Products); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, tools.Message{
			Message: "success",
		})
	}
}

//@Tags Order
//@Summary Create Order Route
//@Description Creating Order
//@Accept json
//@Produce json
//@Param count query int true "Count of elements on page"
//@Param page query int true "Number of page"
//@Success 200 array models.Order
//@Failure 400 object tools.Error
//@Failure 500 object tools.Error
//@Security basicAuth
//@Router /orders [get]
func (oH *OrderHandler) GetUserOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := oH.middlewareC.GetUser(c)

		if err != nil {
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		count, cerr := strconv.ParseUint(c.Query("count"), 10, 64)
		page, perr := strconv.ParseUint(c.Query("page"), 10, 64)

		if cerr != nil || perr != nil {
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadQueryParams.Error(),
			})

			return
		}

		userOrders, total, err := oH.orderUseCase.GetAllUserOrders(user.ID, count, page)

		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusInternalServerError, tools.Error{
				ErrorMessage: tools.GetOrdersError.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, tools.Body{
			"orders": userOrders,
			"total":  total,
		})
	}
}

//@Tags Order
//@Summary Create Order Route
//@Description Creating Order
//@Accept json
//@Produce json
//@Param order_id path integer true "ID of order"
//@Success 200 object models.Order
//@Failure 400 object tools.Error
//@Failure 404 object tools.Error
//@Security basicAuth
//@Router /orders/:order_id [get]
func (oH *OrderHandler) GetUserOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := oH.middlewareC.GetUser(c)

		if err != nil {
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		orderID, err := strconv.ParseUint(c.Param("orderID"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})
		}

		userOrders, err := oH.orderUseCase.GetOrderByID(orderID, user.ID)

		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusNotFound, tools.Error{
				ErrorMessage: tools.NotFound.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, tools.Body{
			"order": userOrders,
		})
	}
}

func (oH *OrderHandler) ChangeStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := oH.middlewareC.GetUser(c)

		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		if !user.IsManager() && !user.IsAdmin() {
			c.JSON(http.StatusForbidden, tools.Error{
				ErrorMessage: tools.PermissionError.Error(),
			})

			return
		}

		orderID, err := strconv.ParseUint(c.Param("orderID"), 10, 64)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		newStatusCode, err := strconv.ParseUint(c.Query("status"), 10, 64)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		newStatus := tools.StatusCodes[newStatusCode]
		if newStatus == "" {
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.UnknownOrderStatus.Error(),
			})

			return
		}

		note, err := oH.orderUseCase.ChangeOrderStatus(orderID, newStatus)

		if err != nil {
			if err == tools.NewStatusIsTheSame {
				c.JSON(http.StatusConflict, tools.Error{
					ErrorMessage: err.Error(),
				})

				return
			}
			if err == tools.OrderNotFound {
				c.JSON(http.StatusBadRequest, tools.Error{
					ErrorMessage: err.Error(),
				})

				return
			}

			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		oH.noteServer.SendNotification(note)

		c.JSON(http.StatusOK, tools.Message{Message: "Order status successfully updated"})
	}
}

//@Tags Order
//@Summary Create Order Route
//@Description Creating Order
//@Accept json
//@Produce json
//@Param order_id path integer true "ID of order"
//@Success 200 object tools.Message
//@Failure 400 object tools.Error
//@Failure 404 object tools.Error
//@Security basicAuth
//@Router /orders/:order_id [delete]
func (oH *OrderHandler) DeleteOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := oH.middlewareC.GetUser(c)

		if err != nil {
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		orderID, err := strconv.ParseUint(c.Param("orderID"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		if err := oH.orderUseCase.DeleteOrder(orderID, user.ID); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusNotFound, tools.Error{
				ErrorMessage: tools.NotFound.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, tools.Message{
			Message: "success",
		})
	}
}
