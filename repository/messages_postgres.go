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

//func (repo *MessagesPostgres) GetMessages(chatId string) ([]entities.MessageResponse, error) {
//	messages := make([]entities.MessageResponse, 0)
//
//	query := fmt.Sprintf("SELECT m.id, m.chat_id, m.user_id, m.content, m.created_at, u.username AS user_name FROM messages m LEFT JOIN users u ON m.user_id = u.id WHERE m.chat_id = $1 ORDER BY m.created_at")
//
//	err := repo.db.Select(&messages, query, chatId)
//
//	return messages, err
//}

func (repo *MessagesPostgres) GetMessages(chatId string) ([]entities.MessageResponseWithStatus, error) {
	messages := make([]entities.MessageResponseWithStatus, 0)

	var limit = 100
	var offset = 0

	query := fmt.Sprintf(`
		SELECT 
			m.id,
			m.chat_id,
			m.user_id,
			m.content,
			m.created_at,
			
			-- Статус ДОСТАВКИ: сообщение доставлено, если у всех ОСТАЛЬНЫХ участников 
			-- маркер last_delivered_at >= даты создания сообщения.
			-- (Для логики "доставлено хотя бы одному" замените MIN на MAX)
		   	COALESCE(
				MIN(CASE WHEN cu.user_id <> m.user_id THEN (cu.last_delivered_at >= m.created_at)::int END)::boolean, 
				TRUE
			) AS is_delivered,
		
			-- Статус ПРОЧТЕНИЯ: сообщение прочитано, если у всех ОСТАЛЬНЫХ участников 
			-- маркер last_read_at >= даты создания сообщения.
			-- (Для логики "прочитано хотя бы одному" замените MIN на MAX)
			COALESCE(
				MIN(CASE WHEN cu.user_id <> m.user_id THEN (cu.last_read_at >= m.created_at)::int END)::boolean, 
				FALSE
			) AS is_read
		
		FROM messages m
		-- Присоединяем участников этого же чата для проверки их маркеров
		LEFT JOIN chat_users cu ON cu.chat_id = m.chat_id AND cu.deleted_at IS NULL
		WHERE m.chat_id = $1
		GROUP BY m.id, m.chat_id, m.user_id, m.content, m.created_at
		ORDER BY m.created_at DESC
		LIMIT $2 OFFSET $3;
	`)

	err := repo.db.Select(&messages, query, chatId, limit, offset)

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
