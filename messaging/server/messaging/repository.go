package messaging

import "github.com/google/uuid"

type Repository interface {
	SaveMessage(Message) error
	GetMessagesBetweenUsers(uuid.UUID, uuid.UUID) ([]Message, error)
}
