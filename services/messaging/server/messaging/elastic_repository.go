package messaging

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	MessageRepoKeyPrefix = "repository:message"
)

type ESRepository struct {
	logger *zap.Logger
	//	redisClient *redis.Client
	es *elasticsearch.Client
}

func ProvideESRepository(
	logger *zap.Logger,
	// rc *redis.Client,
	es *elasticsearch.Client,
) Repository {
	return &ESRepository{
		logger: logger,
		//	redisClient: rc,
		es: es,
	}
}

func keyForUUID(id uuid.UUID) string {
	return MessageRepoKeyPrefix + ":" + id.String()
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

	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      "messages",
		DocumentID: docId.String(),
		Body:       bytes.NewReader(b),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, r.es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("[%s] Error indexing document ID=%s", res.Status(), docId)
	}

	return nil
}

func (r *ESRepository) GetMessages(
	ctx context.Context,
	from uuid.UUID,
	to uuid.UUID,
	start time.Time,
	end time.Time,
) ([]*Message, error) {
	messages := make([]*Message, 0)
	var buf bytes.Buffer
	q := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{"term": map[string]interface{}{"From.keyword": from}},
					map[string]interface{}{"term": map[string]interface{}{"To.keyword": to}},
					map[string]interface{}{"range": map[string]interface{}{
						"Timestamp": map[string]interface{}{
							"gte": start.Unix(),
							"lte": time.Now().Unix(),
						},
					}},
				},
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(q); err != nil {
		return messages, err
	}

	// Perform the search request.
	res, err := r.es.Search(
		r.es.Search.WithContext(context.Background()),
		r.es.Search.WithIndex("messages"),
		r.es.Search.WithBody(&buf),
		r.es.Search.WithTrackTotalHits(true),
		r.es.Search.WithPretty(),
	)
	if err != nil {
		return messages, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return messages, fmt.Errorf(
				`failed to perform search.
				Error: can't display error:
				error parsing the response body: %s`,
				err,
			)
		} else {
			// Print the response status and error information.
			return messages, fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	var rp map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&rp); err != nil {
		return messages, fmt.Errorf("Error parsing the response body: %s", err)
	}

	// Print the ID and document source for each hit.
	for _, hit := range rp["hits"].(map[string]interface{})["hits"].([]interface{}) {
		log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
		source := hit.(map[string]interface{})["_source"]

		messageUuid, err := uuid.Parse(source.(map[string]interface{})["Uuid"].(string))
		if err != nil {
			r.logger.Sugar().Error(err)
			continue
		}
		fromUuid, err := uuid.Parse(source.(map[string]interface{})["From"].(string))
		if err != nil {
			r.logger.Sugar().Error(err)
			continue
		}
		toUuid, err := uuid.Parse(source.(map[string]interface{})["To"].(string))
		if err != nil {
			r.logger.Sugar().Error(err)
			continue
		}
		timestamp := int64(source.(map[string]interface{})["Timestamp"].(float64))
		content := source.(map[string]interface{})["Content"].(string)
		messages = append(messages, &Message{
			Uuid:      messageUuid,
			From:      fromUuid,
			To:        toUuid,
			Timestamp: timestamp,
			Content:   content,
		})
	}

	return messages, nil
}
