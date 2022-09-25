package friends

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/panagiotisptr/hermes-messenger/friends/server/connection"
	"github.com/panagiotisptr/hermes-messenger/protos"
	"go.uber.org/zap"
)

type Friend struct {
	FriendUuid uuid.UUID
	Status     string
}

type Service struct {
	logger     *zap.Logger
	connRepo   connection.Repository
	userClient protos.UserClient
}

func ProvideFriendsService(
	logger *zap.Logger,
	connRepo connection.Repository,
	userClient protos.UserClient,
) *Service {
	return &Service{
		logger:     logger,
		connRepo:   connRepo,
		userClient: userClient,
	}
}

// Returns error if any user with the specified uuid doesn't exist
func (s *Service) usersExists(ctx context.Context, ids []uuid.UUID) error {
	for _, id := range ids {
		userResp, err := s.userClient.GetUser(ctx, &protos.GetUserRequest{
			Uuid: id.String(),
		})
		if err != nil {
			return err
		}
		if userResp.User == nil {
			return fmt.Errorf("could not find user with UUID: %s", id.String())
		}
	}

	return nil
}

func (s *Service) AddFriend(
	ctx context.Context,
	userUuid uuid.UUID,
	friendUuid uuid.UUID,
) error {
	// Check that both users exist
	if err := s.usersExists(ctx, []uuid.UUID{userUuid, friendUuid}); err != nil {
		return err
	}

	return s.connRepo.AddConnection(ctx, userUuid, friendUuid)
}

func (s *Service) RemoveFriend(
	ctx context.Context,
	userUuid uuid.UUID,
	friendUuid uuid.UUID,
) error {
	// Check that both users exist
	if err := s.usersExists(ctx, []uuid.UUID{userUuid, friendUuid}); err != nil {
		return err
	}

	return s.connRepo.RemoveConnection(ctx, userUuid, friendUuid)
}

func (s *Service) GetFriends(
	ctx context.Context,
	userUuid uuid.UUID,
) ([]*Friend, error) {
	fs := make([]*Friend, 0)
	// Check that both users exist
	if err := s.usersExists(ctx, []uuid.UUID{userUuid}); err != nil {
		return fs, err
	}
	connections, err := s.connRepo.GetConnections(ctx, userUuid)
	if err != nil {
		return fs, err
	}

	for _, c := range connections {
		friendUuid := c.From
		if friendUuid == userUuid {
			friendUuid = c.To
		}
		fs = append(fs, &Friend{
			FriendUuid: friendUuid,
			Status:     c.Status,
		})
	}

	return fs, nil
}
