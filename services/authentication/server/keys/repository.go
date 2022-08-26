package keys

import (
	"context"
	"crypto/rsa"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	StoreKeyPair(context.Context, KeyPair, time.Duration) error
	GetPublicKey(context.Context, uuid.UUID, string) (*rsa.PublicKey, error)
	GetPrivateKey(context.Context, uuid.UUID, string) (*rsa.PrivateKey, error)
	GetAllPublicKeys(context.Context, string) (map[uuid.UUID]*rsa.PublicKey, error)
}
