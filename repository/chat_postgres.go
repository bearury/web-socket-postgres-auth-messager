package repository

import (
	"fmt"
	"web-socket-postgres-auth-messager/entities"

	"github.com/jmoiron/sqlx"
)

type ChatPostgres struct {
	db *sqlx.DB
}

func NewChatPostgres(db *sqlx.DB) *ChatPostgres {
	return &ChatPostgres{db: db}
}

func (repo *ChatPostgres) GetChats(userId string) ([]entities.Chat, error) {
	chats := make([]entities.Chat, 0)

	query := fmt.Sprintf("SELECT c.* FROM chats c JOIN chat_users cu ON cu.chat_id=c.id WHERE cu.user_id=$1")

	err := repo.db.Select(&chats, query, userId)

	return chats, err
}
