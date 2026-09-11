package repository

import (
	"web-socket-postgres-auth-messager/entities"

	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	CreateUser(user entities.User) (string, error)
	GetUser(username string, password string) (entities.User, error)
	GetUserById(userId string) (entities.UserResponse, error)
}

type ChatRepository interface {
	GetChats(userId string) ([]entities.Chat, error)
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
