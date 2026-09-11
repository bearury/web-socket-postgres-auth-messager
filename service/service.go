package service

import (
	"web-socket-postgres-auth-messager/entities"
	"web-socket-postgres-auth-messager/repository"
)

type Authorization interface {
	CreateUser(user entities.User) (string, error)
	GenerateToken(username string, password string) (string, error)
	ParseToken(token string) (string, error)
	GetUserById(userId string) (entities.UserResponse, error)
}

type Chat interface {
	GetChats(userId string) ([]entities.Chat, error)
}

type Service struct {
	Authorization
	Chat
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.AuthRepository),
		Chat:          NewChatService(repo.ChatRepository),
	}
}
