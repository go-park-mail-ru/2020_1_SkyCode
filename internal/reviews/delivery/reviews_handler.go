package delivery

import (
	"github.com/2020_1_Skycode/internal/middlewares"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/reviews"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/2020_1_Skycode/internal/tools/requestValidator"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type ReviewsHandler struct {
	reviewsUseCase reviews.UseCase
	middlewareC    *middlewares.MWController
	v              *requestValidator.RequestValidator
}

func NewReviewsHandler(private *gin.RouterGroup, public *gin.RouterGroup, rUC reviews.UseCase,
	v *requestValidator.RequestValidator, middlewareC *middlewares.MWController) *ReviewsHandler {
	rh := &ReviewsHandler{
		reviewsUseCase: rUC,
		middlewareC:    middlewareC,
		v:              v,
	}

	public.GET("/reviews", rh.GetUserReviews())
	public.GET("/reviews/:rev_id", rh.GetReview())

	private.PUT("/reviews/:rev_id", rh.UpdateReview())
	private.DELETE("/reviews/:rev_id", rh.DeleteReview())

	return rh
}

type reviewUpdateRequest struct {
	Text string  `json:"text, omitempty" binding:"required"`
	Rate float64 `json:"rate" binding:"required" validate:"min=0,max=5"`
}

func (rh *ReviewsHandler) GetUserReviews() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := rh.middlewareC.GetUser(c)

		if err != nil {
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		count, err := strconv.ParseUint(c.Query("count"), 10, 64)
		page, err := strconv.ParseUint(c.Query("page"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadQueryParams.Error(),
			})

			return
		}

		returnReviews, total, err := rh.reviewsUseCase.GetUserReviews(user.ID, count, page)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusInternalServerError, tools.Error{
				ErrorMessage: tools.GetOrdersError.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"reviews": returnReviews,
			"total":   total,
		})
	}
}

func (rh *ReviewsHandler) GetReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		reviewID, err := strconv.ParseUint(c.Param("rev_id"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		returnReview, err := rh.reviewsUseCase.GetReview(reviewID)
		if err != nil {
			if err == tools.ReviewNotFoundError {
				c.JSON(http.StatusNotFound, tools.Error{
					ErrorMessage: tools.ReviewNotFoundError.Error(),
				})

				return
			}

			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, tools.Body{
			"review": returnReview,
		})
	}
}

func (rh *ReviewsHandler) UpdateReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		reviewID, err := strconv.ParseUint(c.Param("rev_id"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
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

		req := &reviewUpdateRequest{}
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

		reviewReq := &models.Review{
			ID:   reviewID,
			Text: req.Text,
			Author: &models.User{
				ID: user.ID,
			},
			Rate: req.Rate,
		}

		if err := rh.reviewsUseCase.UpdateReview(reviewReq, user); err != nil {
			if err == tools.ReviewNotFoundError {
				c.JSON(http.StatusNotFound, tools.Error{
					ErrorMessage: tools.ReviewNotFoundError.Error(),
				})

				return
			}
			if err == tools.PermissionError {
				c.JSON(http.StatusForbidden, tools.Error{
					ErrorMessage: tools.PermissionError.Error(),
				})

				return
			}

			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, tools.Message{"Updated"})
	}
}

func (rh *ReviewsHandler) DeleteReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		reviewID, err := strconv.ParseUint(c.Param("rev_id"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
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

		if err := rh.reviewsUseCase.DeleteReview(reviewID, user); err != nil {
			if err == tools.ReviewNotFoundError {
				c.JSON(http.StatusNotFound, tools.Error{
					ErrorMessage: tools.ReviewNotFoundError.Error(),
				})

				return
			}
			if err == tools.PermissionError {
				c.JSON(http.StatusForbidden, tools.Error{
					ErrorMessage: tools.PermissionError.Error(),
				})

				return
			}

			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, tools.Message{"Deleted"})
	}
}
