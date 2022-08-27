package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	logger   *zap.Logger
	userRepo Repository
}

func ProvideUserService(
	logger *zap.Logger,
	userRepo Repository,
) *Service {
	return &Service{
		logger:   logger,
		userRepo: userRepo,
	}
}

func (s *Service) RegisterUser(
	ctx context.Context,
	email string,
) (*User, error) {
	return s.userRepo.AddUser(ctx, email)
}

func (s *Service) GetUser(
	ctx context.Context,
	id uuid.UUID,
) (*User, error) {
	return s.userRepo.GetUser(ctx, id)
}

func (s *Service) GetUserByEmail(
	ctx context.Context,
	email string,
) (*User, error) {
	u, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, fmt.Errorf("Could not find user with email \"%s\"", email)
	}

	return u, nil
}
