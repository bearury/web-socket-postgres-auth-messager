package entities

type WSMessage struct {
	Type    string                   `json:"type"`
	Payload WSPayloadMessageResponse `json:"payload"`
}

type WSPayloadToMessage struct {
	Type ChatType `json:"type"`
	Id   string   `json:"id"`
}

type WSPayloadMessageResponse struct {
	From    string             `json:"from"`
	To      WSPayloadToMessage `json:"to"`
	Content string             `json:"content"`
}

type WSMessageResponse struct {
	Type    string                   `json:"type"`
	Payload WSPayloadMessageResponse `json:"payload"`
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
