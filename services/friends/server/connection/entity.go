package connection

import (
	"github.com/google/uuid"
)

type Connection struct {
	ID     uuid.UUID `bson:"_id" json:"ID"`
	From   uuid.UUID `bson:"From" json:"From"`
	To     uuid.UUID `bson:"To" json:"To"`
	Status string    `bson:"Status" json:"Status"`
}
