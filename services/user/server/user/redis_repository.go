package user

import (
	"context"
	"encoding/json"

	redis "github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	UserRepoKeyPrefix = "repository:user"
)

type RedisRepository struct {
	logger      *zap.Logger
	redisClient *redis.Client
}

func ProvideRedisKeysRepository(
	logger *zap.Logger,
	rc *redis.Client,
) Repository {
	return &RedisRepository{
		logger:      logger,
		redisClient: rc,
	}
}

func (r *RedisRepository) AddUser(email string) (string, error) {
	user := User{
		Uuid:  uuid.New(),
		Email: email,
	}
	bytes, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	_, err = r.redisClient.SetNX(
		context.Background(),
		UserRepoKeyPrefix+":"+user.Uuid.String(),
		string(bytes),
		redis.KeepTTL,
	).Result()

	if err != nil {
		return "", err
	}

	return user.Uuid.String(), nil
}
