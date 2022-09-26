package user

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	// Create creates a new user
	Create(context.Context, UserDetails) (*User, error)

	// Get gets a user by their uuid
	// returns nil if not found
	Get(context.Context, uuid.UUID) (*User, error)

	// GetByEmail finds a user from their email
	// returns nil if not found
	GetByEmail(context.Context, string) (*User, error)
}
