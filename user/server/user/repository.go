package user

import "github.com/google/uuid"

type Repository interface {
	AddUser(uuid.UUID, string) error
	GetUser(uuid.UUID) (*User, error)
}
