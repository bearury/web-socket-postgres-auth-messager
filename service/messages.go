package service

import (
	"context"
	"web-socket-postgres-auth-messager/entities"
	"web-socket-postgres-auth-messager/repository"
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
	return s.repo.AddMessage(user.Id, chatId, content)
}
