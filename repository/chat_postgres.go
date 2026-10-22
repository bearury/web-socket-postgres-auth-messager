package repository

import (
	"context"
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

	query := fmt.Sprintf("SELECT c.* FROM chats c JOIN chat_users cu ON cu.chat_id = c.id WHERE cu.user_id = $1 AND cu.deleted_at IS NULL AND (c.creator = $1 OR EXISTS (SELECT 1 FROM messages m WHERE m.chat_id = c.id))")

	err := repo.db.Select(&chats, query, userId)

	return chats, err
}

func (repo *ChatPostgres) GetChatByUsersIds(user1Id, user2Id string) (entities.Chat, error) {
	chat := entities.Chat{}

	query := fmt.Sprintf("SELECT c.id, c.name, c.type FROM chats c JOIN chat_users cu1 ON cu1.chat_id = c.id JOIN chat_users cu2 ON cu2.chat_id = c.id WHERE c.type = 'chat' AND cu1.user_id = $1 AND cu2.user_id = $2")

	err := repo.db.Get(&chat, query, user1Id, user2Id)

	return chat, err
}

func (repo *ChatPostgres) GetUsersChatByIdWithoutUser(chatId string, userId string) ([]entities.ChatUsers, error) {
	chats := make([]entities.ChatUsers, 0)

	query := fmt.Sprintf("SELECT cu.chat_id, cu.user_id, cu.deleted_at FROM chat_users cu WHERE cu.chat_id = $1 AND cu.user_id != $2")

	err := repo.db.Select(&chats, query, chatId, userId)

	return chats, err
}

func (repo *ChatPostgres) AddChat(name string, chatType entities.ChatType, creator string) (entities.Chat, error) {
	chat := entities.Chat{}

	query := fmt.Sprintf("INSERT INTO chats (name, type, creator) VALUES ($1, $2, $3) RETURNING id, name, type")

	row := repo.db.QueryRow(query, name, chatType, creator)

	if err := row.Scan(&chat.Id, &chat.Name, &chat.Type); err != nil {
		return chat, err
	}

	return chat, nil
}

func (repo *ChatPostgres) AddChatUser(chatId, userId string) (string, error) {
	chat := entities.ChatUsers{}

	query := fmt.Sprintf("INSERT INTO chat_users (chat_id, user_id, deleted_at) VALUES ($1, $2, null) RETURNING chat_id, user_id")

	row := repo.db.QueryRow(query, chatId, userId)

	if err := row.Scan(&chat.ChatId, &chat.UserId); err != nil {
		return chat.ChatId, err
	}

	return chat.ChatId, nil
}

func (repo *ChatPostgres) UpdateChatUserDelevered(ctx context.Context, chatId, userId string) (entities.ChatUsers, error) {
	chatUsers := entities.ChatUsers{}

	query := fmt.Sprintf(`
		UPDATE chat_users 
		SET 
			last_delivered_at = now(),
			deleted_at = null
		WHERE chat_id = $1 AND user_id = $2
		RETURNING chat_id, user_id, last_delivered_at, last_read_at, deleted_at;
	`)

	err := repo.db.QueryRowContext(ctx, query, chatId, userId).Scan(
		&chatUsers.ChatId,
		&chatUsers.UserId,
		&chatUsers.LastDeleveredAt,
		&chatUsers.LastReadAt,
		&chatUsers.DeletedAt,
	)
	if err != nil {
		return entities.ChatUsers{}, err
	}

	return chatUsers, nil
}

func (repo *ChatPostgres) UpdateChatUserRead(ctx context.Context, chatId, userId string) (entities.ChatUsers, error) {
	chatUsers := entities.ChatUsers{}

	query := fmt.Sprintf(`
		UPDATE chat_users 
		SET 
			last_read_at = now(),
			deleted_at = null
		WHERE chat_id = $1 AND user_id = $2
		RETURNING chat_id, user_id, last_delivered_at, last_read_at, deleted_at;
	`)

	err := repo.db.QueryRowContext(ctx, query, chatId, userId).Scan(
		&chatUsers.ChatId,
		&chatUsers.UserId,
		&chatUsers.LastDeleveredAt,
		&chatUsers.LastReadAt,
		&chatUsers.DeletedAt,
	)
	if err != nil {
		return entities.ChatUsers{}, err
	}

	return chatUsers, nil
}
