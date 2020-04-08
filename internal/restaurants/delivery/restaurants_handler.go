package delivery

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/restaurants"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/2020_1_Skycode/internal/tools/requestValidator"
	"github.com/gin-gonic/gin"
	"github.com/renstrom/shortuuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
)

type RestaurantHandler struct {
	restUseCase restaurants.UseCase
	v           *requestValidator.RequestValidator
}

func NewRestaurantHandler(private *gin.RouterGroup, public *gin.RouterGroup,
	validator *requestValidator.RequestValidator, rUC restaurants.UseCase) *RestaurantHandler {
	rh := &RestaurantHandler{
		restUseCase: rUC,
		v: validator,
	}

	public.GET("/restaurants", rh.GetRestaurants())
	public.GET("/restaurants/:rest_id", rh.GetRestaurantByID())

	private.POST("/restaurants", rh.CreateRestaurant())
	private.PUT("/restaurants/:rest_id/update", rh.UpdateRestaurant())
	private.PUT("/restaurants/:rest_id/image", rh.UpdateImage())
	private.DELETE("/restaurants/:rest_id", rh.DeleteRestaurant())

	return rh
}

type restaurantRequest struct {
	Name        string `json:"name, omitempty" binding:"required" validate:"min=3"`
	Description string `json:"description, omitempty" binding:"required" validate:"min=10"`
}

//@Tags Restaurant
//@Summary Get Restaurants List Route
//@Description Returning list of all restaurants
//@Accept json
//@Produce json
//@Success 200 array models.Restaurant
//@Failure 400 object tools.Error
//@Router /restaurants [get]
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

//@Tags Restaurant
//@Summary Get Restaurant By ID Route
//@Description Returning Restaurant Model
//@Accept json
//@Produce json
//@Param rest_id path uint64 true "Restaurant ID"
//@Success 200 object models.Restaurant
//@Failure 400 object tools.Error
//@Failure 404 object tools.Error
//@Router /restaurants/:rest_id [get]
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

//@Tags Restaurant
//@Summary Create New Restaurant Route
//@Description Add new restaurant
//@Accept json
//@Produce json
//@Param RestReq body restaurantRequest true "New restaurant data"
//@Success 200 object tools.Message
//@Failure 400 object tools.Error
//@Router /restaurants [post]
func (rh *RestaurantHandler) CreateRestaurant() gin.HandlerFunc {

	return func(c *gin.Context) {
		req := &restaurantRequest{}

		if err := c.Bind(req); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		errorsList := rh.v.ValidateRequest(req)

		if len(*errorsList) > 0 {
			logrus.Info(tools.NotRequiredFields)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.NotRequiredFields.Error(),
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

//@Tags Restaurant
//@Summary Update Restaurant Route
//@Description Updating Restaurant
//@Accept json
//@Produce json
//@Param rest_id path uint64 true "Restaurant ID"
//@Param RestReq body restaurantRequest true "New restaurant data"
//@Success 200 object tools.Message
//@Failure 400 object tools.Error
//@Router /restaurants/:rest_id/update [put]
func (rh *RestaurantHandler) UpdateRestaurant() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &restaurantRequest{}

		if err := c.Bind(req); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		errorsList := rh.v.ValidateRequest(req)

		if len(*errorsList) > 0 {
			logrus.Info(tools.NotRequiredFields)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.NotRequiredFields.Error(),
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

//@Tags Restaurant
//@Summary Update Restaurant Image Route
//@Description Updating Restaurant Image
//@Accept json
//@Produce mpfd
//@Param rest_id path uint64 true "Restaurant ID"
//@Param image formData file true "New restaurant image"
//@Success 200 object tools.Message
//@Failure 400 object tools.Error
//@Failure 500 object tools.Error
//@Router /restaurants/:rest_id/image [put]
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

//@Tags Restaurant
//@Summary Delete Restaurant Route
//@Description Deleting Restaurant
//@Accept json
//@Produce json
//@Param rest_id path uint64 true "Restaurant ID"
//@Success 200 object tools.Message
//@Failure 400 object tools.Error
//@Failure 500 object tools.Error
//@Router /restaurants/:rest_id [delete]
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
