package dto

import "web-socket-postgres-auth-messager/entities"

type CreateChatRequest struct {
	Type    entities.ChatType `json:"type" binding:"required"`
	Name    string            `json:"name"`
	UserIds []string          `json:"userIds" binding:"required"`
}
