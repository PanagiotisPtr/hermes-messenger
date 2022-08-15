package friends

import (
	"context"

	"github.com/google/uuid"

	"github.com/panagiotisptr/hermes-messenger/friends/server/connection"
	"go.uber.org/zap"
)

type Friend struct {
	FriendUuid uuid.UUID
	Status     string
}

type Service struct {
	logger   *zap.Logger
	connRepo connection.Repository
}

func ProvideFriendsService(
	logger *zap.Logger,
	connRepo connection.Repository,
) *Service {
	return &Service{
		logger:   logger,
		connRepo: connRepo,
	}
}

func (s *Service) AddFriend(
	ctx context.Context,
	userUuid uuid.UUID,
	friendUuid uuid.UUID,
) error {
	return s.connRepo.AddConnection(ctx, userUuid, friendUuid)
}

func (s *Service) RemoveFriend(
	ctx context.Context,
	userUuid uuid.UUID,
	friendUuid uuid.UUID,
) error {
	return s.connRepo.RemoveConnection(ctx, userUuid, friendUuid)
}

func (s *Service) GetFriends(
	ctx context.Context,
	userUuid uuid.UUID,
) ([]*Friend, error) {
	fs := make([]*Friend, 0)
	connections, err := s.connRepo.GetConnections(ctx, userUuid)
	if err != nil {
		return fs, err
	}

	for _, c := range connections {
		fs = append(fs, &Friend{
			FriendUuid: c.From,
			Status:     c.Status,
		})
	}

	return fs, nil
}
