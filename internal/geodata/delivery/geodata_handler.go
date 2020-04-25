package delivery

import (
	"github.com/2020_1_Skycode/internal/geodata"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GeoDataHandler struct {
	geoDataUseCase geodata.UseCase
}

func NewGeoDataHandler(private *gin.RouterGroup, public *gin.RouterGroup, gdUC geodata.UseCase) *GeoDataHandler {
	gdh := &GeoDataHandler{geoDataUseCase: gdUC}

	public.GET("/check_address", gdh.CheckAddress())

	return gdh
}

func (gdh *GeoDataHandler) CheckAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Query("address")
		if address == "" {
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadQueryParams.Error(),
			})

			return
		}

		check, err := gdh.geoDataUseCase.CheckGeoPos(address)
		if err != nil || !check {
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, tools.Message{"OK"})
	}
}
