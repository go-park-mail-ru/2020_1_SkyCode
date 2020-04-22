package delivery

import (
	"github.com/2020_1_Skycode/internal/chats"
	"github.com/2020_1_Skycode/internal/middlewares"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ChatHandler struct {
	cU chats.UseCase
	mw *middlewares.MWController
}

func NewChatsHandler(private *gin.RouterGroup, public *gin.RouterGroup,
	cU chats.UseCase, mw *middlewares.MWController) *ChatHandler {
	cH := &ChatHandler{
		cU: cU,
		mw: mw,
	}

	public.GET("/chat", cH.StartUserChat())
	private.GET("/chats", cH.GetSupChatList())
	public.GET("/chats/:chatID/join", cH.JoinSupport())

	return cH
}

func (cH *ChatHandler) StartUserChat() gin.HandlerFunc {
	return func (c *gin.Context) {
		user, err := cH.mw.GetUser(c)

		if user == nil {
			logrus.Info(err)
			c.JSON(http.StatusUnauthorized, tools.Error{
				ErrorMessage: tools.Unauthorized.Error(),
			})

			return
		}

		conn, joinMsg, err := cH.cU.StartChat(c.Writer, c.Request)

		if err != nil {
			logrus.Error(err)
		}

		err = cH.cU.JoinUserToChat(conn, joinMsg.FullName, joinMsg.ChatID)

		if err != nil {
			logrus.Error(err)
			return
		}

		for {
			message, err := cH.cU.ReadMessageFromUSer(conn)
			logrus.Info(message, err)

			if err != nil {
				logrus.Error(err)
				if err := cH.cU.LeaveUserChat(joinMsg.ChatID); err != nil {
					logrus.Error(err)
				}
				break
			}

			if err := cH.cU.StoreMessage(&models.ChatMessage{
				UserID: user.ID,
				ChatID: message.ChatID,
				Message: message.Message,
			}); err != nil {
				logrus.Error(err)
			}

			cH.cU.WriteFromUserMessage(message)
		}

		return
	}
}

func (cH *ChatHandler) GetSupChatList() gin.HandlerFunc {
	return func (c *gin.Context) {
		user, err := cH.mw.GetUser(c)

		if user == nil || !user.IsSupport() {
			logrus.Info(err)
			c.JSON(http.StatusUnauthorized, tools.Error{
				ErrorMessage: tools.Unauthorized.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, cH.cU.GetChats())
		return
	}
}

func (cH *ChatHandler) JoinSupport() gin.HandlerFunc {
	return func (c *gin.Context) {
		user, err := cH.mw.GetUser(c)

		if user == nil || !user.IsSupport() {
			logrus.Info(err)
			c.JSON(http.StatusUnauthorized, tools.Error{
				ErrorMessage: tools.Unauthorized.Error(),
			})

			return
		}

		chatID := c.Param("chatID")

		if chatID == "" {
			logrus.Info(err)
			return
		}

		conn, joinMsg, err := cH.cU.FindChat(c.Writer, c.Request, chatID)

		if err != nil {
			logrus.Error(err)
			return
		}

		err = cH.cU.JoinSupportToChat(conn, joinMsg.FullName, joinMsg.ChatID)

		if err != nil {
			logrus.Error(err)

			return
		}

		for {
			message, err := cH.cU.ReadMessageFromUSer(conn)
			logrus.Info(message, err)

			if err != nil {
				logrus.Error(err)
				if err := cH.cU.LeaveSupportChat(joinMsg.ChatID); err != nil {
					logrus.Error(err)
				}
				break
			}

			if err := cH.cU.StoreMessage(&models.ChatMessage{
				UserID: user.ID,
				ChatID: message.ChatID,
				Message: message.Message,
			}); err != nil {
				logrus.Error(err)
			}

			cH.cU.WriteFromUserMessage(message)
		}

		return
	}
}

