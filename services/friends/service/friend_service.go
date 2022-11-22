package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/friends/model"
	"github.com/panagiotisptr/hermes-messenger/friends/repository"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/entityutils/filter"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/loggingutils"
	"go.uber.org/zap"
)

// FriendService represents a friend service
type FriendService struct {
	logger     *zap.Logger
	friendRepo repository.FriendRepository
}

// ProvideFriendService provides an instance of the friend
// service
func ProvideFriendService(
	logger *zap.Logger,
	friendRepo repository.FriendRepository,
) *FriendService {
	return &FriendService{
		logger: logger.With(
			zap.String("service", "FriendService"),
		),
		friendRepo: friendRepo,
	}
}

func (s *FriendService) AddFriend(
	ctx context.Context,
	args *model.Friend,
) (*model.Friend, error) {
	l := loggingutils.
		LoggerWithRequestID(ctx, s.logger).
		With(
			zap.String("method", "AddFriend"),
		).Sugar()

	l.Info("hello")

	return nil, nil
}

// CreateFriend creates a new friend
func (s *FriendService) CreateFriend(
	ctx context.Context,
	args *model.Friend,
) (*model.Friend, error) {
	return s.friendRepo.Create(ctx, args)
}

// GetFriend returns a friend from their (uu)id
func (s *FriendService) GetFriend(
	ctx context.Context,
	id uuid.UUID,
) (*model.Friend, error) {
	f := filter.NewFilter()
	f.Add("id", filter.Eq, id)

	return s.friendRepo.FindOne(ctx, f)
}

func (s *FriendService) GetFriendByEmail(
	ctx context.Context,
	email string,
) (*model.Friend, error) {
	f := filter.NewFilter()
	f.Add("email", filter.Eq, email)

	u, err := s.friendRepo.FindOne(ctx, f)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, fmt.Errorf("Could not find friend with email \"%s\"", email)
	}

	return u, nil
}
