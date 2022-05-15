package messaging

import (
	"log"
	"time"

	"github.com/google/uuid"
)

type MessagingService struct {
	logger *log.Logger
	mr     *MemoryRepository
}

func NewMessagingService(l *log.Logger) *MessagingService {
	return &MessagingService{
		logger: l,
		mr:     NewMemoryRepository(l),
	}
}

func (ms *MessagingService) SendMessage(
	from string,
	to string,
	content string,
) error {
	fromUuid, err := uuid.Parse(from)
	if err != nil {
		return err
	}

	toUuid, err := uuid.Parse(to)
	if err != nil {
		return err
	}

	return ms.mr.SaveMessage(Message{
		From:      fromUuid,
		To:        toUuid,
		Timestamp: time.Now().Unix(),
		Content:   content,
	})
}

func (ms *MessagingService) GetMessagesBetweenUsers(
	from string,
	to string,
) ([]Message, error) {
	messages := make([]Message, 0)

	fromUuid, err := uuid.Parse(from)
	if err != nil {
		return messages, err
	}

	toUuid, err := uuid.Parse(to)
	if err != nil {
		return messages, err
	}

	messages, err = ms.mr.GetMessagesBetweenUsers(fromUuid, toUuid)
	if err != nil {
		return messages, err
	}

	return messages, nil
}
