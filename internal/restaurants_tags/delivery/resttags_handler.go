package delivery

import (
	"github.com/2020_1_Skycode/internal/middlewares"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/restaurants_tags"
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

type RestTagsHandler struct {
	restTagsUCase restaurants_tags.UseCase
	middlewareC   *middlewares.MWController
	v             *requestValidator.RequestValidator
}

func NewRestTagHandler(private *gin.RouterGroup, public *gin.RouterGroup, rtUC restaurants_tags.UseCase,
	validator *requestValidator.RequestValidator, mw *middlewares.MWController) *RestTagsHandler {
	rth := &RestTagsHandler{
		restTagsUCase: rtUC,
		middlewareC:   mw,
		v:             validator,
	}

	public.GET("/rest_tags", rth.GetAllTags())
	public.GET("/rest_tags/:id", rth.GetTagByID())

	private.POST("/rest_tags", rth.CreateTag())
	private.PUT("/rest_tags/:id/image", rth.UpdateTagImage())
	private.PUT("/rest_tags/:id", rth.UpdateTag())
	private.DELETE("/rest_tags/:id", rth.DeleteTag())

	return rth
}

type tagRequest struct {
	Name string `json:"name, omitempty" binding:"required" validate:"min=2"`
}

func (rth *RestTagsHandler) GetAllTags() gin.HandlerFunc {
	return func(c *gin.Context) {
		tags, err := rth.restTagsUCase.GetAllTags()
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"rest_tags": tags,
		})
	}
}

func (rth *RestTagsHandler) GetTagByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		tag, err := rth.restTagsUCase.GetTagByID(id)
		if err != nil {
			if err == tools.RestTagNotFound {
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
			"rest_tag": tag,
		})
	}
}

func (rth *RestTagsHandler) CreateTag() gin.HandlerFunc {
	rootDir, _ := os.Getwd()
	return func(c *gin.Context) {
		user, err := rth.middlewareC.GetUser(c)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		if !user.IsAdmin() {
			c.JSON(http.StatusForbidden, tools.Error{
				ErrorMessage: tools.PermissionError.Error(),
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

		req := &tagRequest{}
		if err := req.UnmarshalJSON(data); err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		errorsList := rth.v.ValidateRequest(req)
		if len(*errorsList) > 0 {
			logrus.Error(tools.NotRequiredFields)
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

		if err := c.SaveUploadedFile(file, filepath.Join(rootDir, tools.RestTagsImagesPath, filename)); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		tag := &models.RestTag{
			Name:  req.Name,
			Image: filename,
		}

		if err := rth.restTagsUCase.CreateTag(tag); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, tools.Message{
			Message: "Restaurant tag has been created",
		})
	}
}

func (rth *RestTagsHandler) UpdateTag() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := rth.middlewareC.GetUser(c)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		if !user.IsAdmin() {
			c.JSON(http.StatusForbidden, tools.Error{
				ErrorMessage: tools.PermissionError.Error(),
			})

			return
		}

		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		req := &tagRequest{}
		if err := c.Bind(req); err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		errorsList := rth.v.ValidateRequest(req)
		if len(*errorsList) > 0 {
			logrus.Error(tools.NotRequiredFields)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.NotRequiredFields.Error(),
			})

			return
		}

		tag := &models.RestTag{
			ID:   id,
			Name: req.Name,
		}

		if err := rth.restTagsUCase.UpdateTag(tag); err != nil {
			if err == tools.RestTagNotFound {
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

		c.JSON(http.StatusOK, tools.Message{
			Message: "Restaurant tag has been updated",
		})
	}
}

func (rth *RestTagsHandler) UpdateTagImage() gin.HandlerFunc {
	rootDir, _ := os.Getwd()
	return func(c *gin.Context) {
		user, err := rth.middlewareC.GetUser(c)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		if !user.IsAdmin() {
			c.JSON(http.StatusForbidden, tools.Error{
				ErrorMessage: tools.PermissionError.Error(),
			})

			return
		}

		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
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

		if err := c.SaveUploadedFile(file, filepath.Join(rootDir, tools.RestTagsImagesPath, filename)); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		tag := &models.RestTag{
			ID:    id,
			Image: filename,
		}

		if err := rth.restTagsUCase.UpdateTag(tag); err != nil {
			if err == tools.RestTagNotFound {
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

		c.JSON(http.StatusOK, tools.Message{
			Message: "Restaurant tag image has been updated",
		})
	}
}

func (rth *RestTagsHandler) DeleteTag() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := rth.middlewareC.GetUser(c)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		if !user.IsAdmin() {
			c.JSON(http.StatusForbidden, tools.Error{
				ErrorMessage: tools.PermissionError.Error(),
			})

			return
		}

		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		if err := rth.restTagsUCase.DeleteTag(id); err != nil {
			if err == tools.RestTagNotFound {
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

		c.JSON(http.StatusOK, tools.Message{
			Message: "Success",
		})
	}
}
