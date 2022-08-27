package credentials

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateCredentials(context.Context, uuid.UUID, string) error
	GetCredentials(context.Context, uuid.UUID) (*Credentials, error)
}
