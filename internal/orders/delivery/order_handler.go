package delivery

import (
	"fmt"
	"github.com/2020_1_Skycode/internal/middlewares"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/orders"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/2020_1_Skycode/internal/tools/requestValidator"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type OrderHandler struct {
	OrderUseCase orders.UseCase
	MiddlewareC  *middlewares.MWController
	v              *requestValidator.RequestValidator
}

func NewOrderHandler(private *gin.RouterGroup, public *gin.RouterGroup, orderUC orders.UseCase,
	validator *requestValidator.RequestValidator, mw *middlewares.MWController) *OrderHandler {
	oh := &OrderHandler{
		OrderUseCase: orderUC,
		MiddlewareC:  mw,
		v: validator,
	}

	private.POST("/orders/checkout", oh.Checkout())

	public.GET("/orders/:orderID", oh.Checkout())

	return oh
}

type orderRequest struct {
	UserID    uint64                 `json:"userId" binding:"required" validate:"regexp=numberRegexString"`
	Address   string                 `json:"address" binding:"required" validate:"min=5"`
	Comment   string                 `json:"comment"`
	Phone     string                 `json:"phone" binding:"required" validate:"min=11,max=15"`
	PersonNum uint16                 `json:"personNum" binding:"required" validate:"regexp=numberRegexString"`
	Products  []*models.OrderProduct `json:"products" binding:"required" required:"dive,required"`
	Price     float32                `json:"price" binding:"required" validate:"regexp=numberRegexString"`
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

		_, err := oH.MiddlewareC.GetUser(c)

		if err != nil {
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		if err := c.Bind(req); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
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
			Address:   req.Address,
			Comment:   req.Comment,
			PersonNum: req.PersonNum,
			Products:  req.Products,
			Price:     req.Price,
		}

		if err := oH.OrderUseCase.CheckoutOrder(order); err != nil {
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
