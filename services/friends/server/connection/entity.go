package connection

import (
	"github.com/google/uuid"
)

type Connection struct {
	From   uuid.UUID
	To     uuid.UUID
	Status string
}
