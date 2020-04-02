package delivery

import (
	"github.com/2020_1_Skycode/internal/middlewares"
	"github.com/2020_1_Skycode/internal/orders"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	OrderUseCase orders.UseCase
	MiddlewareC    *middlewares.MWController
}

func NewOrderHandler(router *gin.Engine, orderUC orders.UseCase, mw *middlewares.MWController) *OrderHandler {
	oh := &OrderHandler{
		OrderUseCase: orderUC,
		MiddlewareC: mw,
	}

	router.POST("api/v1/orders/checkout", oh.MiddlewareC.CheckAuth(), )
	router.GET("api/v1/orders/:orderID")

	return oh
}

func (oH *OrderHandler) Checkout() gin.HandlerFunc {
	type OrderStructRequest struct {

	}
	return func(c *gin.Context) {

	}
}
