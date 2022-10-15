package memory

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/entityutils"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/entityutils/filter"
	"github.com/panagiotisptr/hermes-messenger/user/model"
	"github.com/panagiotisptr/hermes-messenger/user/repository"
	"go.uber.org/zap"
)

type MemoryRepository struct {
	users []*model.User
	entityutils.RepoHelper
	logger *zap.Logger
}

func ProvideUserRepository(
	logger *zap.Logger,
) repository.UserRepository {
	return &MemoryRepository{
		users:      make([]*model.User, 0),
		RepoHelper: entityutils.RepoHelper{},
		logger:     logger,
	}
}

func (r *MemoryRepository) Find(
	ctx context.Context,
	f filter.Filter,
) ([]*model.User, error) {
	users := []*model.User{}
	for _, u := range r.users {
		b, err := json.Marshal(u)
		if err != nil {
			return users, err
		}
		m := map[string]interface{}{}
		err = json.Unmarshal(b, &m)
		if err != nil {
			return nil, err
		}
		if f.Match(m) {
			users = append(users, u)
		}
	}

	return users, nil
}

func (r *MemoryRepository) FindOne(
	ctx context.Context,
	f filter.Filter,
) (*model.User, error) {
	for _, u := range r.users {
		b, err := json.Marshal(u)
		if err != nil {
			return nil, err
		}
		m := map[string]interface{}{}
		err = json.Unmarshal(b, &m)
		if err != nil {
			return nil, err
		}
		if f.Match(m) {
			return u, nil
		}
	}

	return nil, nil
}

func (r *MemoryRepository) Create(
	ctx context.Context,
	args *model.User,
) (*model.User, error) {
	if args.Email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	f := filter.NewFilter()
	f.Add("email", filter.Eq, args.Email)
	u, err := r.FindOne(ctx, f)
	if err != nil {
		return nil, err
	}
	if u != nil {
		return nil, fmt.Errorf("email already in use")
	}
	newUser := args
	id := uuid.New()
	newUser.ID = &id
	newUser.Meta = entityutils.Meta{}
	r.RepoHelper.UpdateMeta(
		ctx,
		&newUser.Meta,
		entityutils.CreateOp,
	)
	r.users = append(r.users, newUser)

	return newUser, nil
}

func (r *MemoryRepository) Update(
	ctx context.Context,
	f filter.Filter,
	args *model.User,
) ([]*model.User, error) {
	users := []*model.User{}
	users, err := r.Find(ctx, f)
	if err != nil {
		return users, nil
	}

	newUsers := []*model.User{}
	for i, u := range r.users {
		match := false
		for _, uu := range users {
			if uu.ID == u.ID {
				match = true
			}
		}
		if match {
			newU := *args
			newU.ID = u.ID
			newU.Meta = u.Meta
			r.RepoHelper.UpdateMeta(
				ctx,
				&newU.Meta,
				entityutils.UpdateOp,
			)
			r.users[i] = &newU
		}
	}
	r.users = newUsers

	return r.Find(ctx, f)
}

func (r *MemoryRepository) Delete(
	ctx context.Context,
	f filter.Filter,
) ([]*model.User, error) {
	users := []*model.User{}
	users, err := r.Find(ctx, f)
	if err != nil {
		return users, nil
	}

	newUsers := []*model.User{}
	for _, u := range r.users {
		match := false
		for _, uu := range users {
			if uu.ID == u.ID {
				match = true
			}
		}
		if !match {
			newUsers = append(newUsers, u)
		}
	}
	r.users = newUsers

	return users, nil
}
