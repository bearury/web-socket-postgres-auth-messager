package entities

import "time"

type ChatType string

const (
	ChatTypeChat  ChatType = "chat"
	ChatTypeGroup ChatType = "group"
	ChatTypeChain ChatType = "chain"
)

type Chat struct {
	Id   string   `json:"id" db:"id"` // json:"-" - при сериализации в JSON поле игнорируется
	Name string   `json:"name" db:"name" binding:"required"`
	Type ChatType `json:"type" db:"type" binding:"required"`
}

type ChatUsers struct {
	ChatId    string    `json:"chatId" db:"chat_id" binding:"required"`
	UserId    string    `json:"userId" db:"user_id" binding:"required"`
	DeletedAt time.Time `json:"deletedAt" db:"deleted_at" binding:"required"`
}
