package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/user/model"
	"github.com/panagiotisptr/hermes-messenger/user/repository"
	"go.uber.org/zap"
)

// Service represents a user service
type Service struct {
	logger   *zap.Logger
	userRepo repository.Repository
}

// ProvideUserService provides an instance of the user
// service
func ProvideUserService(
	logger *zap.Logger,
	userRepo repository.Repository,
) *Service {
	return &Service{
		logger:   logger,
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (s *Service) CreateUser(
	ctx context.Context,
	args model.UserDetails,
) (*model.User, error) {
	return s.userRepo.Create(ctx, args)
}

// GetUser returns a user from their (uu)id
func (s *Service) GetUser(
	ctx context.Context,
	id uuid.UUID,
) (*model.User, error) {
	return s.userRepo.Get(ctx, id)
}

func (s *Service) GetUserByEmail(
	ctx context.Context,
	email string,
) (*model.User, error) {
	u, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, fmt.Errorf("Could not find user with email \"%s\"", email)
	}

	return u, nil
}
