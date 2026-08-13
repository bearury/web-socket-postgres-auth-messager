package service

import (
	"crypto/sha1"
	"fmt"
	"time"
	"web-socket-postgres-auth-messager/entities"
	"web-socket-postgres-auth-messager/repository"

	"github.com/dgrijalva/jwt-go"
)

const (
	solt       = "dsffsfdsgv45ttrdsfg"
	signingKey = "secret"
	tokenTTL   = time.Hour * 12
)

type AuthService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
}

func (s *AuthService) CreateUser(user entities.User) (string, error) {
	user.Password = s.generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, s.generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey))
}
func (s *AuthService) ParseToken(token string) (string, error) {
	accessToken, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := accessToken.Claims.(*tokenClaims)
	if !ok {
		return "", fmt.Errorf("Invalid token claims")
	}

	return claims.UserId, nil
}

func (service *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(solt)))
}

func (s *AuthService) GetUserById(userId string) (entities.UserResponse, error) {
	user, err := s.repo.GetUserById(userId)
	return user, err
}
