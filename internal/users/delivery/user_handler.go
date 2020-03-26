package delivery

import (
	"fmt"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/2020_1_Skycode/internal/users"
	"github.com/gin-gonic/gin"
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
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
		Password  string `json:"password"`
	}

	return func(c *gin.Context) {
		req := &SignUpRequest{}

		fmt.Println("SignUp")

		if err := c.Bind(req); err != nil {
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		u := &models.User{
			FirstName: req.FirstName,
			LastName: req.LastName,
			Phone: req.Phone,
			Password: req.Password,
		}

		if err := uh.userUseCase.CreateUser(u); err != nil {
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
