package redis_repository

import (
	"context"
	"crypto/x509"
	"strings"

	"github.com/pkg/errors"

	"github.com/golang-jwt/jwt"
	"github.com/panagiotisptr/hermes-messenger/services/authentication/repository"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	PrivateKeyNotFound = errors.New("private key not found")
	PublicKeyNotFound  = errors.New("public key not found")
	InvalidToken       = errors.New("invalid token")
)

type RedisTokenRepository struct {
	logger *zap.Logger
	rc     *redis.Client
}

func ProvideTokenRepository(
	logger *zap.Logger,
) repository.TokenRepository {
	return &RedisTokenRepository{
		logger: logger.With(
			zap.String("repository", "TokenRepository"),
			zap.String("type", "redis"),
		),
	}
}

func (r *RedisTokenRepository) SignTokenWithClaims(
	ctx context.Context,
	claims jwt.Claims,
) (string, error) {
	for {
		privateKeyUUID, err := r.rc.SRandMember(ctx, PrivateKeyUUIDsSetName).Result()
		if err != nil {
			if err == redis.Nil {
				return "", PrivateKeyNotFound
			}
			return "", errors.Wrap(err, "getting random set member from private keys set")
		}

		keyName := strings.Join([]string{PrivateKeyPrefix, privateKeyUUID}, ":")
		privateKeyString, err := r.rc.Get(
			ctx,
			keyName,
		).Result()
		if err == redis.Nil {
			// try again with another key
			continue
		}
		if err != nil {
			return "", errors.Wrap(err, "getting private key with UUID")
		}

		privateKey, err := x509.ParsePKCS1PrivateKey([]byte(privateKeyString))
		if err != nil {
			return "", errors.Wrap(err, "failed to parse private key")
		}

		token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privateKey)

		return token, errors.Wrap(err, "failed to sign token")
	}
}

func (r *RedisTokenRepository) ValidateToken(
	ctx context.Context,
	tokenString string,
) (jwt.Claims, error) {
	publicKeyUUIDs, err := r.rc.SMembers(ctx, PublicKeyUUIDsSetName).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, PublicKeyNotFound
		}
		return nil, errors.Wrap(err, "getting public keys from redis set")
	}
	for _, publicKeyUUID := range publicKeyUUIDs {
		keyName := strings.Join([]string{PublicKeyPrefix, publicKeyUUID}, ":")
		publicKeyString, err := r.rc.Get(
			ctx,
			keyName,
		).Result()
		if err == redis.Nil {
			// try again with another key
			continue
		}
		if err != nil {
			return nil, errors.Wrap(err, "getting public key with UUID")
		}

		publicKey, err := x509.ParsePKCS1PublicKey([]byte(publicKeyString))
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse public key")
		}

		token, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
			if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.Errorf("unexpected method: %s", jwtToken.Header["alg"])
			}

			return publicKey, nil
		})
		if err != nil {
			// failed to validate token
			continue
		}

		if token.Claims.Valid() != nil {
			return nil, errors.Errorf("Invalid token claims")
		}

		return token.Claims, nil
	}

	return nil, InvalidToken
}
