package keys

import (
	"crypto/rsa"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type MemoryRepository struct {
	logger   *zap.Logger
	keyStore map[string]KeyPair
}

func ProvideMemoryKeysRepository(logger *zap.Logger) *MemoryRepository {
	return &MemoryRepository{
		logger:   logger,
		keyStore: make(map[string]KeyPair),
	}
}

func (mr *MemoryRepository) StoreKeyPair(keyName string, keyPair KeyPair, ttl time.Duration) error {
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
