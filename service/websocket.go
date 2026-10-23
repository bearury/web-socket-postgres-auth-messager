package service

import (
	"encoding/json"
	"web-socket-postgres-auth-messager/entities"
	ws_hub "web-socket-postgres-auth-messager/handler/ws-hub"

	"github.com/sirupsen/logrus"
)

type WebSocketService struct {
	msgSrv  *MessagesService
	chatSrv *ChatService
	hub     NotificationHub
}

func NewWebsocketService(msgSrv *MessagesService, chatSrv *ChatService, hub NotificationHub) *WebSocketService {
	return &WebSocketService{msgSrv: msgSrv, chatSrv: chatSrv, hub: hub}
}

func (s *WebSocketService) HandleMessage(client *ws_hub.Client, msg *entities.WSMessage) error {
	var payload entities.WSPayloadMessageResponse

	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		logrus.Errorf("Failed to unmarshal message payload: %s", err)
		return err
	}

	content := payload.Content
	chatId := payload.To.Id

	message, err := s.msgSrv.AddMessage(client.User, chatId, content)

	if err != nil {
		logrus.Error("Failed to create message in database: %s", err)
		return err
	}

	s.SendToChatMessage(message.ChatId, client.User, message.Content)
	return nil
}

func (s *WebSocketService) SendToChatMessage(chatId string, user entities.UserResponse, content string) {

	users, err := s.chatSrv.GetUsersChatByIdWithoutUser(chatId, user.Id)

	if err != nil {
		logrus.Error("Error Get Chat_users By Id Without current User. Error:", err)
	}

	userIDs := make([]string, len(users))

	for i, user := range users {
		userIDs[i] = user.UserId
	}

	// 1. Создаем и наполняем payload конкретными данными
	payloadData := entities.WSPayloadMessageResponse{
		From:    user.Username,
		Content: content,
	}
	payloadData.To.Id = chatId
	payloadData.To.Type = entities.ChatTypeChat

	// 2. Маршалим payload в сырые байты
	payloadBytes, err := json.Marshal(payloadData)
	if err != nil {
		logrus.Errorf("Failed to marshal payload: %s", err)
		return
	}

	// 3. Собираем финальное сообщение для отправки
	msg := entities.WSMessage{
		Type:    entities.WSTypeMessage,        // или просто "message"
		Payload: json.RawMessage(payloadBytes), // Явное приведение []byte -> json.RawMessage
	}

	s.hub.SendToChat(userIDs, msg)
}
