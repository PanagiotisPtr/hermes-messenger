package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/entityutils/filter"
	"github.com/panagiotisptr/hermes-messenger/user/model"
	"github.com/panagiotisptr/hermes-messenger/user/repository"
	"go.uber.org/zap"
)

// UserService represents a user service
type UserService struct {
	logger   *zap.Logger
	userRepo repository.UserRepository
}

// ProvideUserService provides an instance of the user
// service
func ProvideUserService(
	logger *zap.Logger,
	userRepo repository.UserRepository,
) *UserService {
	return &UserService{
		logger: logger.With(
			zap.String("service", "UserService"),
		),
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(
	ctx context.Context,
	args *model.User,
) (*model.User, error) {
	return s.userRepo.Create(ctx, args)
}

// GetUser returns a user from their (uu)id
func (s *UserService) GetUser(
	ctx context.Context,
	id uuid.UUID,
) (*model.User, error) {
	f := filter.NewFilter()
	f.Add("id", filter.Eq, id)

	return s.userRepo.FindOne(ctx, f)
}

func (s *UserService) GetUserByEmail(
	ctx context.Context,
	email string,
) (*model.User, error) {
	f := filter.NewFilter()
	f.Add("email", filter.Eq, email)

	u, err := s.userRepo.FindOne(ctx, f)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, fmt.Errorf("Could not find user with email \"%s\"", email)
	}

	return u, nil
}
