package handler

import (
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

	logrus.Infof("getChats: %v", userId)

}

func (handler *Handler) removeChat(c *gin.Context) {

}
