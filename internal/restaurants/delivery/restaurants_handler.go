package delivery

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/restaurants"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/gin-gonic/gin"
	"github.com/renstrom/shortuuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
)

type RestaurantHandler struct {
	restUseCase restaurants.UseCase
}

func NewRestaurantHandler(router *gin.Engine, rUC restaurants.UseCase) *RestaurantHandler {
	rh := &RestaurantHandler{
		restUseCase: rUC,
	}

	router.GET("api/v1/restaurants", rh.GetRestaurants())
	router.GET("api/v1/restaurants/:rest_id", rh.GetRestaurantByID())
	router.POST("api/v1/restaurants", rh.CreateRestaurant())
	router.POST("api/v1/restaurants/:rest_id/update", rh.UpdateRestaurant())
	router.POST("api/v1/restaurants/:rest_id/image", rh.UpdateImage())
	router.POST("api/v1/restaurants/:rest_id", rh.DeleteRestaurant())

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
		id, err := strconv.ParseUint(c.Param("rest_id"), 10, 64)
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

func (rh *RestaurantHandler) CreateRestaurant() gin.HandlerFunc {
	type RestaurantRequest struct {
		Name        string `json:"name, omitempty" binding:"required"`
		Description string `json:"description, omitempty" binding:"required"`
	}

	return func(c *gin.Context) {
		req := &RestaurantRequest{}

		if err := c.Bind(req); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		rest := &models.Restaurant{
			Name:        req.Name,
			Description: req.Description,
		}

		if err := rh.restUseCase.CreateRestaurant(rest); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, tools.Message{
			Message: "Restaurant has been created",
		})
	}
}

func (rh *RestaurantHandler) UpdateRestaurant() gin.HandlerFunc {
	type RestaurantRequest struct {
		Name        string `json:"name, omitempty" binding:"required"`
		Description string `json:"description, omitempty" binding:"required"`
	}

	return func(c *gin.Context) {
		req := &RestaurantRequest{}

		if err := c.Bind(req); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		prodID, err := strconv.ParseUint(c.Param("rest_id"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		rest := &models.Restaurant{
			ID:          prodID,
			Name:        req.Name,
			Description: req.Description,
		}

		if err = rh.restUseCase.UpdateRestaurant(rest); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, tools.Message{
			Message: "Restaurant has been updated",
		})
	}
}

func (rh *RestaurantHandler) UpdateImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("image")

		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		restID, err := strconv.ParseUint(c.Param("rest_id"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		filename := shortuuid.New() + "-" + file.Filename

		if err := c.SaveUploadedFile(file, tools.RestaurantImagesPath+filename); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		rest, err := rh.restUseCase.GetRestaurantByID(restID)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		if rest.Image != "" {
			if err := os.Remove(tools.RestaurantImagesPath + rest.Image); err != nil {
				logrus.Info(err)
				c.JSON(http.StatusInternalServerError, tools.Error{
					ErrorMessage: tools.DeleteAvatarError.Error(),
				})

				return
			}
		}

		if err = rh.restUseCase.UpdateImage(restID, filename); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, tools.Message{
			Message: "success",
		})
	}
}

func (rh *RestaurantHandler) DeleteRestaurant() gin.HandlerFunc {
	return func(c *gin.Context) {
		restID, err := strconv.ParseUint(c.Param("rest_id"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		rest, err := rh.restUseCase.GetRestaurantByID(restID)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		if rest.Image != "" {
			if err := os.Remove(tools.RestaurantImagesPath + rest.Image); err != nil {
				logrus.Info(err)
				c.JSON(http.StatusInternalServerError, tools.Error{
					ErrorMessage: tools.DeleteAvatarError.Error(),
				})

				return
			}
		}

		if err = rh.restUseCase.Delete(restID); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, tools.Message{
			Message: "success",
		})
	}
}
