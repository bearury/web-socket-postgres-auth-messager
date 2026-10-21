package service

import (
	"web-socket-postgres-auth-messager/entities"
	"web-socket-postgres-auth-messager/repository"
)

type MessagesService struct {
	repo repository.MessagesRepository
}

func NewMessagesService(repo repository.MessagesRepository) *MessagesService {
	return &MessagesService{repo: repo}
}

func (s *MessagesService) GetMessages(chatId string) ([]entities.MessageResponse, error) {
	return s.repo.GetMessages(chatId)
}

func (s *MessagesService) AddMessage(userId, chatId, content string) (entities.MessageResponse, error) {
	return s.repo.AddMessage(userId, chatId, content)
}
