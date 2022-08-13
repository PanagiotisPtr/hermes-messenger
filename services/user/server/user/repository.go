package user

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	AddUser(context.Context, string) (*User, error)
	GetUser(context.Context, uuid.UUID) (*User, error)
}
