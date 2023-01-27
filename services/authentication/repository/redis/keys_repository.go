package redis_repository

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"strings"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/services/authentication/repository"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const PublicKeyUUIDsSetName = "auth:publickeyuuids"
const PrivateKeyUUIDsSetName = "auth:privatekeyuuids"
const PrivateKeyPrefix = "auth:privatekeys"
const PublicKeyPrefix = "auth:publickeys"
const MaxNumberOfKeys = 5

type RedisKeysRepository struct {
	logger *zap.Logger
	rc     *redis.Client
}

func ProvideKeysRepository(
	logger *zap.Logger,
) repository.TokenRepository {
	return &RedisTokenRepository{
		logger: logger.With(
			zap.String("repository", "KeysRepository"),
			zap.String("type", "redis"),
		),
	}
}

type keyPair struct {
	kid        uuid.UUID
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func (r *RedisKeysRepository) GetPublicKeys(ctx context.Context) ([]string, error) {
	keys := []string{}
	publicKeyUUIDs, err := r.rc.SMembers(ctx, PublicKeyUUIDsSetName).Result()
	if err != nil {
		if err == redis.Nil {
			return keys, PublicKeyNotFound
		}
		return keys, errors.Wrap(err, "getting public keys from redis set")
	}

	for _, publicKeyUUID := range publicKeyUUIDs {
		keyName := strings.Join([]string{PublicKeyPrefix, publicKeyUUID}, ":")
		publicKeyString, err := r.rc.Get(
			ctx,
			keyName,
		).Result()
		if err == redis.Nil {
			// key expired - it will be re-generated automatically later
			continue
		}
		if err != nil {
			return keys, errors.Wrap(err, "getting public key with UUID")
		}

		keys = append(keys, publicKeyString)
	}

	return keys, nil
}

func (r *RedisKeysRepository) GenerateKeys(ctx context.Context) error {
	err := r.maintainKeysInSet(ctx, PrivateKeyUUIDsSetName, PrivateKeyPrefix)
	if err != nil {
		return err
	}

	err = r.maintainKeysInSet(ctx, PublicKeyUUIDsSetName, PublicKeyPrefix)
	if err != nil {
		return err
	}

	size, err := r.rc.SCard(ctx, PrivateKeyPrefix).Result()
	if err == redis.Nil {
		size = 0
	} else if err != nil {
		return errors.Wrap(err, "could not get set cardinality")
	}

	if size >= MaxNumberOfKeys {
		return nil
	}

	kp, err := generateRSAKeyPair()
	if err != nil {
		return errors.Wrap(err, "failed to generate rsa key pair")
	}

	_, err = r.rc.SAdd(
		ctx,
		PrivateKeyUUIDsSetName,
		kp.kid.String(),
	).Result()
	if err != nil {
		return errors.Wrap(err, "failed to store private key uuid")
	}

	_, err = r.rc.SAdd(
		ctx,
		PublicKeyUUIDsSetName,
		kp.kid.String(),
	).Result()
	if err != nil {
		return errors.Wrap(err, "failed to store public key uuid")
	}

	privateKeyName := strings.Join([]string{PrivateKeyPrefix, kp.kid.String()}, ":")
	privateKeyString := x509.MarshalPKCS1PrivateKey(kp.privateKey)
	_, err = r.rc.Set(
		ctx,
		privateKeyName,
		privateKeyString,
		time.Hour*24,
	).Result()
	if err != nil {
		return errors.Wrap(err, "failed storing private key")
	}

	publicKeyName := strings.Join([]string{PublicKeyPrefix, kp.kid.String()}, ":")
	publicKeyString := x509.MarshalPKCS1PublicKey(kp.publicKey)
	_, err = r.rc.Set(
		ctx,
		publicKeyName,
		publicKeyString,
		time.Hour*24,
	).Result()
	if err != nil {
		return errors.Wrap(err, "failed storing public key")
	}

	return nil
}

func (r *RedisKeysRepository) maintainKeysInSet(
	ctx context.Context,
	setName string,
	keyPrefix string,
) error {
	kids, err := r.rc.SMembers(ctx, setName).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return err
	}

	shouldDelete := []string{}
	for _, kid := range kids {
		keyName := strings.Join([]string{keyPrefix, kid}, ":")
		_, err = r.rc.Get(
			ctx,
			keyName,
		).Result()
		// key expired - need to delete it from set
		if err == redis.Nil {
			shouldDelete = append(shouldDelete, kid)
			continue
		}
		if err != nil {
			return errors.Wrap(err, "failed to get key from redis")
		}
	}

	for _, kid := range shouldDelete {
		_, err = r.rc.SRem(ctx, setName, kid).Result()
		if err != nil {
			return errors.Wrap(err, "failed to delete uuid from set")
		}
	}

	return nil
}

func generateRSAKeyPair() (keyPair, error) {
	kid := uuid.New()
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return keyPair{
			kid:        kid,
			publicKey:  nil,
			privateKey: nil,
		}, err
	}
	publicKey := &privateKey.PublicKey

	return keyPair{
		kid:        kid,
		publicKey:  publicKey,
		privateKey: privateKey,
	}, nil
}
