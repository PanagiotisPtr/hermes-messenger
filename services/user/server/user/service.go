package user

import (
	"context"

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
