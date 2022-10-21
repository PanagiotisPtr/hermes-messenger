package memory

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/friends/model"
	"github.com/panagiotisptr/hermes-messenger/friends/repository"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/entityutils"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/entityutils/filter"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/loggingutils"
	"go.uber.org/zap"
)

type MemoryRepository struct {
	friends []*model.Friend
	logger  *zap.Logger
}

func ProvideFriendRepository(
	logger *zap.Logger,
) repository.FriendRepository {
	return &MemoryRepository{
		friends: make([]*model.Friend, 0),
		logger: logger.With(
			zap.String("repository", "FriendRepository"),
			zap.String("type", "memory"),
		),
	}
}

func (r *MemoryRepository) Find(
	ctx context.Context,
	f filter.Filter,
) (<-chan *model.Friend, error) {
	ch := make(chan *model.Friend)
	l := loggingutils.
		LoggerWithRequestID(ctx, r.logger).
		With(
			zap.String("method", "Find"),
		).Sugar()
	go func() {
		defer close(ch)
		for _, friend := range r.friends {
			b, err := json.Marshal(friend)
			if err != nil {
				l.Error(
					"marshalling friend",
					err,
				)
			}
			m := map[string]interface{}{}
			err = json.Unmarshal(b, &m)
			if err != nil {
				l.Error(
					"unmarshalling friend",
					err,
				)
			}
			if f.Match(m) {
				ch <- friend
			}
		}
	}()

	return ch, nil
}

func (r *MemoryRepository) FindOne(
	ctx context.Context,
	f filter.Filter,
) (*model.Friend, error) {
	l := loggingutils.
		LoggerWithRequestID(ctx, r.logger).
		With(
			zap.String("method", "FindOne"),
		).Sugar()
	ch, err := r.Find(ctx, f)
	if err != nil {
		l.Error("failed to find friend", err)
		return nil, err
	}
	for friend := range ch {
		return friend, nil
	}

	return nil, nil
}

func (r *MemoryRepository) Create(
	ctx context.Context,
	args *model.Friend,
) (*model.Friend, error) {
	f := filter.NewFilter()
	f.Add("friendId", filter.Eq, args.FriendID)
	friend, err := r.FindOne(ctx, f)
	if err != nil {
		return nil, err
	}
	if friend != nil {
		return nil, fmt.Errorf("connection between users already exists")
	}
	newFriend := args
	id := uuid.New()
	newFriend.ID = &id
	newFriend.Meta = entityutils.Meta{}
	newFriend.UpdateMeta(
		ctx,
		entityutils.CreateOp,
	)
	r.friends = append(r.friends, newFriend)

	return newFriend, nil
}

func (r *MemoryRepository) Update(
	ctx context.Context,
	f filter.Filter,
	args *model.Friend,
) (int64, error) {
	l := loggingutils.
		LoggerWithRequestID(ctx, r.logger).
		With(
			zap.String("method", "Update"),
		).Sugar()
	updated := int64(0)
	friends, err := r.Find(ctx, f)
	if err != nil {
		return updated, err
	}
	ids := []uuid.UUID{}
	for friend := range friends {
		if friend.ID == nil {
			l.Error("found friend with nil UUID", "filter:", f)
		}
		ids = append(ids, *friend.ID)
	}

	newFriends := []*model.Friend{}
	for i, friend := range r.friends {
		match := false
		for _, id := range ids {
			if friend.ID != nil && *friend.ID == id {
				match = true
			}
		}
		if match {
			newF := *args
			newF.ID = friend.ID
			newF.Meta = friend.Meta
			newF.UpdateMeta(
				ctx,
				entityutils.UpdateOp,
			)
			r.friends[i] = &newF
			updated++
		}
	}
	r.friends = newFriends

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
	friends, err := r.Find(ctx, f)
	if err != nil {
		return deleted, err
	}
	ids := []uuid.UUID{}
	for friend := range friends {
		if friend.ID == nil {
			l.Error("found friend with nil UUID", "filter:", f)
		}
		ids = append(ids, *friend.ID)
	}

	newFriends := []*model.Friend{}
	for _, friend := range r.friends {
		match := false
		for _, id := range ids {
			if friend.ID != nil && *friend.ID == id {
				match = true
			}
		}
		if !match {
			newFriends = append(newFriends, friend)
		} else {
			deleted++
		}
	}
	r.friends = newFriends

	return deleted, nil
}
