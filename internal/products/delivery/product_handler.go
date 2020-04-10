package delivery

import (
	"github.com/2020_1_Skycode/internal/middlewares"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/products"
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
)

type ProductHandler struct {
	productUseCase products.UseCase
	restUseCase    restaurants.UseCase
	middlewareC    *middlewares.MWController
	v              *requestValidator.RequestValidator
}

func NewProductHandler(private *gin.RouterGroup, public *gin.RouterGroup, pUC products.UseCase,
	validator *requestValidator.RequestValidator, rUC restaurants.UseCase, mw *middlewares.MWController) *ProductHandler {
	ph := &ProductHandler{
		productUseCase: pUC,
		middlewareC:    mw,
		restUseCase:    rUC,
		v:              validator,
	}

	public.GET("/products/:prod_id", ph.GetProduct())
	public.GET("/restaurants/:rest_id/product", ph.GetProducts())

	private.POST("/restaurants/:rest_id/product", ph.CreateProduct())
	private.PUT("/products/:prod_id/update", ph.UpdateProduct())
	private.PUT("/products/:prod_id/image", ph.UpdateImage())
	private.DELETE("/products/:prod_id/delete", ph.DeleteProduct())

	return ph
}

type productRequest struct {
	Name  string  `json:"name, omitempty" binding:"required" validate:"min=2"`
	Price float32 `json:"price, omitempty" binding:"required"`
}

//@Tags Product
//@Summary Get Product Route
//@Description Returning Product Model
//@Accept json
//@Produce json
//@Param prod_id path int true "Product ID"
//@Success 200 object models.Product
//@Failure 400 object tools.Error
//@Failure 404 object tools.Error
//@Router /product/prod:id [get]
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

//@Tags Product
//@Summary Get Products Of Restaurant Route
//@Description Returning Products List of Restaurant
//@Accept json
//@Produce json
//@Param count query int true "Count of elements on page"
//@Param page query int true "Number of page"
//@Param rest_id path int true "Id of restaurant"
//@Success 200 array models.Product
//@Failure 400 object tools.Error
//@Router /restaurants/rest:id/product [get]
func (ph *ProductHandler) GetProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		count, err := strconv.ParseUint(c.Query("count"), 10, 64)
		page, err := strconv.ParseUint(c.Query("page"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: "Bad params",
			})

			return
		}

		id, err := strconv.ParseUint(c.Param("rest_id"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		products, total, err := ph.productUseCase.GetProductsByRestaurantID(id, count, page)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusNotFound, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"products": products,
			"total":    total,
		})
	}
}

//@Tags Product
//@Summary Create Product Route
//@Description Creating Product
//@Accept json
//@Produce json
//@Param rest_id path int true "Restaurant ID"
//@Param Name formData string true "New product name"
//@Param Price formData number true "New product price"
//@Param image formData file true "New product image"
//@Success 200 object tools.Message
//@Failure 400 object tools.Error
//@Router /restaurants/rest:id/product [post]
func (ph *ProductHandler) CreateProduct() gin.HandlerFunc {
	rootDir, _ := os.Getwd()
	return func(c *gin.Context) {
		req := &productRequest{}

		user, err := ph.middlewareC.GetUser(c)

		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		if !user.IsManager() && !user.IsAdmin() {
			c.JSON(http.StatusForbidden, tools.Error{
				ErrorMessage: "User doesn't have permissions",
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

		rest, err := ph.restUseCase.GetRestaurantByID(restID)

		if err != nil {
			c.JSON(http.StatusNotFound, tools.Error{
				ErrorMessage: tools.RestaurantNotFoundError.Error(),
			})

			return
		}

		if rest.ManagerID != user.ID && !user.IsAdmin() {
			logrus.Info(rest, user)
			c.JSON(http.StatusForbidden, tools.Error{
				ErrorMessage: tools.RestaurantPermissionsError.Error(),
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

		errorsList := ph.v.ValidateRequest(req)

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

		if err := c.SaveUploadedFile(file, filepath.Join(rootDir, tools.ProductImagesPath, filename)); err != nil {
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
			Image:  filename,
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

//@Tags Product
//@Summary Update Product Route
//@Description Updating Product
//@Accept json
//@Produce json
//@Param prod_id path int true "Product ID"
//@Param ProdReq body productRequest true "New product data"
//@Success 200 object tools.Message
//@Failure 400 object tools.Error
//@Router /product/:prod_id/update [put]
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

		errorsList := ph.v.ValidateRequest(req)

		if len(*errorsList) > 0 {
			logrus.Info(tools.NotRequiredFields)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.NotRequiredFields.Error(),
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

//@Tags Product
//@Summary Update Product Image Route
//@Description Updating Product Image
//@Accept mpfd
//@Produce json
//@Param prod_id path int true "Product ID"
//@Param ProdReq formData file true "New product image"
//@Success 200 object tools.Message
//@Failure 400 object tools.Error
//@Failure 500 object tools.Error
//@Router /product/:prod_id/image [put]
func (ph *ProductHandler) UpdateImage() gin.HandlerFunc {
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

		filename := shortuuid.New()

		if err := c.SaveUploadedFile(file, filepath.Join(rootDir, tools.AvatarPath, filename)); err != nil {
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
			if err := os.Remove(filepath.Join(rootDir, tools.ProductImagesPath, product.Image)); err != nil {
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

//@Tags Product
//@Summary Update Product Image Route
//@Description Updating Product Image
//@Accept json
//@Produce json
//@Param prod_id path int true "Product ID"
//@Success 200 object tools.Message
//@Failure 400 object tools.Error
//@Failure 500 object tools.Error
//@Router /product/:prod_id [delete]
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
