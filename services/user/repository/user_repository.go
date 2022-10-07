package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/user/model"
)

type Repository interface {
	// Create creates a new user
	Create(context.Context, model.UserDetails) (*model.User, error)

	// Get gets a user by their uuid
	// returns nil if not found
	Get(context.Context, uuid.UUID) (*model.User, error)

	// GetByEmail finds a user from their email
	// returns nil if not found
	GetByEmail(context.Context, string) (*model.User, error)
}
