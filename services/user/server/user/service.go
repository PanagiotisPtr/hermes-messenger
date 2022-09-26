package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Service represents a user service
type Service struct {
	logger   *zap.Logger
	userRepo Repository
}

// ProvideUserService provides an instance of the user
// service
func ProvideUserService(
	logger *zap.Logger,
	userRepo Repository,
) *Service {
	return &Service{
		logger:   logger,
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (s *Service) CreateUser(
	ctx context.Context,
	args UserDetails,
) (*User, error) {
	return s.userRepo.Create(ctx, args)
}

// GetUser returns a user from their (uu)id
func (s *Service) GetUser(
	ctx context.Context,
	id uuid.UUID,
) (*User, error) {
	return s.userRepo.Get(ctx, id)
}

func (s *Service) GetUserByEmail(
	ctx context.Context,
	email string,
) (*User, error) {
	u, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, fmt.Errorf("Could not find user with email \"%s\"", email)
	}

	return u, nil
}
