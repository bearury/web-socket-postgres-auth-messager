package handler

import (
	"fmt"
	"net/http"
	"web-socket-postgres-auth-messager/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (handler *Handler) getMessages(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		logrus.Error(err)
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("%w: %s", service.ErrInvalidUserID, err.Error()))
		return
	}

	chatIDStr, err := getParam(c)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid chat id format: must be a valid UUID")
		return
	}

	ctx := c.Request.Context()

	messages, err := handler.services.Messages.GetMessages(ctx, chatIDStr, userId)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, messages)
}

func (handler *Handler) addMessage(c *gin.Context) {
	newErrorResponse(c, http.StatusBadRequest, "Route for internal implementation")
}

func (handler *Handler) updateMessage(c *gin.Context) {}

func (handler *Handler) removeMessage(c *gin.Context) {}
