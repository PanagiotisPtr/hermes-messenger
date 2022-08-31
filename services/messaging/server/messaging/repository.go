package messaging

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	SaveMessage(context.Context, uuid.UUID, uuid.UUID, string) error
	GetMessages(context.Context, uuid.UUID, uuid.UUID, int64, int64) ([]*Message, error)
}
