package handler

import (
	"fmt"
	"net/http"
	"web-socket-postgres-auth-messager/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (handler *Handler) getMe(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		logrus.Error(err)
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("%w: %s", service.ErrInvalidUserID, err.Error()))
		return
	}

	user, err := handler.services.GetUserById(userId)

	if err != nil {
		logrus.Error(err)
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error get user: %s", err.Error()))
	}

	c.JSON(http.StatusOK, user)
}

func (handler *Handler) getUsers(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		logrus.Error(err)
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("%w: %s", service.ErrInvalidUserID, err.Error()))
		return
	}

	users, err := handler.services.GetAllUsersWithoutUser(userId)

	if err != nil {
		logrus.Error(err)
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Error get users: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, users)
}
