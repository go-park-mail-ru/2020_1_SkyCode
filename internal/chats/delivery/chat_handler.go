package delivery

import (
	"github.com/2020_1_Skycode/internal/chats"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ChatHandler struct {
	cU chats.UseCase
}

func NewChatsHandler(public *gin.RouterGroup, cU chats.UseCase) *ChatHandler {
	cH := &ChatHandler{
		cU: cU,
	}

	public.GET("/chat", cH.StartUserChat())

	return cH
}

func (cH *ChatHandler) StartUserChat() gin.HandlerFunc {
	return func (c *gin.Context) {
		chatID := cH.cU.StartChat()

		ws, err := cH.cU.JoinUserToChat(c.Writer, c.Request, chatID)

		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusInternalServerError, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		for {
			message, err := cH.cU.ReadMessageFromUSer(ws)
			logrus.Info(message, err)

			if err != nil {
				logrus.Error(err)
				if err := cH.cU.LeaveUserChat(chatID); err != nil {
					logrus.Error(err)
				}
				break
			}

			cH.cU.WriteFromUserMessage(message)
		}

		c.JSON(http.StatusOK, tools.Message{
			Message: err.Error(),
		})

		return
	}
}