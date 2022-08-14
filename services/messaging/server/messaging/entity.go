package messaging

import "github.com/google/uuid"

type Message struct {
	Uuid      uuid.UUID `json:"Uuid"`
	From      uuid.UUID `json:"From"`
	To        uuid.UUID `json:"To"`
	Timestamp int64     `json:"Timestamp"`
	Content   string    `json:"Content"`
}
