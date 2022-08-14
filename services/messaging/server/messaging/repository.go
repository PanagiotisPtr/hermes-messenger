package messaging

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	SaveMessage(context.Context, uuid.UUID, uuid.UUID, string) error
	GetMessages(context.Context, uuid.UUID, uuid.UUID, time.Time, time.Time) ([]*Message, error)
}
