package delivery

import (
	"github.com/2020_1_Skycode/internal/middlewares"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/sessions"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/2020_1_Skycode/internal/tools/CSRFManager"
	"github.com/2020_1_Skycode/internal/tools/requestValidator"
	"github.com/2020_1_Skycode/internal/users"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type SessionHandler struct {
	SessionUseCase sessions.UseCase
	UserUseCase    users.UseCase
	MiddlewareC    *middlewares.MWController
	tM             *CSRFManager.CSRFManager
	v              *requestValidator.RequestValidator
}

func NewSessionHandler(private *gin.RouterGroup, public *gin.RouterGroup, sessionUC sessions.UseCase, usersUC users.UseCase,
	validator *requestValidator.RequestValidator, tM *CSRFManager.CSRFManager, mwareC *middlewares.MWController) *SessionHandler {
	sh := &SessionHandler{
		SessionUseCase: sessionUC,
		UserUseCase:    usersUC,
		MiddlewareC:    mwareC,
		tM:             tM,
		v:              validator,
	}

	public.POST("/signin", sh.SignIn())

	private.POST("/logout", sh.LogOut())

	return sh
}

type signInRequest struct {
	Phone    string `json:"phone" binding:"required" validate:"min=11,max=15"`
	Password string `json:"password" binding:"required" validate:"passwd"`
}

//@Tags Session
//@Summary Sign In Route
//@Description Signing in user
//@Accept json
//@Produce json
//@Param SignInReq body signInRequest true "User data"
//@Success 200 object models.User
//@Failure 400 object tools.Error
//@Failure 404 object tools.Error
//@Router /signin [post]
func (sh *SessionHandler) SignIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &signInRequest{}

		if err := c.Bind(req); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		errorsList := sh.v.ValidateRequest(req)

		if len(*errorsList) > 0 {
			logrus.Info(tools.NotRequiredFields)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.NotRequiredFields.Error(),
			})

			return
		}

		user, err := sh.UserUseCase.GetUserByPhone(req.Phone)

		if err != nil {
			logrus.Info(err, req)
			c.JSON(http.StatusNotFound, tools.Error{
				ErrorMessage: tools.NoSuchUser.Error(),
			})

			return
		}

		session, err := sh.MiddlewareC.GetSession(c)

		if err == nil && session.UserId == user.ID {
			logrus.Info(tools.HashingError)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.Authorized.Error(),
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

		csrfToken, err := sh.tM.GenerateCSRF(session.UserId, session.Token)

		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusInternalServerError, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		c.Writer.Header().Set("X-Csrf-Token", csrfToken)

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

//@Tags Session
//@Summary Logout Route
//@Description Logouting user
//@Accept json
//@Produce json
//@Success 200 object tools.Message
//@Failure 401 object tools.Error
//@Failure 500 object tools.Error
//@Security basicAuth
//@Router /logout [post]
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
