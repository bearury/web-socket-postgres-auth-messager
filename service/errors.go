package service

import "errors"

var (
	ErrChatAlreadyExists = errors.New("chat already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserIDsEmpty      = errors.New("userIds is empty")
	ErrInvalidChatType   = errors.New("invalid chat type")
	ErrInvalidUserID     = errors.New("invalid user id")
)
