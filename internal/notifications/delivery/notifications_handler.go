package delivery

import (
	"github.com/2020_1_Skycode/internal/middlewares"
	"github.com/2020_1_Skycode/internal/notifications"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/2020_1_Skycode/internal/tools/notificationsWS"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type NotificationsHandler struct {
	notificationsUC notifications.UseCase
	middlewareC     *middlewares.MWController
	noteServer      *notificationsWS.NotificationServer
}

func NewNotificationsHandler(private *gin.RouterGroup, public *gin.RouterGroup, nUC notifications.UseCase,
	mw *middlewares.MWController, ns *notificationsWS.NotificationServer) *NotificationsHandler {
	nh := &NotificationsHandler{
		notificationsUC: nUC,
		middlewareC:     mw,
		noteServer:      ns,
	}

	public.GET("/notifications", nh.GeUserNotifications())
	private.POST("/notifications/:id", nh.ChangeReadStatus())

	public.GET("/notification_server", nh.JoinNotificationServer())
	private.DELETE("/notifications/:id", nh.DeleteNotification())

	return nh
}

func (nh *NotificationsHandler) GeUserNotifications() gin.HandlerFunc {
	return func(c *gin.Context) {
		u, err := nh.middlewareC.GetUser(c)

		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		notes, err := nh.notificationsUC.GetAllByUser(u.ID)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"notifications": notes,
		})
	}
}

func (nh *NotificationsHandler) ChangeReadStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		u, err := nh.middlewareC.GetUser(c)

		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})
		}

		if err := nh.notificationsUC.ChangeReadStatus(noteID, u.ID); err != nil {
			if err == tools.PermissionError {
				c.JSON(http.StatusForbidden, tools.Error{
					ErrorMessage: err.Error(),
				})

				return
			}

			if err == tools.NotificationNotFound {
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

		c.JSON(http.StatusOK, tools.Message{Message: "Success"})
	}
}

func (nh *NotificationsHandler) DeleteNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
		u, err := nh.middlewareC.GetUser(c)

		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})
		}

		if err := nh.notificationsUC.Delete(noteID, u.ID); err != nil {
			if err == tools.PermissionError {
				c.JSON(http.StatusForbidden, tools.Error{
					ErrorMessage: err.Error(),
				})

				return
			}

			if err == tools.NotificationNotFound {
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

		c.JSON(http.StatusOK, tools.Message{Message: "Success"})
	}
}

func (nh *NotificationsHandler) JoinNotificationServer() gin.HandlerFunc {
	return func(c *gin.Context) {
		u, err := nh.middlewareC.GetUser(c)

		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		_, err = nh.noteServer.CreateClient(c.Writer, c.Request, u.ID)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}
	}
}
