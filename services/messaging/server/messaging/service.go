package messaging

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/protos"
	"github.com/panagiotisptr/hermes-messenger/services/friends/server/connection/status"
	"go.uber.org/zap"
)

type Service struct {
	logger        *zap.Logger
	messageRepo   Repository
	friendsClient protos.FriendsClient
}

func ProvideMessagingService(
	logger *zap.Logger,
	messageRepo Repository,
	friendsClient protos.FriendsClient,
) *Service {
	return &Service{
		logger:        logger,
		messageRepo:   messageRepo,
		friendsClient: friendsClient,
	}
}

// Check if two users are friends, returns error if they are not
func (s *Service) usersAreFriends(
	ctx context.Context,
	a uuid.UUID,
	b uuid.UUID,
) error {
	resp, err := s.friendsClient.GetFriends(ctx, &protos.GetFriendsRequest{
		UserUuid: a.String(),
	})
	if err != nil {
		return err
	}

	friendsUuid := b.String()
	for _, f := range resp.Friends {
		if f.Status != status.Accepted {
			continue
		}
		if f.FriendUuid == friendsUuid {
			return nil
		}
	}

	return fmt.Errorf(
		"Users %s and %s are not friends",
		a.String(),
		b.String(),
	)
}

func (s *Service) SendMessage(
	ctx context.Context,
	from uuid.UUID,
	to uuid.UUID,
	content string,
) error {
	if err := s.usersAreFriends(ctx, from, to); err != nil {
		return err
	}
	return s.messageRepo.SaveMessage(
		ctx,
		from,
		to,
		content,
	)
}

func (s *Service) GetMessagesBetweenUsers(
	ctx context.Context,
	from uuid.UUID,
	to uuid.UUID,
	size int64,
	offset int64,
) ([]*Message, error) {
	ms := make([]*Message, 0)
	if err := s.usersAreFriends(ctx, from, to); err != nil {
		return ms, err
	}

	return s.messageRepo.GetMessages(
		ctx,
		from,
		to,
		size,
		offset,
	)
}
