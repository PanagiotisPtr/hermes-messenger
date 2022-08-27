package keys

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"strings"
	"time"

	redis "github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	PublicKeySuffix   = "public"
	PrivateKeySuffix  = "private"
	KeysRepoKeyPrefix = "repository:keys"
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

func getKeyPairName(keyPair KeyPair) string {
	return KeysRepoKeyPrefix + ":" + keyPair.keyType + ":" + keyPair.Uuid.String()
}

func (rr *RedisRepository) StoreKeyPair(
	ctx context.Context,
	keyPair KeyPair,
	ttl time.Duration,
) error {
	publicKeyString := x509.MarshalPKCS1PublicKey(keyPair.publicKey)
	privateKeyString := x509.MarshalPKCS1PrivateKey(keyPair.privateKey)
	_, err := rr.redisClient.SetNX(
		ctx,
		getKeyPairName(keyPair)+":"+PublicKeySuffix,
		string(publicKeyString),
		ttl,
	).Result()
	if err != nil {
		return err
	}

	_, err = rr.redisClient.SetNX(
		ctx,
		getKeyPairName(keyPair)+":"+PrivateKeySuffix,
		string(privateKeyString),
		ttl,
	).Result()
	if err != nil {
		return err
	}

	return nil
}

func (rr *RedisRepository) GetPublicKey(
	ctx context.Context,
	kid uuid.UUID,
	keyType string,
) (*rsa.PublicKey, error) {
	keyName := KeysRepoKeyPrefix + ":" + keyType + ":" + kid.String() + ":" + PublicKeySuffix
	publicKeyString, err := rr.redisClient.Get(
		ctx,
		keyName,
	).Result()
	if err != nil {
		return nil, fmt.Errorf("Could not find public key pair under the name %s", keyName)
	}
	publicKey, err := x509.ParsePKCS1PublicKey([]byte(publicKeyString))
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

func (rr *RedisRepository) GetPrivateKey(
	ctx context.Context,
	kid uuid.UUID,
	keyType string,
) (*rsa.PrivateKey, error) {
	keyName := KeysRepoKeyPrefix + ":" + keyType + ":" + kid.String() + ":" + PrivateKeySuffix
	privateKeyString, err := rr.redisClient.Get(ctx, keyName).Result()
	if err != nil {
		return nil, fmt.Errorf("Could not find private key pair under the name %s", keyName)
	}
	privateKey, err := x509.ParsePKCS1PrivateKey([]byte(privateKeyString))
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func (rr *RedisRepository) getKeysWithPattern(
	ctx context.Context,
	pattern string,
) ([]string, error) {
	keyStrings := make([]string, 0)
	rr.logger.Sugar().Infof("PATTERN: %s", pattern)
	iter := rr.redisClient.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		rr.logger.Sugar().Infof("FOUND KEY: %s", iter.Val())
		keyStrings = append(keyStrings, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return keyStrings, err
	}

	return keyStrings, nil
}

func (rr *RedisRepository) GetAllPublicKeys(
	ctx context.Context,
	keyType string,
) (map[uuid.UUID]*rsa.PublicKey, error) {
	keyPattern := KeysRepoKeyPrefix + ":" + keyType + ":*:"
	publicKeys := make(map[uuid.UUID]*rsa.PublicKey, 0)
	publicKeyStrings, err := rr.getKeysWithPattern(ctx, keyPattern+PublicKeySuffix)
	if err != nil {
		return publicKeys, err
	}
	for _, publicKeyName := range publicKeyStrings {
		rr.logger.Sugar().Infof("KEY NAME: %s", publicKeyName)
		// 					0 	   1	2 	   3  	     4
		// KeyFormat: repository:keys:KeyType:KID:(Public/Private)
		parts := strings.Split(publicKeyName, ":")
		kid, err := uuid.Parse(parts[3])
		if err != nil {
			rr.logger.Sugar().Error(fmt.Errorf(
				"failed to parse uuid of public key %s. Reason: %v",
				publicKeyName,
				err,
			))
		}
		publicKey, err := rr.GetPublicKey(ctx, kid, parts[2])
		if err != nil {
			return publicKeys, err
		}
		publicKeys[kid] = publicKey
	}

	return publicKeys, nil
}
