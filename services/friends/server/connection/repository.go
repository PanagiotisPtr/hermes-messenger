package connection

import "github.com/google/uuid"

type Repository interface {
	AddConnection(uuid.UUID, uuid.UUID) error
	UpdateConnectionStatus(Connection, string) error
	GetConnection(uuid.UUID, uuid.UUID) (*Connection, error)
	GetConnectionsForUser(uuid.UUID) ([]*Connection, error)
}
