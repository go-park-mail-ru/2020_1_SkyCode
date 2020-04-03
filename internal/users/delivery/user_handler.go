package delivery

import (
	"github.com/2020_1_Skycode/internal/middlewares"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/2020_1_Skycode/internal/users"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/renstrom/shortuuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

type UserHandler struct {
	userUseCase users.UseCase
	middlewareC *middlewares.MWController
}

func NewUserHandler(router *gin.Engine, uUC users.UseCase, middlewareC *middlewares.MWController) *UserHandler {
	uh := &UserHandler{
		userUseCase: uUC,
		middlewareC: middlewareC,
	}

	router.GET("api/v1/profile", middlewareC.CheckAuth(), uh.GetProfile())
	router.POST("api/v1/signup", uh.SignUp())
	router.PUT("api/v1/profile/bio", middlewareC.CheckAuth(), uh.EditBio())
	router.PUT("api/v1/profile/avatar", middlewareC.CheckAuth(), uh.EditAvatar())
	router.PUT("api/v1/profile/password", middlewareC.CheckAuth(), uh.ChangePassword())
	router.PUT("api/v1/profile/phone", middlewareC.CheckAuth(), uh.ChangePhoneNumber())

	return uh
}

func (uh *UserHandler) SignUp() gin.HandlerFunc {
	type SignUpRequest struct {
		FirstName string `json:"firstName, omitempty" binding:"required"`
		LastName  string `json:"lastName, omitempty" binding:"required"`
		Phone     string `json:"phone, omitempty" binding:"required"`
		Password  string `json:"password, omitempty" binding:"required"`
	}

	return func(c *gin.Context) {
		req := &SignUpRequest{}

		if err := c.Bind(req); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		u := &models.User{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Phone:     req.Phone,
			Password:  req.Password,
		}

		if err := uh.userUseCase.CreateUser(u); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, tools.Message{
			Message: "User has been registered",
		})
	}
}

func (uh *UserHandler) EditBio() gin.HandlerFunc {
	type EditBioRequest struct {
		FirstName string `json:"firstName" binding:"required"`
		LastName  string `json:"lastName" binding:"required"`
		Email     string `json:"email" binding:"required"`
	}

	return func(c *gin.Context) {
		updProfile := &EditBioRequest{}

		if err := c.Bind(updProfile); err != nil {
			logrus.Info("Bad params")
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		usr, exists := c.Get("user")

		if !exists {
			c.JSON(http.StatusUnauthorized, tools.Error{
				ErrorMessage: tools.Unauthorized.Error(),
			})

			return
		}

		user, ok := usr.(*models.User)

		if !ok {
			c.JSON(http.StatusInternalServerError, tools.Error{
				ErrorMessage: tools.UserTypeAssertionErr.Error(),
			})
			return
		}

		user.FirstName = updProfile.FirstName
		user.LastName = updProfile.LastName
		user.Email = updProfile.Email

		if err := uh.userUseCase.UpdateBio(user); err != nil {
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

func (uh *UserHandler) EditAvatar() gin.HandlerFunc {
	return func(c *gin.Context) {
		usr, exists := c.Get("user")

		if !exists {
			c.JSON(http.StatusUnauthorized, tools.Error{
				ErrorMessage: tools.Unauthorized.Error(),
			})

			return
		}

		user, ok := usr.(*models.User)

		if !ok {
			c.JSON(http.StatusInternalServerError, tools.Error{
				ErrorMessage: tools.UserTypeAssertionErr.Error(),
			})
			return
		}

		file, err := c.FormFile("avatar")

		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		filename := shortuuid.New() + "-" + file.Filename

		if err := c.SaveUploadedFile(file, tools.AvatarPath+filename); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		if user.Avatar != "" {
			if err := os.Remove(tools.AvatarPath + user.Avatar); err != nil {
				logrus.Info(err)
				c.JSON(http.StatusInternalServerError, tools.Error{
					ErrorMessage: tools.DeleteAvatarError.Error(),
				})

				return
			}
		}

		if err := uh.userUseCase.UpdateAvatar(user.ID, filename); err != nil {
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

func (uh *UserHandler) ChangePhoneNumber() gin.HandlerFunc {
	type ChangePhoneNumberRequest struct {
		NewPhone string `json:"newPhone" binding:"required" validate:"numeric"`
	}
	return func(c *gin.Context) {
		req := &ChangePhoneNumberRequest{}

		if err := c.Bind(req); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		usr, exists := c.Get("user")

		if !exists {
			c.JSON(http.StatusUnauthorized, tools.Error{
				ErrorMessage: tools.Unauthorized.Error(),
			})

			return
		}

		user, ok := usr.(*models.User)

		if !ok {
			c.JSON(http.StatusInternalServerError, tools.Error{
				ErrorMessage: tools.UserTypeAssertionErr.Error(),
			})
			return
		}

		validate := validator.New()

		if err := validate.Struct(req); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		if err := uh.userUseCase.UpdatePhoneNumber(user.ID, req.NewPhone); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusInternalServerError, tools.Error{
				ErrorMessage: tools.UpdatePhoneError.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, tools.Message{
			Message: "success",
		})
	}
}

func (uh *UserHandler) ChangePassword() gin.HandlerFunc {
	type ChangePasswordRequest struct {
		NewPassword string `json:"newPassword" binding:"required"`
	}
	return func(c *gin.Context) {
		req := &ChangePasswordRequest{}

		if err := c.Bind(req); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		usr, exists := c.Get("user")

		if !exists {
			c.JSON(http.StatusUnauthorized, tools.Error{
				ErrorMessage: tools.Unauthorized.Error(),
			})

			return
		}

		user, ok := usr.(*models.User)

		if !ok {
			c.JSON(http.StatusInternalServerError, tools.Error{
				ErrorMessage: tools.UserTypeAssertionErr.Error(),
			})
			return
		}

		if err := uh.userUseCase.UpdatePassword(user.ID, req.NewPassword); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusInternalServerError, tools.Error{
				ErrorMessage: tools.UpdatePhoneError.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, tools.Message{
			Message: "success",
		})
	}
}

func (uh *UserHandler) GetProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		usr, exists := c.Get("user")

		if !exists {
			c.JSON(http.StatusUnauthorized, tools.Error{
				ErrorMessage: tools.Unauthorized.Error(),
			})

			return
		}

		user, ok := usr.(*models.User)

		if !ok {
			c.JSON(http.StatusInternalServerError, tools.Error{
				ErrorMessage: tools.UserTypeAssertionErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, tools.UserMessage{
			User: user,
		})
	}
}


