package keys

import (
	"crypto/rsa"
	"time"
)

type Repository interface {
	StoreKeyPair(string, KeyPair, time.Duration) error
	GetPublicKey(string) (*rsa.PublicKey, error)
	GetPrivateKey(string) (*rsa.PrivateKey, error)
	GetAllPublicKeys(string) ([]*rsa.PublicKey, error)
	GetAllPrivateKeys(string) ([]*rsa.PrivateKey, error)
}
