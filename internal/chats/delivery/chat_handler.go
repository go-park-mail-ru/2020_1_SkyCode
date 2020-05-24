package delivery

import (
	"github.com/2020_1_Skycode/internal/chats"
	"github.com/2020_1_Skycode/internal/middlewares"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/2020_1_Skycode/internal/tools/supportChat"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type ChatHandler struct {
	cU chats.UseCase
	sC *supportChat.ChatServer
	mw *middlewares.MWController
}

func NewChatsHandler(private *gin.RouterGroup, public *gin.RouterGroup,
	cU chats.UseCase, mw *middlewares.MWController) *ChatHandler {
	cH := &ChatHandler{
		cU: cU,
		mw: mw,
		sC: supportChat.NewChatServer(),
	}

	public.GET("/chat", cH.StartUserChat())
	private.GET("/chats", cH.GetSupChatList())
	public.GET("/chats/:chatID/join", cH.JoinSupport())
	private.GET("/chats/:chatID/details", cH.GetChatMessages())
	private.DELETE("/chats/:chatID", cH.CloseChat())

	return cH
}

func (cH *ChatHandler) StartUserChat() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := cH.mw.GetUser(c)

		if user == nil {
			logrus.Info(err)
			c.JSON(http.StatusUnauthorized, tools.Error{
				ErrorMessage: tools.Unauthorized.Error(),
			})

			return
		}

		conn, chat, err := cH.sC.JoinUserToChat(user, c.Writer, c.Request)

		if err != nil {
			logrus.Error(err)
			return
		}

		for {
			message := &supportChat.InputMessage{}

			if err := conn.ReadJSON(message); err != nil {
				logrus.Error(err)
				break
			}
			logrus.Debug(message, err)

			if err := cH.cU.StoreMessage(&models.ChatMessage{
				UserID:   user.ID,
				UserName: user.FirstName,
				ChatID:   message.ChatID,
				Message:  message.Message,
			}); err != nil {
				logrus.Error(err)
			}

			chat.NotifyMembers(message)
		}
	}
}

func (cH *ChatHandler) GetSupChatList() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := cH.mw.GetUser(c)

		if user == nil || !user.IsSupport() {
			logrus.Info(err)
			c.JSON(http.StatusUnauthorized, tools.Error{
				ErrorMessage: tools.Unauthorized.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"chats": cH.sC.GetSupportChats(),
		})
	}
}

func (cH *ChatHandler) JoinSupport() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := cH.mw.GetUser(c)

		if user == nil || !user.IsSupport() {
			logrus.Info(err)
			c.JSON(http.StatusUnauthorized, tools.Error{
				ErrorMessage: tools.Unauthorized.Error(),
			})

			return
		}

		chatID, err := strconv.ParseUint(c.Param("chatID"), 10, 64)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		conn, chat, err := cH.sC.JoinSupportToChat(chatID, user, c.Writer, c.Request)
		if err != nil {
			if err == tools.ChatNotFound {
				c.JSON(http.StatusNotFound, tools.Error{
					ErrorMessage: tools.ChatNotFound.Error(),
				})
				return
			}
			if err == tools.NewSupportJoinError {
				c.JSON(http.StatusConflict, tools.Error{
					ErrorMessage: err.Error(),
				})
			}

			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		for {
			message := &supportChat.InputMessage{}

			if err := conn.ReadJSON(message); err != nil {
				logrus.Error(err)
				break
			}
			logrus.Debug(message, err)

			if err := cH.cU.StoreMessage(&models.ChatMessage{
				UserName: user.FirstName,
				UserID:   user.ID,
				ChatID:   message.ChatID,
				Message:  message.Message,
			}); err != nil {
				logrus.Error(err)
			}

			chat.NotifyMembers(message)
		}

	}
}

func (cH *ChatHandler) GetChatMessages() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := cH.mw.GetUser(c)

		if user == nil {
			logrus.Info(err)
			c.JSON(http.StatusUnauthorized, tools.Error{
				ErrorMessage: tools.Unauthorized.Error(),
			})

			return
		}

		chatID, err := strconv.ParseUint(c.Param("chatID"), 10, 64)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		chat := cH.sC.GetChat(chatID)

		if chat == nil {
			c.JSON(http.StatusNotFound, tools.Error{
				ErrorMessage: tools.ChatNotFound.Error(),
			})
			return
		}

		messages, err := cH.cU.GetChatMessages(chatID)

		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusInternalServerError, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, messages)
	}
}

func (cH *ChatHandler) CloseChat() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := cH.mw.GetUser(c)

		if user == nil || !user.IsSupport() {
			logrus.Info(err)
			c.JSON(http.StatusUnauthorized, tools.Error{
				ErrorMessage: tools.Unauthorized.Error(),
			})

			return
		}

		chatID, err := strconv.ParseUint(c.Param("chatID"), 10, 64)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		if err := cH.sC.DeleteChat(chatID); err != nil {
			if err == tools.ChatNotFound {
				c.JSON(http.StatusNotFound, tools.Error{
					ErrorMessage: tools.ChatNotFound.Error(),
				})

				return
			}

			if err == tools.ChatInUse {
				c.JSON(http.StatusBadRequest, tools.Error{
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
	}
}
