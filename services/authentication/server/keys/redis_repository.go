package keys

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v9"
	"go.uber.org/zap"
)

const (
	PublicKeySuffix  = "-public"
	PrivateKeySuffix = "-private"
)

type RedisRepository struct {
	logger      *zap.Logger
	redisClient *redis.Client
}

func ProvideRedisKeysRepository(logger *zap.Logger, rc *redis.Client) Repository {
	return &RedisRepository{
		logger:      logger,
		redisClient: rc,
	}
}

func (rr *RedisRepository) StoreKeyPair(keyPairName string, keyPair KeyPair, ttl time.Duration) error {
	publicKeyString := x509.MarshalPKCS1PublicKey(keyPair.publicKey)
	privateKeyString := x509.MarshalPKCS1PrivateKey(keyPair.privateKey)
	_, err := rr.redisClient.SetNX(
		context.Background(),
		keyPairName+PublicKeySuffix,
		string(publicKeyString),
		ttl,
	).Result()
	if err != nil {
		return err
	}

	_, err = rr.redisClient.SetNX(
		context.Background(),
		keyPairName+PrivateKeySuffix,
		string(privateKeyString),
		ttl,
	).Result()
	if err != nil {
		return err
	}

	return nil
}

func (rr *RedisRepository) GetPublicKey(keyName string) (*rsa.PublicKey, error) {
	publicKeyString, err := rr.redisClient.Get(context.Background(), keyName).Result()
	if err != nil {
		return nil, fmt.Errorf("Could not find public key pair under the name %s", keyName)
	}
	publicKey, err := x509.ParsePKCS1PublicKey([]byte(publicKeyString))
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

func (rr *RedisRepository) GetPrivateKey(keyName string) (*rsa.PrivateKey, error) {
	privateKeyString, err := rr.redisClient.Get(context.Background(), keyName).Result()
	if err != nil {
		return nil, fmt.Errorf("Could not find private key pair under the name %s", keyName)
	}
	privateKey, err := x509.ParsePKCS1PrivateKey([]byte(privateKeyString))
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func (rr *RedisRepository) getKeysWithPattern(pattern string) ([]string, error) {
	keyStrings := make([]string, 0)
	ctx := context.Background()
	iter := rr.redisClient.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		keyStrings = append(keyStrings, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return keyStrings, err
	}

	return keyStrings, nil
}

func (rr *RedisRepository) GetAllPublicKeys(keyPrefix string) ([]*rsa.PublicKey, error) {
	publicKeys := make([]*rsa.PublicKey, 0)
	publicKeyStrings, err := rr.getKeysWithPattern(keyPrefix + "*" + PublicKeySuffix)
	if err != nil {
		return publicKeys, err
	}
	for _, publicKeyName := range publicKeyStrings {
		publicKey, err := rr.GetPublicKey(publicKeyName)
		if err != nil {
			return publicKeys, err
		}
		publicKeys = append(publicKeys, publicKey)
	}

	return publicKeys, nil
}

func (rr *RedisRepository) GetAllPrivateKeys(keyPrefix string) ([]*rsa.PrivateKey, error) {
	privateKeys := make([]*rsa.PrivateKey, 0)
	privateKeyStrings, err := rr.getKeysWithPattern(keyPrefix + "*" + PrivateKeySuffix)
	if err != nil {
		return privateKeys, err
	}
	for _, privateKeyName := range privateKeyStrings {
		privateKey, err := rr.GetPrivateKey(privateKeyName)
		if err != nil {
			return privateKeys, err
		}
		privateKeys = append(privateKeys, privateKey)
	}

	return privateKeys, nil
}
