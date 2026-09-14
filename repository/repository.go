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
	AddChat(name string, chatType entities.ChatType) (entities.Chat, error)
	GetChatByUsersIds(user1Id, user2Id string) (entities.Chat, error)
	AddChatUser(chatId string, userId string) (string, error)
}

type Repository struct {
	AuthRepository
	ChatRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		AuthRepository: NewAuthPostgres(db),
		ChatRepository: NewChatPostgres(db),
	}
}
