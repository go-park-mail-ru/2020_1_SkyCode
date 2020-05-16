package delivery

import (
	"github.com/2020_1_Skycode/internal/middlewares"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/restaurants"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/2020_1_Skycode/internal/tools/requestValidator"
	"github.com/gin-gonic/gin"
	"github.com/renstrom/shortuuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type RestaurantHandler struct {
	restUseCase restaurants.UseCase
	middlewareC *middlewares.MWController
	v           *requestValidator.RequestValidator
}

func NewRestaurantHandler(private *gin.RouterGroup, public *gin.RouterGroup,
	validator *requestValidator.RequestValidator, rUC restaurants.UseCase,
	mw *middlewares.MWController) *RestaurantHandler {
	rh := &RestaurantHandler{
		restUseCase: rUC,
		middlewareC: mw,
		v:           validator,
	}

	public.GET("/restaurants", rh.GetRestaurants())
	public.GET("/restaurants_point", rh.GetRestaurantsWithCloserPoint())
	public.GET("/restaurants/:rest_id", rh.GetRestaurantByID())

	private.POST("/restaurants", rh.CreateRestaurant())
	private.PUT("/restaurants/:rest_id/update", rh.UpdateRestaurant())
	private.PUT("/restaurants/:rest_id/image", rh.UpdateImage())
	private.DELETE("/restaurants/:rest_id", rh.DeleteRestaurant())

	private.POST("/restaurants/:rest_id/points", rh.AddPoint())
	public.GET("/restaurants/:rest_id/points", rh.GetPoints())

	private.POST("/restaurants/:rest_id/reviews", rh.AddReview())
	public.GET("/restaurants/:rest_id/reviews", rh.GetReviews())

	private.POST("/restaurants/:rest_id/tag/:tag_id", rh.AddTag())
	private.DELETE("/restaurants/:rest_id/tag/:tag_id", rh.DeleteTag())

	public.GET("/restaurants/:rest_id/prod_tags", rh.GetProductTags())
	private.POST("/restaurants/:rest_id/prod_tags", rh.AddProductTag())
	private.DELETE("/restaurants/:rest_id/prod_tags/:tag_id", rh.DeleteProductTag())

	return rh
}

type restaurantRequest struct {
	Name        string `json:"name, omitempty" binding:"required" validate:"min=3"`
	Description string `json:"description, omitempty" binding:"required" validate:"min=10"`
}

type reviewRequest struct {
	Text string   `json:"text, omitempty" binding:"required"`
	Rate *float64 `json:"rate" binding:"required" validate:"min=0,max=5"`
}

type pointRequest struct {
	Address string  `json:"address" binding:"required"`
	Radius  float64 `json:"radius" binding:"required"`
}

type productTagRequest struct {
	Name string `json:"name" binding:"required" validate:"min=2"`
}

//@Tags Restaurant
//@Summary Get Restaurants List Route
//@Description Returning list of all restaurants
//@Accept json
//@Produce json
//@Param count query int true "Count of elements on page"
//@Param page query int true "Number of page"
//@Success 200 array models.Restaurant
//@Failure 400 object tools.Error
//@Router /restaurants [get]
func (rh *RestaurantHandler) GetRestaurants() gin.HandlerFunc {
	return func(c *gin.Context) {
		count, err := strconv.ParseUint(c.Query("count"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadQueryParams.Error(),
			})

			return
		}
		page, err := strconv.ParseUint(c.Query("page"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadQueryParams.Error(),
			})

			return
		}
		tag := c.Query("tag")
		tagID := uint64(0)
		if tag != "" {
			tagID, err = strconv.ParseUint(tag, 10, 64)
			if err != nil {
				logrus.Info(err)
				c.JSON(http.StatusBadRequest, tools.Error{
					ErrorMessage: tools.BadQueryParams.Error(),
				})

				return
			}
		}

		restList, total, err := rh.restUseCase.GetRestaurants(count, page, tagID)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: "text",
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"restaurants": restList,
			"total":       total,
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
//@Param Name formData string true "New restaurant name"
//@Param Description formData string true "New restaurant price"
//@Param image formData file true "New restaurant image"
//@Success 200 object tools.Message
//@Failure 400 object tools.Error
//@Router /restaurants [post]
func (rh *RestaurantHandler) CreateRestaurant() gin.HandlerFunc {
	rootDir, _ := os.Getwd()
	return func(c *gin.Context) {
		req := &restaurantRequest{}

		user, err := rh.middlewareC.GetUser(c)

		if err != nil {
			logrus.Info(err)
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

		file, err := c.FormFile("image")

		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		filename := shortuuid.New()

		if err := c.SaveUploadedFile(file, filepath.Join(rootDir, tools.RestaurantImagesPath, filename)); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		rest := &models.Restaurant{
			ManagerID:   user.ID,
			Name:        req.Name,
			Description: req.Description,
			Image:       filename,
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
	rootDir, _ := os.Getwd()
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

		filename := shortuuid.New()

		if err := c.SaveUploadedFile(file, filepath.Join(rootDir, tools.RestaurantImagesPath, filename)); err != nil {
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
			if err := os.Remove(filepath.Join(rootDir, tools.RestaurantImagesPath, rest.Image)); err != nil {
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
	rootDir, _ := os.Getwd()
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
			if err := os.Remove(filepath.Join(rootDir, tools.RestaurantImagesPath, rest.Image)); err != nil {
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

func (rh *RestaurantHandler) AddPoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := rh.middlewareC.GetUser(c)
		if err != nil {
			logrus.Info(err)
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

		restID, err := strconv.ParseUint(c.Param("rest_id"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		req := &pointRequest{}
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

		p := &models.RestaurantPoint{
			Address:       req.Address,
			ServiceRadius: req.Radius,
			RestID:        restID,
		}

		err = rh.restUseCase.AddPoint(p)
		if err != nil {
			if err == tools.ApiAnswerEmptyResult || err == tools.RestaurantNotFoundError {
				c.JSON(http.StatusNotFound, tools.Error{
					ErrorMessage: err.Error(),
				})

				return
			}

			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, tools.Message{"Created"})
	}
}

func (rh *RestaurantHandler) GetRestaurantsWithCloserPoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		count, err := strconv.ParseUint(c.Query("count"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadQueryParams.Error(),
			})

			return
		}
		page, err := strconv.ParseUint(c.Query("page"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadQueryParams.Error(),
			})

			return
		}
		tag := c.Query("tag")
		tagID := uint64(0)
		if tag != "" {
			tagID, err = strconv.ParseUint(tag, 10, 64)
			if err != nil {
				logrus.Info(err)
				c.JSON(http.StatusBadRequest, tools.Error{
					ErrorMessage: tools.BadQueryParams.Error(),
				})

				return
			}
		}

		address := c.Query("address")
		if address == "" {
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadQueryParams.Error(),
			})

			return
		}

		returnRestaurants, total, err := rh.restUseCase.GetRestaurantsInServiceRadius(address, count, page, tagID)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"restaurants": returnRestaurants,
			"total":       total,
		})
	}
}

func (rh *RestaurantHandler) GetPoints() gin.HandlerFunc {
	return func(c *gin.Context) {
		count, err := strconv.ParseUint(c.Query("count"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadQueryParams.Error(),
			})

			return
		}
		page, err := strconv.ParseUint(c.Query("page"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadQueryParams.Error(),
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

		returnPoints, total, err := rh.restUseCase.GetPoints(restID, count, page)
		if err != nil {
			if err == tools.RestaurantNotFoundError {
				c.JSON(http.StatusNotFound, tools.Error{
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

		c.JSON(http.StatusOK, gin.H{
			"points": returnPoints,
			"total":  total,
		})
	}
}

func (rh *RestaurantHandler) AddReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		restID, err := strconv.ParseUint(c.Param("rest_id"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		req := &reviewRequest{}

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

		errorsList := rh.v.ValidateRequest(req)
		if len(*errorsList) > 0 {
			logrus.Info(tools.NotRequiredFields)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.ErrorRequestValidation.Error(),
			})

			return
		}

		user, err := rh.middlewareC.GetUser(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		newReview := &models.Review{
			RestID: restID,
			Text:   req.Text,
			Author: &models.User{
				ID: user.ID,
			},
			CreationDate: time.Now(),
			Rate:         *req.Rate,
		}

		if err := rh.restUseCase.AddReview(newReview); err != nil {
			if err == tools.RestaurantNotFoundError {
				c.JSON(http.StatusNotFound, tools.Error{
					ErrorMessage: tools.RestaurantNotFoundError.Error(),
				})

				return
			}

			if err == tools.ReviewAlreadyExists {
				c.JSON(http.StatusConflict, tools.Error{
					ErrorMessage: tools.ReviewAlreadyExists.Error(),
				})

				return
			}

			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, tools.Message{"Created"})
	}
}

func (rh *RestaurantHandler) GetReviews() gin.HandlerFunc {
	return func(c *gin.Context) {
		restID, err := strconv.ParseUint(c.Param("rest_id"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		count, err := strconv.ParseUint(c.Query("count"), 10, 64)
		page, err := strconv.ParseUint(c.Query("page"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadQueryParams.Error(),
			})

			return
		}

		user, err := rh.middlewareC.GetUser(c)
		if err != nil && err != tools.Unauthorized {
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		var userID uint64
		if user == nil {
			userID = 0
		} else {
			userID = user.ID
		}

		reviews, current, total, err := rh.restUseCase.GetReviews(restID, userID, count, page)
		if err != nil {
			if err == tools.RestaurantNotFoundError {
				c.JSON(http.StatusNotFound, tools.Error{
					ErrorMessage: tools.RestaurantNotFoundError.Error(),
				})

				return
			}

			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"reviews": reviews,
			"current": current,
			"total":   total,
		})
	}
}

func (rh *RestaurantHandler) AddTag() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := rh.middlewareC.GetUser(c)
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

		restID, err := strconv.ParseUint(c.Param("rest_id"), 10, 64)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		tagID, err := strconv.ParseUint(c.Param("tag_id"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadQueryParams.Error(),
			})

			return
		}

		if err := rh.restUseCase.AddTag(restID, tagID); err != nil {
			if err == tools.TagRestComboAlreadyExist {
				c.JSON(http.StatusConflict, tools.Error{
					ErrorMessage: err.Error(),
				})

				return
			}

			if err == tools.RestaurantNotFoundError || err == tools.RestTagNotFound {
				c.JSON(http.StatusNotFound, tools.Error{
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

		c.JSON(http.StatusOK, tools.Message{
			Message: "Tag added to the restaurant",
		})
	}
}

func (rh *RestaurantHandler) DeleteTag() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := rh.middlewareC.GetUser(c)
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

		restID, err := strconv.ParseUint(c.Param("rest_id"), 10, 64)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		tagID, err := strconv.ParseUint(c.Param("tag_id"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadQueryParams.Error(),
			})

			return
		}

		if err := rh.restUseCase.DeleteTag(restID, tagID); err != nil {
			if err == tools.TagRestComboDoesntExist {
				c.JSON(http.StatusConflict, tools.Error{
					ErrorMessage: err.Error(),
				})

				return
			}

			if err == tools.RestaurantNotFoundError || err == tools.RestTagNotFound {
				c.JSON(http.StatusNotFound, tools.Error{
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

		c.JSON(http.StatusOK, tools.Message{
			Message: "Tag deleted from the restaurant",
		})
	}
}

func (rh *RestaurantHandler) GetProductTags() gin.HandlerFunc {
	return func(c *gin.Context) {
		restID, err := strconv.ParseUint(c.Param("rest_id"), 10, 64)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		tags, err := rh.restUseCase.GetProductTagsByID(restID)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"Tags": tags,
		})
	}
}

func (rh *RestaurantHandler) AddProductTag() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := rh.middlewareC.GetUser(c)
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

		restID, err := strconv.ParseUint(c.Param("rest_id"), 10, 64)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		data, err := c.GetRawData()

		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BindingError.Error(),
			})

			return
		}

		req := &productTagRequest{}

		if err := req.UnmarshalJSON(data); err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		errorsList := rh.v.ValidateRequest(req)
		if len(*errorsList) > 0 {
			logrus.Info(tools.NotRequiredFields)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.ErrorRequestValidation.Error(),
			})

			return
		}

		tag := &models.ProductTag{
			Name:   req.Name,
			RestID: restID,
		}

		if err := rh.restUseCase.AddProductTag(tag); err != nil {
			if err == tools.TagRestComboAlreadyExist {
				c.JSON(http.StatusConflict, tools.Error{
					ErrorMessage: err.Error(),
				})

				return
			}

			if err == tools.RestaurantNotFoundError || err == tools.RestTagNotFound {
				c.JSON(http.StatusNotFound, tools.Error{
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

		c.JSON(http.StatusOK, tools.Message{
			Message: "Product tag added",
		})
	}
}

func (rh *RestaurantHandler) DeleteProductTag() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := rh.middlewareC.GetUser(c)
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

		tagID, err := strconv.ParseUint(c.Param("tag_id"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadQueryParams.Error(),
			})

			return
		}

		if err := rh.restUseCase.DeleteProductTag(tagID); err != nil {
			if err == tools.ProductTagNotFound {
				c.JSON(http.StatusConflict, tools.Error{
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

		c.JSON(http.StatusOK, tools.Message{
			Message: "Product tag deleted",
		})
	}
}
