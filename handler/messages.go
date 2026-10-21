package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (handler *Handler) getMessages(c *gin.Context) {
	chatIDStr, err := getParam(c)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid chat id format: must be a valid UUID")
		return
	}

	messages, err := handler.services.Messages.GetMessages(chatIDStr)

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
