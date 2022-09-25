package messaging

import "github.com/google/uuid"

type Message struct {
	ID        uuid.UUID `bson:"_id" json:"ID"`
	From      uuid.UUID `bson:"From" json:"From"`
	To        uuid.UUID `bson:"To" json:"To"`
	Timestamp int64     `bson:"Timestamp" json:"Timestamp"`
	Content   string    `bson:"Content" json:"Content"`
}
