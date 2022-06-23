package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

type KeyPair struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func GenerateRSAKeyPair() (KeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return KeyPair{
			publicKey:  nil,
			privateKey: nil,
		}, fmt.Errorf("Could not generate RSA key pair")
	}
	publicKey := &privateKey.PublicKey

	return KeyPair{
		publicKey:  publicKey,
		privateKey: privateKey,
	}, nil
}
