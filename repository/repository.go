package repository

import (
	"web-socket-postgres-auth-messager/entities"

	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	CreateUser(user entities.User) (string, error)
	GetUser(username string, password string) (entities.User, error)
	GetUserById(userId string) (entities.UserResponse, error)
	GetAllUsersWithoutUser(userId string) ([]entities.UserResponse, error)
}

type ChatRepository interface {
	GetChats(userId string) ([]entities.Chat, error)
	AddChat(name string, chatType entities.ChatType, creator string) (entities.Chat, error)
	GetChatByUsersIds(user1Id, user2Id string) (entities.Chat, error)
	AddChatUser(chatId string, userId string) (string, error)
	GetUsersChatByIdWithoutUser(chatId string, userId string) ([]entities.ChatUsers, error)
}

type MessagesRepository interface {
	GetMessages(chatId string) ([]entities.MessageResponse, error)
	AddMessage(userId, chatId, content string) (entities.MessageResponse, error)
}

type Repository struct {
	AuthRepository
	ChatRepository
	MessagesRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		AuthRepository:     NewAuthPostgres(db),
		ChatRepository:     NewChatPostgres(db),
		MessagesRepository: NewMessagesPostgres(db),
	}
}
