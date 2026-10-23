package entities

import "encoding/json"

type WSType string

const (
	WSTypeMessage WSType = "message"
	WSTypeNotify  WSType = "notify"
	WSTypePing    WSType = "ping"
)

type WSNotifyType string

const (
	WSNotifyDelivery WSNotifyType = "notify_delivery"
	WSNotifyRead     WSNotifyType = "notify_read"
)

type WSMessage struct {
	Type    WSType          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type WSPayloadMessageResponse struct {
	From    string `json:"from"`
	Content string `json:"content"`
	To      struct {
		Id   string   `json:"id"`
		Type ChatType `json:"type"`
	} `json:"to"`
}

type WSPayloadNotifyResponse struct {
	From string       `json:"from"`
	Type WSNotifyType `json:"type"`
	To   struct {
		Id   string   `json:"id"`
		Type ChatType `json:"type"`
	} `json:"to"`
}

type MessageResponse struct {
	Id        string `json:"id" db:"id" binding:"required"`
	ChatId    string `json:"chatId" db:"chat_id" binding:"required"`
	UserId    string `json:"userId" db:"user_id" binding:"required"`
	UserName  string `json:"userName" db:"user_name"`
	Content   string `json:"content" db:"content" binding:"required"`
	CreatedAt string `json:"createdAt" db:"created_at" binding:"required"`
}

type MessageResponseWithStatus struct {
	Id          string `json:"id" db:"id" binding:"required"`
	ChatId      string `json:"chatId" db:"chat_id" binding:"required"`
	UserId      string `json:"userId" db:"user_id" binding:"required"`
	UserName    string `json:"userName" db:"user_name"`
	Content     string `json:"content" db:"content" binding:"required"`
	CreatedAt   string `json:"createdAt" db:"created_at" binding:"required"`
	IsDelivered bool   `json:"isDelivered" db:"is_delivered"`
	IsRead      bool   `json:"isRead" db:"is_read"`
}
