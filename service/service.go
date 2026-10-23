package service

import (
	"context"
	"web-socket-postgres-auth-messager/entities"
	ws_hub "web-socket-postgres-auth-messager/handler/ws-hub"
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

type WebSocket interface {
	HandleMessage(client *ws_hub.Client, msg *entities.WSMessage) error
}

type NotificationHub interface {
	SendToChat(userIDs []string, message entities.WSMessage)
}

type Service struct {
	Authorization
	Chat
	Messages
	WebSocket
}

func NewService(repo *repository.Repository, hub NotificationHub) *Service {
	authSrv := NewAuthService(repo.AuthRepository)
	chatSrv := NewChatService(repo.ChatRepository, repo.AuthRepository)
	msgSrv := NewMessagesService(repo.MessagesRepository, repo.ChatRepository, hub)
	wsSrv := NewWebsocketService(msgSrv, chatSrv, hub)

	return &Service{
		Authorization: authSrv,
		Chat:          chatSrv,
		Messages:      msgSrv,
		WebSocket:     wsSrv,
	}
}
