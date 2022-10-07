package memory

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/user/model"
	"github.com/panagiotisptr/hermes-messenger/user/repository"
	"go.uber.org/zap"
)

type MemoryRepository struct {
	users  []*model.User
	logger *zap.Logger
}

func ProvideUserRepository(
	logger *zap.Logger,
) repository.Repository {
	return &MemoryRepository{
		users:  make([]*model.User, 0),
		logger: logger,
	}
}

func (r *MemoryRepository) Create(
	ctx context.Context,
	args model.UserDetails,
) (*model.User, error) {
	if args.Email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	u, err := r.GetByEmail(
		ctx,
		args.Email,
	)
	if err != nil {
		return nil, err
	}
	if u != nil {
		return nil, fmt.Errorf("email already in use")
	}
	r.users = append(r.users, &model.User{
		ID:          uuid.New(),
		UserDetails: args,
	})

	return r.users[len(r.users)-1], nil
}

func (r *MemoryRepository) Get(
	ctx context.Context,
	id uuid.UUID,
) (*model.User, error) {
	for _, u := range r.users {
		if u.ID == id {
			return u, nil
		}
	}

	return nil, nil
}

func (r *MemoryRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*model.User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}

	return nil, nil
}
