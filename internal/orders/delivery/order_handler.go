package delivery

import (
	"github.com/2020_1_Skycode/internal/middlewares"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/orders"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type OrderHandler struct {
	OrderUseCase orders.UseCase
	MiddlewareC  *middlewares.MWController
}

func NewOrderHandler(router *gin.Engine, orderUC orders.UseCase, mw *middlewares.MWController) *OrderHandler {
	oh := &OrderHandler{
		OrderUseCase: orderUC,
		MiddlewareC:  mw,
	}

	router.POST("api/v1/orders/checkout", oh.MiddlewareC.CheckAuth(), oh.Checkout())
	router.GET("api/v1/orders/:orderID")

	return oh
}

func (oH *OrderHandler) Checkout() gin.HandlerFunc {
	type OrderStructRequest struct {
		UserID    uint64                 `json:"userId" binding:"required"`
		Address   string                 `json:"address" binding:"required"`
		Comment   string                 `json:"comment" binding:"required"`
		PersonNum uint16                 `json:"personNum" binding:"required"`
		Products  []*models.OrderProduct `json:"products" binding:"required"`
		Price     float32                `json:"price" binding:"required"`
	}
	return func(c *gin.Context) {
		req := &OrderStructRequest{}

		usr, exists := c.Get("user")

		if !exists {
			c.JSON(http.StatusUnauthorized, tools.Error{
				ErrorMessage: tools.Unauthorized.Error(),
			})

			return
		}

		_, ok := usr.(*models.User)

		if !ok {
			c.JSON(http.StatusInternalServerError, tools.Error{
				ErrorMessage: tools.UserTypeAssertionErr.Error(),
			})
			return
		}

		if err := c.Bind(req); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusOK, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})
		}

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
		}

		c.JSON(http.StatusOK, tools.Message{
			Message: "success",
		})
	}
}
