package handler

import (
	"errors"
	"fmt"
	"net/http"
	"web-socket-postgres-auth-messager/dto"
	"web-socket-postgres-auth-messager/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (handler *Handler) addToChat(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		logrus.Error(err)
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("%w: %s", service.ErrInvalidUserID, err.Error()))
		return
	}

	var input dto.CreateChatRequest

	logrus.Info("Запрос от юзера с ID: ", userId)

	if err := c.BindJSON(&input); err != nil {
		logrus.Errorf("Error binding JSON: %s, Handler: %s, Body: %s", err, c.Request.RequestURI, c.Request.Body)
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Неверные параметры запроса: %s", err.Error()))
		return
	}

	chat, err := handler.services.Chat.AddChat(input.Name, input.Type, userId, input.UserIds)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrChatAlreadyExists):
			newErrorResponse(c, http.StatusConflict, err.Error())

		case errors.Is(err, service.ErrUserNotFound):
			newErrorResponse(c, http.StatusNotFound, err.Error())

		case errors.Is(err, service.ErrUserIDsEmpty),
			errors.Is(err, service.ErrInvalidUserID),
			errors.Is(err, service.ErrInvalidChatType):
			newErrorResponse(c, http.StatusBadRequest, err.Error())

		default:
			logrus.Error(err)
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}

		return
	}

	c.JSON(http.StatusOK, chat)
}

func (handler *Handler) getChats(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		logrus.Error(err)
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("%w: %s", service.ErrInvalidUserID, err.Error()))
		return
	}

	logrus.Info("Запрос от юзера с ID: ", userId)

	chats, err := handler.services.Chat.GetChats(userId)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, chats)
}

func (handler *Handler) removeChat(c *gin.Context) {

}
