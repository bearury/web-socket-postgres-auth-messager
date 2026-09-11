package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"web-socket-postgres-auth-messager/entities"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (handler *Handler) signUp(c *gin.Context) {
	var input entities.User

	if err := c.BindJSON(&input); err != nil {
		logrus.Errorf("Error binding JSON: %s", err)
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Неверные параметры запроса: %s", err.Error()))
		return
	}

	id, err := handler.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (handler *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		logrus.Errorf("Handler: %s, Params: %b,Error binding JSON: %v", c.Request.RequestURI, c.Request.GetBody, err)
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Неверные параметры запроса: %s", err.Error()))
		return
	}

	token, err := handler.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Errorf("Handler: %s, UserName: %s, Password: %s, Error generate Token: %s", c.Request.RequestURI, input.Username, input.Password, err)
			newErrorResponse(c, http.StatusNotFound, "Неверный логин или пароль")
			return
		} else {
			logrus.Errorf("Error generate Token: %s", err)
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})

}
