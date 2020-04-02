package delivery

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/products"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/gin-gonic/gin"
	"github.com/renstrom/shortuuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
)

type ProductHandler struct {
	productUseCase products.UseCase
}

func NewProductHandler(router *gin.Engine, pUC products.UseCase) *ProductHandler {
	ph := &ProductHandler{
		productUseCase: pUC,
	}

	router.GET("api/v1/products/:prod_id", ph.GetProduct())
	router.POST("api/v1/restaurants/:rest_id/product", ph.CreateProduct())
	router.POST("api/v1/products/:prod_id/update", ph.UpdateProduct())
	router.POST("api/v1/products/:prod_id/image", ph.UpdateImage())
	router.POST("api/v1/products/:prod_id/delete", ph.DeleteProduct())

	return ph
}

func (ph *ProductHandler) GetProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("prod_id"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		product, err := ph.productUseCase.GetProductByID(id)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusNotFound, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, product)
	}
}

func (ph *ProductHandler) CreateProduct() gin.HandlerFunc {
	type ProductRequest struct {
		Name  string  `json:"name, omitempty" binding:"required"`
		Price float32 `json:"price, omitempty" binding:"required"`
	}

	return func(c *gin.Context) {
		req := &ProductRequest{}

		if err := c.Bind(req); err != nil {
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

		product := &models.Product{
			Name:   req.Name,
			Price:  req.Price,
			RestId: restID,
		}

		if err = ph.productUseCase.CreateProduct(product); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, tools.Message{
			Message: "Product has been created",
		})
	}
}

func (ph *ProductHandler) UpdateProduct() gin.HandlerFunc {
	type UpdateProductRequest struct {
		Name  string  `json:"name, omitempty" binding:"required"`
		Price float32 `json:"price, omitempty" binding:"required"`
	}

	return func(c *gin.Context) {
		req := &UpdateProductRequest{}

		if err := c.Bind(req); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		prodID, err := strconv.ParseUint(c.Param("prod_id"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		product := &models.Product{
			ID:    prodID,
			Name:  req.Name,
			Price: req.Price,
		}

		if err = ph.productUseCase.UpdateProduct(product); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, tools.Message{
			Message: "Product has been updated",
		})
	}
}

func (ph *ProductHandler) UpdateImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("image")

		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		filename := shortuuid.New() + "-" + file.Filename

		if err := c.SaveUploadedFile(file, tools.ProductImagesPath+filename); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		prodID, err := strconv.ParseUint(c.Param("prod_id"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		product, err := ph.productUseCase.GetProductByID(prodID)

		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		if product.Image != "" {
			if err := os.Remove(tools.ProductImagesPath + product.Image); err != nil {
				logrus.Info(err)
				c.JSON(http.StatusInternalServerError, tools.Error{
					ErrorMessage: tools.DeleteAvatarError.Error(),
				})

				return
			}
		}

		if err = ph.productUseCase.UpdateProductImage(prodID, filename); err != nil {
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

func (ph *ProductHandler) DeleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		prodID, err := strconv.ParseUint(c.Param("prod_id"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})
		}

		product, err := ph.productUseCase.GetProductByID(prodID)

		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		if product.Image != "" {
			if err := os.Remove(tools.ProductImagesPath + product.Image); err != nil {
				logrus.Info(err)
				c.JSON(http.StatusInternalServerError, tools.Error{
					ErrorMessage: tools.DeleteAvatarError.Error(),
				})

				return
			}
		}

		if err = ph.productUseCase.DeleteProduct(prodID); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})
		}

		c.JSON(http.StatusOK, tools.Message{
			Message: "success",
		})
	}
}
