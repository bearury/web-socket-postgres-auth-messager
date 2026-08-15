package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (handler *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "No authorization header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header")
		return
	}

	userId, err := handler.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header")
		return
	}

	c.Set(userCtx, userId)
}

func getUserId(c *gin.Context) (string, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "User not found in context")
		return "", errors.New("User not found in context")
	}

	idString, ok := id.(string)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "Userid is invalid type")
		return "", errors.New("Userid is invalid type")
	}

	return idString, nil
}

func getParam(c *gin.Context) (string, error) {
	id := c.Param("id")

	if id == "" {
		return "", errors.New("Param id is empty")
	}

	// TODO надо подключить пакет UUID для валидации uuid
	//if err := uuid.Validate(id); err == nil {
	//	return id, nil
	//} else {
	//	return "", errors.New("Param id is not uuid")
	//}

	return id, nil
}
