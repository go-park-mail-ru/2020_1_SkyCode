package delivery

import (
	"github.com/2020_1_Skycode/internal/middlewares"
	"github.com/2020_1_Skycode/internal/restaurant_points"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type RestPointsHandler struct {
	restPointsUseCase restaurant_points.UseCase
	middlewareC       *middlewares.MWController
}

func NewRestPointsHandler(private *gin.RouterGroup, public *gin.RouterGroup,
	rpUC restaurant_points.UseCase, mw *middlewares.MWController) *RestPointsHandler {
	rph := &RestPointsHandler{
		restPointsUseCase: rpUC,
		middlewareC:       mw,
	}

	public.GET("/points/:id", rph.GetPoint())

	private.DELETE("/points/:id", rph.DeletePoint())

	return rph
}

func (rph *RestPointsHandler) GetPoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		pID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		returnPoint, err := rph.restPointsUseCase.GetPoint(pID)
		if err != nil {
			if err == tools.RestPointNotFound {
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

		c.JSON(http.StatusOK, returnPoint)
	}
}

func (rph *RestPointsHandler) DeletePoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := rph.middlewareC.GetUser(c)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		if !user.IsManager() && !user.IsAdmin() {
			c.JSON(http.StatusForbidden, tools.Error{
				ErrorMessage: tools.DontEnoughRights.Error(),
			})

			return
		}

		pID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		if err := rph.restPointsUseCase.Delete(pID); err != nil {
			if err == tools.RestPointNotFound {
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

		c.JSON(http.StatusOK, tools.Message{"Deleted"})
	}
}
