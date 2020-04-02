package delivery

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"

	"github.com/2020_1_Skycode/internal/restaurants"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/gin-gonic/gin"
)

type RestaurantHandler struct {
	restUseCase restaurants.UseCase
}

func NewRestaurantHandler(router *gin.Engine, rUC restaurants.UseCase) *RestaurantHandler {
	rh := &RestaurantHandler{
		restUseCase: rUC,
	}

	router.GET("api/v1/restaurants", rh.GetRestaurants())
	router.GET("api/v1/restaurants/:id", rh.GetRestaurantByID())

	return rh
}

func (rh *RestaurantHandler) GetRestaurants() gin.HandlerFunc {
	return func(c *gin.Context) {
		restList, err := rh.restUseCase.GetRestaurants()
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: "text",
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"Restaurants": restList,
		})
	}
}

func (rh *RestaurantHandler) GetRestaurantByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})
		}
		rest, err := rh.restUseCase.GetRestaurantByID(id)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusNotFound, tools.Error{ // change error
				ErrorMessage: "text",
			})

			return
		}

		c.JSON(http.StatusOK, rest)
	}
}
