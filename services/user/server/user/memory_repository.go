package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MemoryRepository struct {
	users  []*User
	logger *zap.Logger
}

func ProvideMemoryRepository(
	logger *zap.Logger,
) Repository {
	return &MemoryRepository{
		users:  make([]*User, 0),
		logger: logger,
	}
}

func (r *MemoryRepository) Create(
	ctx context.Context,
	args UserDetails,
) (*User, error) {
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
	r.users = append(r.users, &User{
		ID:          uuid.New(),
		UserDetails: args,
	})

	return r.users[len(r.users)-1], nil
}

func (r *MemoryRepository) Get(
	ctx context.Context,
	id uuid.UUID,
) (*User, error) {
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
) (*User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}

	return nil, nil
}
