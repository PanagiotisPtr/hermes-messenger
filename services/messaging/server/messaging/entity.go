package messaging

import "github.com/google/uuid"

type Message struct {
	From      uuid.UUID
	To        uuid.UUID
	Timestamp int64
	Content   string
}
