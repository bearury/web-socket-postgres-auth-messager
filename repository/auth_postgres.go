package repository

import (
	"fmt"
	"web-socket-postgres-auth-messager/entities"

	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (repo *AuthPostgres) GetUser(username, password string) (entities.User, error) {
	var user entities.User

	query := fmt.Sprintf("SELECT id FROM users WHERE username=$1 AND password_hash = $2")

	err := repo.db.Get(&user, query, username, password)

	return user, err
}

func (repo *AuthPostgres) CreateUser(user entities.User) (string, error) {
	var id string

	query := fmt.Sprintf("INSERT INTO users (name, username, password_hash) VALUES ($1, $2, $3) RETURNING id")

	row := repo.db.QueryRow(query, user.Name, user.Username, user.Password)

	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (repo *AuthPostgres) GetUserById(userId string) (entities.UserResponse, error) {
	var user entities.UserResponse

	query := fmt.Sprintf("SELECT id, name, username FROM users WHERE id = $1")

	err := repo.db.Get(&user, query, userId)

	return user, err
}

func (repo *AuthPostgres) GetAllUsersWithoutUser(userId string) ([]entities.UserResponse, error) {
	users := make([]entities.UserResponse, 0)

	query := fmt.Sprintf("SELECT id, name, username FROM users WHERE id IS DISTINCT FROM $1")

	err := repo.db.Select(&users, query, userId)

	return users, err
}
