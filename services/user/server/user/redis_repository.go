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

func ProvideRedisRepository(
	logger *zap.Logger,
	rc *redis.Client,
) Repository {
	return &RedisRepository{
		logger:      logger,
		redisClient: rc,
	}
}

func keyForUUID(id uuid.UUID) string {
	return UserRepoKeyPrefix + ":" + id.String()
}

func (r *RedisRepository) AddUser(ctx context.Context, email string) (*User, error) {
	user := &User{
		Uuid:  uuid.New(),
		Email: email,
	}
	bytes, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	_, err = r.redisClient.SetNX(
		ctx,
		keyForUUID(user.Uuid),
		string(bytes),
		redis.KeepTTL,
	).Result()

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *RedisRepository) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
	userJson, err := r.redisClient.Get(
		ctx,
		keyForUUID(id),
	).Result()
	if err != nil {
		return nil, err
	}

	var user User
	err = json.Unmarshal([]byte(userJson), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
