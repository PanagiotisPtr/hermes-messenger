package messaging

import (
	"context"
	"fmt"
	"sort"
	"time"

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
	start time.Time,
	end time.Time,
) ([]*Message, error) {
	ms := make([]*Message, 0)
	if err := s.usersAreFriends(ctx, from, to); err != nil {
		return ms, err
	}
	left, err := s.messageRepo.GetMessages(
		ctx,
		from,
		to,
		start,
		end,
	)
	if err != nil {
		return ms, err
	}

	right, err := s.messageRepo.GetMessages(
		ctx,
		to,
		from,
		start,
		end,
	)
	if err != nil {
		return ms, err
	}

	// Probably don't need to sort the results here
	sort.Slice(left, func(i, j int) bool {
		return left[i].Timestamp < left[j].Timestamp
	})

	sort.Slice(right, func(i, j int) bool {
		return right[i].Timestamp < right[j].Timestamp
	})

	i := 0
	j := 0
	for i < len(left) && j < len(right) {
		if left[i].Timestamp < right[j].Timestamp {
			ms = append(ms, left[i])
			i++
		} else {
			ms = append(ms, right[j])
			j++
		}
	}

	for i < len(left) {
		ms = append(ms, left[i])
		i++
	}

	for j < len(right) {
		ms = append(ms, right[j])
		j++
	}

	return ms, nil
}
