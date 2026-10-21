package repository

import (
	"fmt"
	"web-socket-postgres-auth-messager/entities"

	"github.com/jmoiron/sqlx"
)

type MessagesPostgres struct {
	db *sqlx.DB
}

func NewMessagesPostgres(db *sqlx.DB) *MessagesPostgres {
	return &MessagesPostgres{db: db}
}

func (repo *MessagesPostgres) GetMessages(chatId string) ([]entities.MessageResponse, error) {
	messages := make([]entities.MessageResponse, 0)

	query := fmt.Sprintf("SELECT m.id, m.chat_id, m.user_id, m.content, m.created_at, u.username AS user_name FROM messages m LEFT JOIN users u ON m.user_id = u.id WHERE m.chat_id = $1 ORDER BY m.created_at")

	err := repo.db.Select(&messages, query, chatId)

	return messages, err
}

func (repo *MessagesPostgres) AddMessage(userId, chatId, content string) (entities.MessageResponse, error) {
	message := entities.MessageResponse{}

	query := fmt.Sprintf("INSERT INTO messages (user_id, chat_id, content) VALUES ($1, $2, $3) RETURNING id, user_id, chat_id, content, created_at")

	row := repo.db.QueryRow(query, userId, chatId, content)

	if err := row.Scan(&message.Id, &message.UserId, &message.ChatId, &message.Content, &message.CreatedAt); err != nil {
		return message, err
	}

	return message, nil
}
