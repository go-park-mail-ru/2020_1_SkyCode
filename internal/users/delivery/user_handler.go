package delivery

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/2020_1_Skycode/internal/users"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserHandler struct {
	userUseCase users.UseCase
}

func NewUserHandler(router *gin.Engine, uUC users.UseCase) *UserHandler {
	uh := &UserHandler{
		userUseCase: uUC,
	}

	router.POST("api/v1/signup", uh.SignUp())

	return uh
}

func (uh *UserHandler) SignUp() gin.HandlerFunc {
	type SignUpRequest struct {
		FirstName string `json:"first_name, omitempty" binding:"required"`
		LastName  string `json:"last_name, omitempty" binding:"required"`
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
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"first_name" binding:"required"`
		Email     string `json:"first_name" binding:"required"`
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

	}
}

func (uh *UserHandler) EditAvatar() gin.HandlerFunc {
	return nil
}

func (uh *UserHandler) ChangePhoneNumber() gin.HandlerFunc {
	return nil
}

func (uh *UserHandler) ChangePassword() gin.HandlerFunc {
	return nil
}
