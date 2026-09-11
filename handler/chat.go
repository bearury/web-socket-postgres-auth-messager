package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (handler *Handler) addToChat(c *gin.Context) {

}

func (handler *Handler) getChats(c *gin.Context) {
	userId, err := getUserId(c)

	if err != nil {
		logrus.Error(err)
	}

	chats, err := handler.services.Chat.GetChats(userId)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, chats)
}

func (handler *Handler) removeChat(c *gin.Context) {

}
