package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (handler *Handler) signUp(c *gin.Context) {
	logrus.Info("sign up handler")
}

func (handler *Handler) signIn(c *gin.Context) {

}
