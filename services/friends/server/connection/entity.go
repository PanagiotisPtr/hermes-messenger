package connection

import (
	"github.com/google/uuid"
)

type Connection struct {
	From   uuid.UUID `json:"From"`
	To     uuid.UUID `json:"To"`
	Status string    `json:"Status"`
}
