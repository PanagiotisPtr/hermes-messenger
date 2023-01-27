package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	"github.com/google/uuid"
)

type KeyPair struct {
	kid        uuid.UUID
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func GenerateRSAKeyPair(keyType string) (KeyPair, error) {
	kid := uuid.New()
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return KeyPair{
			kid:        kid,
			publicKey:  nil,
			privateKey: nil,
		}, fmt.Errorf("Could not generate RSA key pair")
	}
	publicKey := &privateKey.PublicKey

	return KeyPair{
		kid:        kid,
		publicKey:  publicKey,
		privateKey: privateKey,
	}, nil
}
