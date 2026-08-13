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
