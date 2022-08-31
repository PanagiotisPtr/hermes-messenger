package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	redis "github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/esutils"
	"go.uber.org/zap"
)

const (
	MessagesIndex = "messages"
)

type ESRepository struct {
	logger      *zap.Logger
	redisClient *redis.Client
	es          *elasticsearch.Client
}

func ProvideESRepository(
	logger *zap.Logger,
	rc *redis.Client,
	es *elasticsearch.Client,
) Repository {
	return &ESRepository{
		logger:      logger,
		redisClient: rc,
		es:          es,
	}
}

func (r *ESRepository) SaveMessage(
	ctx context.Context,
	from uuid.UUID,
	to uuid.UUID,
	content string,
) error {
	docId := uuid.New()
	m := &Message{
		Uuid:      docId,
		From:      from,
		To:        to,
		Timestamp: time.Now().Unix(),
		Content:   content,
	}

	err := esutils.StoreDocument(
		ctx,
		r.es,
		MessagesIndex,
		docId.String(),
		m,
		true,
	)
	if err != nil {
		return err
	}

	b, err := json.Marshal(m)
	if err != nil {
		r.logger.Sugar().Warnf("Failed to publish message to redis channel: %s", err.Error())
		return nil
	}
	err = r.redisClient.Publish(ctx, "messages", string(b)).Err()
	if err != nil {
		r.logger.Sugar().Warnf("Failed to publish message to redis channel: %s", err.Error())
	}

	return nil
}

func (r *ESRepository) GetMessages(
	ctx context.Context,
	from uuid.UUID,
	to uuid.UUID,
	size int64,
	offset int64,
) ([]*Message, error) {
	q := fmt.Sprintf(`{
	    "from": %d,
	    "size": %d,
	    "sort": {
	        "Timestamp": {
	            "order": "desc"
	        }
	    },
	    "query": {
	        "bool": {
	            "should": [
	                {
	                    "bool": {
	                        "must": [
	                            { "term": { "From.keyword": "%s" } },
	                            { "term": { "To.keyword": "%s" } }
	                        ]
	                    }
	                },
	                {
	                    "bool": {
	                        "must": [
	                            { "term": { "From.keyword": "%s" } },
	                            { "term": { "To.keyword": "%s" } }
	                        ]
	                    }
	                }
	            ]
	        }
	    }
	}`, offset, size, from, to, to, from)

	return esutils.GetResults[*Message](
		r.es,
		r.es.Search.WithContext(ctx),
		r.es.Search.WithIndex(MessagesIndex),
		r.es.Search.WithBody(strings.NewReader(q)),
	)
}
