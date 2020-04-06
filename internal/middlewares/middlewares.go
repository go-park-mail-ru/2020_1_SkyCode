package middlewares

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/sessions"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/2020_1_Skycode/internal/users"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type MWController struct {
	sessionUC sessions.UseCase
	userUC users.UseCase
}

func NewMiddleWareController(router *gin.Engine, sessionUC sessions.UseCase, userUC users.UseCase) *MWController {
	mw := &MWController{
		sessionUC: sessionUC,
		userUC: userUC,
	}

	router.Use(mw.AccessLogging())
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
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if c.Request.Method == http.MethodOptions {
			c.JSON(http.StatusOK, tools.Message{
				Message: "Options ok",
			})
			return
		}

		c.Next()
	}
}

func (mw *MWController) CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.Info("Check auth")
		cookie, err := c.Cookie("SkyDelivery");

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

func (mw *MWController) AccessLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := []string{c.Request.Method, c.Request.URL.String(), c.Request.RemoteAddr, time.Now().UTC().String()}
		logrus.Info(strings.Join(data, " "))
		c.Next()
	}
}
