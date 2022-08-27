package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	"github.com/google/uuid"
)

const (
	RefreshKeyType = "refresh"
	AccessKeyType  = "access"
)

type KeyPair struct {
	Uuid       uuid.UUID
	keyType    string
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func GenerateRSAKeyPair(keyType string) (KeyPair, error) {
	kid := uuid.New()
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return KeyPair{
			Uuid:       kid,
			publicKey:  nil,
			privateKey: nil,
		}, fmt.Errorf("Could not generate RSA key pair")
	}
	publicKey := &privateKey.PublicKey

	return KeyPair{
		Uuid:       kid,
		keyType:    keyType,
		publicKey:  publicKey,
		privateKey: privateKey,
	}, nil
}
