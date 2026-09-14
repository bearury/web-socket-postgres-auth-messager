package service

import (
	"database/sql"
	"errors"
	"fmt"
	"web-socket-postgres-auth-messager/entities"
	"web-socket-postgres-auth-messager/repository"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/google/uuid"
)

type ChatService struct {
	repo     repository.ChatRepository
	userRepo repository.AuthRepository
}

func NewChatService(repo repository.ChatRepository, userRepo repository.AuthRepository) *ChatService {
	return &ChatService{repo: repo, userRepo: userRepo}
}

func (s *ChatService) GetChats(userId string) ([]entities.Chat, error) {
	return s.repo.GetChats(userId)
}

func (s *ChatService) AddChat(name string, chatType entities.ChatType, userId string, userIds []string) (entities.Chat, error) {

	if len(userIds) == 0 {
		return entities.Chat{}, ErrUserIDsEmpty
	}

	switch chatType {
	case entities.ChatTypeChat:

		user2Id, err := uuid.Parse(userIds[0])
		if err != nil {
			return entities.Chat{}, fmt.Errorf(
				"%w: %s",
				ErrInvalidUserID,
				userIds[0],
			)
		}

		user2, err := s.userRepo.GetUserById(user2Id.String())

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return entities.Chat{}, ErrUserNotFound
			}
			return entities.Chat{}, fmt.Errorf("Error get user2: %w", err)
		}

		// проверка что чат создан
		chat, err := s.repo.GetChatByUsersIds(userId, user2.Id)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// Чат НЕ найден
				// создать новый чат
				chat, err := s.repo.AddChat(name, chatType)

				if err != nil {
					return entities.Chat{}, fmt.Errorf("failed to add chat: %w", err)
				}
				// добавить текущего юзера к чату
				_, err = s.repo.AddChatUser(chat.Id, userId)

				if err != nil {
					return entities.Chat{}, fmt.Errorf("failed to add chat user1: %w", err)
				}
				// добавить юзера 2 к чату
				_, err = s.repo.AddChatUser(chat.Id, user2.Id)

				if err != nil {
					return entities.Chat{}, fmt.Errorf("failed to add chat user2: %w", err)
				}

				return chat, nil
			} else {
				return entities.Chat{}, fmt.Errorf("error getting chat by userIds: %w", err)
			}
		}
		return chat, nil
	case entities.ChatTypeGroup:
		logger.Info("Croup")
		return s.repo.AddChat(name, chatType)
	case entities.ChatTypeChain:
		logger.Info("Chain")
		return s.repo.AddChat(name, chatType)
	default:
		return entities.Chat{}, fmt.Errorf("unknown chat type: %s", chatType)
	}
}
