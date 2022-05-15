package messaging

import (
	"log"
	"sort"

	"github.com/google/uuid"
)

type MemoryRepository struct {
	logger   *log.Logger
	messages []Message
}

func NewMemoryRepository(l *log.Logger) *MemoryRepository {
	return &MemoryRepository{
		logger:   l,
		messages: make([]Message, 0),
	}
}

func (mr *MemoryRepository) SaveMessage(m Message) error {
	mr.messages = append(mr.messages, m)

	return nil
}

func (mr *MemoryRepository) GetMessagesBetweenUsers(
	from uuid.UUID,
	to uuid.UUID,
) ([]Message, error) {
	list := make([]Message, 0)
	for _, message := range mr.messages {
		if message.From == from || message.To == to {
			list = append(list, message)
		}
		if message.From == to || message.To == from {
			list = append(list, message)
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Timestamp < list[j].Timestamp
	})

	return list, nil
}
