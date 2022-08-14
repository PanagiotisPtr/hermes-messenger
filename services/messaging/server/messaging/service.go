package messaging

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	logger      *zap.Logger
	messageRepo Repository
}

func ProvideMessagingService(
	logger *zap.Logger,
	messageRepo Repository,
) *Service {
	return &Service{
		logger:      logger,
		messageRepo: messageRepo,
	}
}

func (s *Service) SendMessage(
	ctx context.Context,
	from uuid.UUID,
	to uuid.UUID,
	content string,
) error {
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
		return ms[i].Timestamp < ms[j].Timestamp
	})

	sort.Slice(left, func(i, j int) bool {
		return ms[i].Timestamp < ms[j].Timestamp
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
