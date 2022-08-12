package user

import "github.com/google/uuid"

type Repository interface {
	AddUser(string) string
	GetUser(uuid.UUID) (*User, error)
}
