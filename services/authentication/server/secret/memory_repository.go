package secret

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
)

type MemoryRepository struct {
	logger   *log.Logger
	keyStore map[string]KeyPair
}

func NewMemoryRepository(logger *log.Logger) *MemoryRepository {
	return &MemoryRepository{
		logger:   logger,
		keyStore: make(map[string]KeyPair),
	}
}

func (mr *MemoryRepository) GenerateRSAKeyPair() (KeyPair, error) {
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

func (mr *MemoryRepository) StoreKeyPair(keyName string, keyPair KeyPair) error {
	if _, ok := mr.keyStore[keyName]; ok {
		return fmt.Errorf("There's already a key stored under the name %s", keyName)
	}
	mr.keyStore[keyName] = keyPair

	return nil
}

func (mr *MemoryRepository) getKeyPair(keyName string) (KeyPair, error) {
	keyPair, ok := mr.keyStore[keyName]
	if !ok {
		return KeyPair{
			publicKey:  nil,
			privateKey: nil,
		}, fmt.Errorf("Could not find key pair under the name %s", keyName)
	}

	return keyPair, nil
}

func (mr *MemoryRepository) GetPublicKey(keyName string) (*rsa.PublicKey, error) {
	keyPair, err := mr.getKeyPair(keyName)
	if err != nil {
		return keyPair.publicKey, err
	}

	return keyPair.publicKey, nil
}

func (mr *MemoryRepository) GetPrivateKey(keyName string) (*rsa.PrivateKey, error) {
	keyPair, err := mr.getKeyPair(keyName)
	if err != nil {
		return keyPair.privateKey, err
	}

	return keyPair.privateKey, nil
}
