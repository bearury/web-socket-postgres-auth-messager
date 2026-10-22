package service

import (
	"context"
	"web-socket-postgres-auth-messager/entities"
	"web-socket-postgres-auth-messager/repository"
)

type Authorization interface {
	CreateUser(user entities.User) (entities.User, error)
	GenerateToken(username string, password string) (string, error)
	ParseToken(token string) (string, error)
	GetUserById(userId string) (entities.UserResponse, error)
	GetAllUsersWithoutUser(userId string) ([]entities.UserResponse, error)
}

type Chat interface {
	GetChats(userId string) ([]entities.Chat, error)
	AddChat(name string, chatType entities.ChatType, userId string, userIds []string) (entities.Chat, error)
	GetUsersChatByIdWithoutUser(chatId string, userId string) ([]entities.ChatUsers, error)
}

type Messages interface {
	GetMessages(ctx context.Context, chatId, userId string) ([]entities.MessageResponseWithStatus, error)
	AddMessage(user entities.UserResponse, chatId, content string) (entities.MessageResponse, error)
}

type NotificationHub interface {
	SendToChat(userIDs []string, message entities.WSMessageResponse)
}

type Service struct {
	Authorization
	Chat
	Messages
}

func NewService(repo *repository.Repository, hub NotificationHub) *Service {
	return &Service{
		Authorization: NewAuthService(repo.AuthRepository),
		Chat:          NewChatService(repo.ChatRepository, repo.AuthRepository),
		Messages:      NewMessagesService(repo.MessagesRepository, repo.ChatRepository, hub),
	}
}
