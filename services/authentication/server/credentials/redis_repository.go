package credentials

import (
	"context"
	"encoding/json"

	redis "github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	CredentialsRepoKeyPrefix = "repository:credentials"
)

type RedisRepository struct {
	logger      *zap.Logger
	redisClient *redis.Client
}

func ProvideRedisRepository(
	logger *zap.Logger,
	rc *redis.Client,
) Repository {
	return &RedisRepository{
		logger:      logger,
		redisClient: rc,
	}
}

func keyForUuid(id uuid.UUID) string {
	return CredentialsRepoKeyPrefix + ":" + id.String()
}

func (r *RedisRepository) CreateCredentials(
	ctx context.Context,
	userUuid uuid.UUID,
	password string,
) error {
	c := &Credentials{
		Uuid:     userUuid,
		Password: password,
	}
	bytes, err := json.Marshal(c)
	if err != nil {
		return err
	}

	_, err = r.redisClient.SetNX(
		ctx,
		keyForUuid(c.Uuid),
		string(bytes),
		redis.KeepTTL,
	).Result()

	if err != nil {
		return err
	}

	return nil
}

func (r *RedisRepository) GetCredentials(
	ctx context.Context,
	userUuid uuid.UUID,
) (*Credentials, error) {
	cJson, err := r.redisClient.Get(
		ctx,
		keyForUuid(userUuid),
	).Result()
	if err != nil {
		return nil, err
	}

	var c Credentials
	err = json.Unmarshal([]byte(cJson), &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
