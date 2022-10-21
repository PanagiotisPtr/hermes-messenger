package memory

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/entityutils"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/entityutils/filter"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/loggingutils"
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
) repository.UserRepository {
	return &MemoryRepository{
		users: make([]*model.User, 0),
		logger: logger.With(
			zap.String("repository", "UserRepository"),
			zap.String("type", "memory"),
		),
	}
}

func (r *MemoryRepository) Find(
	ctx context.Context,
	f filter.Filter,
) (<-chan *model.User, error) {
	ch := make(chan *model.User)
	l := loggingutils.
		LoggerWithRequestID(ctx, r.logger).
		With(
			zap.String("method", "Find"),
		).Sugar()
	go func() {
		defer close(ch)
		for _, u := range r.users {
			b, err := json.Marshal(u)
			if err != nil {
				l.Error(
					"marshalling user",
					err,
				)
			}
			m := map[string]interface{}{}
			err = json.Unmarshal(b, &m)
			if err != nil {
				l.Error(
					"unmarshalling user",
					err,
				)
			}
			if f.Match(m) {
				ch <- u
			}
		}
	}()

	return ch, nil
}

func (r *MemoryRepository) FindOne(
	ctx context.Context,
	f filter.Filter,
) (*model.User, error) {
	l := loggingutils.
		LoggerWithRequestID(ctx, r.logger).
		With(
			zap.String("method", "FindOne"),
		).Sugar()
	ch, err := r.Find(ctx, f)
	if err != nil {
		l.Error("failed to find user:", err)
		return nil, err
	}
	for u := range ch {
		return u, nil
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
	newUser.UpdateMeta(
		ctx,
		entityutils.CreateOp,
	)
	r.users = append(r.users, newUser)

	return newUser, nil
}

func (r *MemoryRepository) Update(
	ctx context.Context,
	f filter.Filter,
	args *model.User,
) (int64, error) {
	l := loggingutils.
		LoggerWithRequestID(ctx, r.logger).
		With(
			zap.String("method", "Update"),
		).Sugar()
	updated := int64(0)
	users, err := r.Find(ctx, f)
	if err != nil {
		return updated, err
	}
	ids := []uuid.UUID{}
	for u := range users {
		if u.ID == nil {
			l.Error("found user with nil UUID", "filter:", f)
		}
		ids = append(ids, *u.ID)
	}

	newUsers := []*model.User{}
	for i, u := range r.users {
		match := false
		for _, id := range ids {
			if u.ID != nil && *u.ID == id {
				match = true
			}
		}
		if match {
			newU := *args
			newU.ID = u.ID
			newU.Meta = u.Meta
			newU.UpdateMeta(
				ctx,
				entityutils.UpdateOp,
			)
			r.users[i] = &newU
			updated++
		}
	}
	r.users = newUsers

	return updated, nil
}

func (r *MemoryRepository) Delete(
	ctx context.Context,
	f filter.Filter,
) (int64, error) {
	l := loggingutils.
		LoggerWithRequestID(ctx, r.logger).
		With(
			zap.String("method", "Delete"),
		).Sugar()
	deleted := int64(0)
	users, err := r.Find(ctx, f)
	if err != nil {
		return deleted, err
	}
	ids := []uuid.UUID{}
	for u := range users {
		if u.ID == nil {
			l.Error("found user with nil UUID", "filter:", f)
		}
		ids = append(ids, *u.ID)
	}

	newUsers := []*model.User{}
	for _, u := range r.users {
		match := false
		for _, id := range ids {
			if u.ID != nil && *u.ID == id {
				match = true
			}
		}
		if !match {
			newUsers = append(newUsers, u)
		} else {
			deleted++
		}
	}
	r.users = newUsers

	return deleted, nil
}
