package connection

import (
	"fmt"
	"log"

	"github.com/panagiotisptr/hermes-messenger/friends/server/connection/status"

	"github.com/google/uuid"
)

type MemoryRepository struct {
	logger      *log.Logger
	connections []*Connection
}

func NewMemoryRepository(logger *log.Logger) *MemoryRepository {
	return &MemoryRepository{
		logger:      logger,
		connections: make([]*Connection, 0),
	}
}

func (mr *MemoryRepository) AddConnection(from uuid.UUID, to uuid.UUID) error {
	mr.connections = append(mr.connections, &Connection{
		From:   from,
		To:     to,
		Status: status.Pending,
	})

	return nil
}

func (mr *MemoryRepository) UpdateConnectionStatus(c Connection, status string) error {
	for _, connection := range mr.connections {
		if connection.From == c.From && connection.To == c.To {
			connection.Status = status
			return nil
		}
		if connection.From == c.To && connection.To == c.From {
			connection.Status = status
			return nil
		}
	}

	return fmt.Errorf("Connection from %s to %s not found.", c.From.String(), c.To.String())
}

func (mr *MemoryRepository) GetConnection(from uuid.UUID, to uuid.UUID) (*Connection, error) {
	for _, connection := range mr.connections {
		if connection.From == from && connection.To == to {
			return connection, nil
		}
		if connection.From == to && connection.To == from {
			return connection, nil
		}
	}

	return nil, fmt.Errorf("Connection from %s to %s not found.", from.String(), to.String())
}

func (mr *MemoryRepository) GetConnectionsForUser(userUuid uuid.UUID) ([]*Connection, error) {
	result := make([]*Connection, 0)
	for _, connection := range mr.connections {
		if connection.From == userUuid || connection.To == userUuid {
			result = append(result, connection)
		}
	}

	return result, nil
}
