package delivery

import (
	"github.com/2020_1_Skycode/internal/middlewares"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/sessions"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/2020_1_Skycode/internal/users"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type SessionHandler struct {
	SessionUseCase sessions.UseCase
	UserUseCase    users.UseCase
	MiddlewareC    *middlewares.MWController
}

func NewSessionHandler(router *gin.Engine, sessionUC sessions.UseCase, usersUC users.UseCase, mwareC *middlewares.MWController) *SessionHandler {
	sh := &SessionHandler{
		SessionUseCase: sessionUC,
		UserUseCase:    usersUC,
		MiddlewareC:    mwareC,
	}

	router.POST("api/v1/signin", sh.SignIn())
	router.POST("api/v1/logout", sh.MiddlewareC.CheckAuth(), sh.LogOut())

	return sh
}

func (sh *SessionHandler) SignIn() gin.HandlerFunc {
	type SignInRequest struct {
		Phone    string `json:"phone" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	return func(c *gin.Context) {
		req := &SignInRequest{}

		if err := c.Bind(req); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		user, err := sh.UserUseCase.GetUserByPhone(req.Phone)

		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusNotFound, tools.Error{
				ErrorMessage: tools.NoSuchUser.Error(),
			})

			return
		}

		if user.Password != req.Password {
			c.JSON(http.StatusNotFound, tools.Error{
				ErrorMessage: tools.WrongPassword.Error(),
			})

			return
		}

		session, cookie := models.GenerateSession(user.ID)

		if err := sh.SessionUseCase.StoreSession(session); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusNotFound, tools.Error{
				ErrorMessage: tools.SessionStoreError.Error(),
			})

			return
		}

		c.SetCookie(cookie.Name,
			cookie.Value,
			cookie.MaxAge,
			cookie.Path,
			cookie.Domain,
			cookie.Secure,
			cookie.HttpOnly)

		c.JSON(http.StatusOK, tools.UserMessage{
			User: user,
		})
	}
}

func (sh *SessionHandler) LogOut() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess, exists := c.Get("session")

		if !exists {
			c.JSON(http.StatusUnauthorized, tools.Error{
				ErrorMessage: tools.Unauthorized.Error(),
			})

			return
		}

		session, ok := sess.(*models.Session)

		if !ok {
			c.JSON(http.StatusInternalServerError, tools.Error{
				ErrorMessage: tools.SessionTypeAssertionErr.Error(),
			})

			return
		}

		if err := sh.SessionUseCase.DeleteSession(session.ID); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusInternalServerError, tools.Error{
				ErrorMessage: tools.DeleteSessionError.Error(),
			})

			return
		}

		cookie := &http.Cookie{
			Name:     "SkyDelivery",
			Value:    session.Token,
			MaxAge:   -1,
			HttpOnly: true,
		}

		c.SetCookie(cookie.Name,
			cookie.Value,
			cookie.MaxAge,
			cookie.Path,
			cookie.Domain,
			cookie.Secure,
			cookie.HttpOnly)

		c.JSON(http.StatusOK, tools.Message{
			Message: "success",
		})
	}
}
