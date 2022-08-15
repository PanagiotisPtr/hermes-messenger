package connection

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	AddConnection(context.Context, uuid.UUID, uuid.UUID) error
	RemoveConnection(context.Context, uuid.UUID, uuid.UUID) error
	GetConnections(context.Context, uuid.UUID) ([]*Connection, error)
}
