package connection

import (
	"context"
	"encoding/json"
	"fmt"

	redis "github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/friends/server/connection/status"
	"go.uber.org/zap"
)

const (
	ConnectionRepoKeyPrefix = "repository:connection"
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

func keyForConnection(c *Connection) string {
	return fmt.Sprintf("%s:%s:%s",
		ConnectionRepoKeyPrefix,
		c.From.String(),
		c.To.String(),
	)
}

func (r *RedisRepository) setConnection(
	ctx context.Context,
	c *Connection,
) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return r.redisClient.SetNX(
		ctx,
		keyForConnection(c),
		b,
		redis.KeepTTL,
	).Err()
}

func (r *RedisRepository) deleteConnection(
	ctx context.Context,
	c *Connection,
) error {
	return r.redisClient.Del(
		ctx,
		keyForConnection(c),
	).Err()
}

func (r *RedisRepository) getConnectionFromRedis(
	ctx context.Context,
	redisKey string,
) (*Connection, error) {
	cs, err := r.redisClient.Get(
		ctx,
		redisKey,
	).Result()

	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var c Connection
	err = json.Unmarshal([]byte(cs), &c)

	return &c, err
}

func (r *RedisRepository) getConnection(
	ctx context.Context,
	from string,
	to string,
) (*Connection, error) {
	return r.getConnectionFromRedis(
		ctx,
		ConnectionRepoKeyPrefix+":"+from+":"+to,
	)
}

func (r *RedisRepository) AddConnection(
	ctx context.Context,
	from uuid.UUID,
	to uuid.UUID,
) error {
	s := status.Pending
	// Has the other person requested to be our frined
	op, err := r.getConnection(ctx, to.String(), from.String())
	if err != nil {
		return err
	}
	if op != nil {
		// if so, then accept
		op.Status = status.Accepted
		s = status.Accepted
		err = r.deleteConnection(ctx, op)
		if err != nil {
			return nil
		}
		err = r.setConnection(ctx, op)
		if err != nil {
			return err
		}
	}

	return r.setConnection(ctx, &Connection{
		From:   from,
		To:     to,
		Status: s,
	})
}

func (r *RedisRepository) GetConnections(
	ctx context.Context,
	to uuid.UUID,
) ([]*Connection, error) {
	cs := make([]*Connection, 0)
	iter := r.redisClient.Scan(
		ctx,
		0,
		ConnectionRepoKeyPrefix+":*:"+to.String(),
		0,
	).Iterator()
	for iter.Next(ctx) {
		c, err := r.getConnectionFromRedis(ctx, iter.Val())
		if err != nil {
			r.logger.Sugar().Error(err)
			continue
		}
		cs = append(cs, c)
	}

	return cs, iter.Err()
}

func (r *RedisRepository) RemoveConnection(
	ctx context.Context,
	from uuid.UUID,
	to uuid.UUID,
) error {
	left, err := r.getConnection(ctx, from.String(), to.String())
	if err != nil {
		return err
	}
	if left != nil {
		err = r.deleteConnection(ctx, left)
		if err != nil {
			return err
		}
	}

	right, err := r.getConnection(ctx, to.String(), from.String())
	if err != nil {
		return err
	}
	if right != nil {
		err = r.deleteConnection(ctx, right)
		if err != nil {
			return err
		}
	}

	return nil
}
