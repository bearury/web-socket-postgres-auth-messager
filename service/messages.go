package service

import (
	"context"
	"web-socket-postgres-auth-messager/entities"
	"web-socket-postgres-auth-messager/repository"

	"github.com/sirupsen/logrus"
)

type MessagesService struct {
	repo      repository.MessagesRepository
	chatsRepo repository.ChatRepository
	hub       NotificationHub
}

func NewMessagesService(repo repository.MessagesRepository, chatsRepo repository.ChatRepository, hub NotificationHub) *MessagesService {
	return &MessagesService{repo: repo, chatsRepo: chatsRepo, hub: hub}
}

func (s *MessagesService) GetMessages(ctx context.Context, chatId, userId string) ([]entities.MessageResponseWithStatus, error) {
	messages, err := s.repo.GetMessages(chatId)
	if err != nil {
		return nil, err
	}

	_, err = s.chatsRepo.UpdateChatUserRead(ctx, chatId, userId)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *MessagesService) AddMessage(user entities.UserResponse, chatId, content string) (entities.MessageResponse, error) {

	message, err := s.repo.AddMessage(user.Id, chatId, content)
	if err != nil {
		return entities.MessageResponse{}, err
	}

	s.sendToChatMessage(message.ChatId, user, message.Content)

	return message, nil
}

func (s *MessagesService) sendToChatMessage(chatId string, user entities.UserResponse, content string) {

	users, err := s.chatsRepo.GetUsersChatByIdWithoutUser(chatId, user.Id)

	if err != nil {
		logrus.Error("Error Get Chat_users By Id Without current User. Error:", err)
	}

	userIDs := make([]string, len(users))

	for i, user := range users {
		userIDs[i] = user.UserId
	}

	s.hub.SendToChat(userIDs,
		entities.WSMessageResponse{
			Type: "message",
			Payload: entities.WSPayloadMessageResponse{
				From: user.Username,
				To: entities.WSPayloadToMessage{
					Type: entities.ChatTypeChat,
					Id:   chatId,
				},
				Content: content,
			},
		})
}
