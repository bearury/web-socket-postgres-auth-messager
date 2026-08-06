package service

import (
	"web-socket-postgres-auth-messager/entities"
	"web-socket-postgres-auth-messager/repository"
)

type AuthService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user entities.User) (string, error) {
	return "Test", nil
}
func (s *AuthService) GenerateToken(username, password string) (string, error) {
	return "Test", nil
}
func (s *AuthService) ParseToken(token string) (string, error) {
	return "Test", nil
}
