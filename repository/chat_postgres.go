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

	query := fmt.Sprintf("SELECT c.* FROM chats c JOIN chat_users cu ON cu.chat_id=c.id WHERE cu.user_id=$1 AND cu.deleted_at IS NULL AND EXISTS ( SELECT 1 FROM messages m WHERE m.chat_id = c.id)")

	err := repo.db.Select(&chats, query, userId)

	return chats, err
}

func (repo *ChatPostgres) GetChatByUsersIds(user1Id, user2Id string) (entities.Chat, error) {
	chat := entities.Chat{}

	query := fmt.Sprintf("SELECT c.id, c.name, c.type FROM chats c JOIN chat_users cu1 ON cu1.chat_id = c.id JOIN chat_users cu2 ON cu2.chat_id = c.id WHERE c.type = 'chat' AND cu1.user_id = $1 AND cu2.user_id = $2")

	err := repo.db.Get(&chat, query, user1Id, user2Id)

	return chat, err
}

func (repo *ChatPostgres) AddChat(name string, chatType entities.ChatType) (entities.Chat, error) {
	chat := entities.Chat{}

	query := fmt.Sprintf("INSERT INTO %s (name, type) VALUES ($1, $2) RETURNING id, name, type", chatsTable)

	row := repo.db.QueryRow(query, name, chatType)

	if err := row.Scan(&chat.Id, &chat.Name, &chat.Type); err != nil {
		return chat, err
	}

	return chat, nil
}

func (repo *ChatPostgres) AddChatUser(chatId, userId string) (string, error) {
	chat := entities.ChatUsers{}

	query := fmt.Sprintf("INSERT INTO %s (chat_id, user_id, deleted_at) VALUES ($1, $2, null) RETURNING chat_id, user_id", chatUsersTable)

	row := repo.db.QueryRow(query, chatId, userId)

	if err := row.Scan(&chat.ChatId, &chat.UserId); err != nil {
		return chat.ChatId, err
	}

	return chat.ChatId, nil
}
