package entities

type WSMessage struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
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

//id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
//chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
//user_id UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
//content TEXT NOT NULL,
//created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
