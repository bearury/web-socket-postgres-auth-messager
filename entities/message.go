package entities

type Message struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

type PayloadMessageResponse struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}

type MessageResponse struct {
	Type    string                 `json:"type"`
	Payload PayloadMessageResponse `json:"payload"`
}

//id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
//chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
//user_id UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
//content TEXT NOT NULL,
//created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
