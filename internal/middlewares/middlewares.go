package middlewares

import (
	"github.com/2020_1_Skycode/internal/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type MWController struct {
	sessionUC sessions.UseCase
}

func NewMiddleWareController(router *gin.Engine, sessionUC sessions.UseCase) *MWController {
	mw := &MWController{
		sessionUC: sessionUC,
	}

	router.Use(mw.CORS())

	return mw
}

func (mw *MWController) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, DELETE, PUT")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if c.Request.Method == http.MethodOptions {
			return
		}

		c.Next()
	}
}

func (mw *MWController) CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		c.Set("session", sess)
		c.Next()
	}
}

func (mw *MWController) AccessLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		data := []string{r.Method, r.URL.String(), r.RemoteAddr, time.Now().UTC().String()}
		logrus.Info(strings.Join(data, " "))
	})
}
