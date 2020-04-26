package middlewares

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/sessions"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/2020_1_Skycode/internal/tools/CSRFManager"
	"github.com/2020_1_Skycode/internal/users"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type MWController struct {
	sessionUC sessions.UseCase
	userUC    users.UseCase
	cM        *CSRFManager.CSRFManager
}

func NewMiddleWareController(router *gin.Engine, sessionUC sessions.UseCase,
	userUC users.UseCase, cM *CSRFManager.CSRFManager) *MWController {
	mw := &MWController{
		sessionUC: sessionUC,
		userUC:    userUC,
		cM:        cM,
	}

	router.Use(mw.AccessLogging())
	router.Use(mw.CheckAuth())
	router.Use(mw.CORS())
	return mw
}

func (mw *MWController) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.Info("CORS handler")
		origin := c.Request.Header.Get("Origin")

		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, DELETE, PUT")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Csrf-Token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "X-Csrf-Token")

		if c.Request.Method == http.MethodOptions {
			logrus.Info("OPTIONS")
			c.JSON(http.StatusOK, tools.Message{
				Message: "Options ok",
			})

			c.Abort()
			return
		}

		c.Next()
	}
}

func (mw *MWController) CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.Info("Check auth")
		cookie, err := c.Cookie("SkyDelivery")

		if err != nil {
			logrus.Info(err)
			c.Next()
			return
		}

		sess, err := mw.sessionUC.GetSession(cookie)

		if err != nil {
			logrus.Info(err)
			c.Next()
			return
		}

		user, err := mw.userUC.GetUserById(sess.UserId)

		if err != nil {
			logrus.Info(err)
			c.Next()
			return
		}

		c.Set("session", sess)
		c.Set("user", user)
		c.Next()
	}
}

func (mw *MWController) GetUser(c *gin.Context) (*models.User, error) {
	usr, exists := c.Get("user")
	if !exists {
		return nil, tools.Unauthorized
	}

	user, ok := usr.(*models.User)
	if !ok {
		return nil, tools.UserTypeAssertionErr
	}

	return user, nil
}

func (mw *MWController) GetSession(c *gin.Context) (*models.Session, error) {
	sess, exists := c.Get("session")

	if !exists {
		return nil, tools.Unauthorized
	}

	session, ok := sess.(*models.Session)

	if !ok {
		return nil, tools.SessionTypeAssertionErr
	}

	return session, nil
}

func (mw *MWController) CSRFControl() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("X-Csrf-Token")
		c.Writer.Header().Set("X-XSS-Protection", "1")

		if token == "" {
			logrus.Info(tools.CSRFNotPresented.Error())
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.CSRFNotPresented.Error(),
			})

			c.Abort()
			return
		}

		session, err := mw.GetSession(c)

		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusUnauthorized, tools.Error{
				ErrorMessage: tools.Unauthorized.Error(),
			})

			c.Abort()
			return
		}

		err = mw.cM.ValidateCSRF(token, session.UserId, session.Token)

		if err != nil {
			if err == tools.ExpiredCSRFError {
				newToken, err := mw.cM.GenerateCSRF(session.UserId, session.Token)

				if err != nil {
					logrus.Info(err)
					c.JSON(http.StatusInternalServerError, tools.Error{
						ErrorMessage: err.Error(),
					})

					c.Abort()
					return
				}

				c.Writer.Header().Set("X-Csrf-Token", newToken)

				c.Next()
				return
			}

			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			c.Abort()
			return
		}

		c.Next()
	}
}

func (mw *MWController) AccessLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := []string{c.Request.Method, c.Request.URL.String(), c.Request.RemoteAddr, time.Now().UTC().String()}
		logrus.Info(strings.Join(data, " "), c.HandlerNames())
		c.Next()
	}
}
