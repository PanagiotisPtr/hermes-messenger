package esutils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type ShardsData struct {
	Total      int64 `json:"total"`
	Successful int64 `json:"successful"`
	Skipped    int64 `json:"skipped"`
	Failed     int64 `json:"failed"`
}

type Hit[T any] struct {
	Index  string  `json:"_index"`
	ID     string  `json:"_id"`
	Score  float64 `json:"_score"`
	Source T       `json:"_source"`
}

type TotalData struct {
	Value    int64  `json:"value"`
	Relation string `json:"relation"`
}

type HitsData[T any] struct {
	MaxScore float64    `json:"max_score"`
	Hits     []Hit[T]   `json:"hits"`
	Total    *TotalData `json:"total"`
}

type ESResult[T any] struct {
	Took    int64       `json:"took"`
	TimeOut bool        `json:"timed_out"`
	Shards  ShardsData  `json:"_shards"`
	Hits    HitsData[T] `json:"hits"`
}

type ErrorData struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

type ESError struct {
	Error ErrorData `json:"error"`
}

// StoreDocument stores a single document on an ES index
func StoreDocument(
	ctx context.Context,
	es *elasticsearch.Client,
	index string,
	documentID string,
	document interface{},
	refresh bool,
) error {
	b, err := json.Marshal(document)
	if err != nil {
		return err
	}

	refreshString := "false"
	if refresh {
		refreshString = "true"
	}
	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: documentID,
		Body:       bytes.NewReader(b),
		Refresh:    refreshString,
	}
	res, err := req.Do(ctx, es)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf(
			"[%s] Error Indexing document ID=%s",
			res.Status(),
			documentID,
		)
	}

	return nil
}

// GetResults returns the search results from an ES search
func GetResults[T any](
	es *elasticsearch.Client,
	o ...func(*esapi.SearchRequest),
) ([]T, error) {
	rv := make([]T, 0)
	res, err := es.Search(o...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e ESError
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return rv, fmt.Errorf(
				`failed to perform search.
			Error: can't display error:
			error parsing the response body: %s`,
				err,
			)
		} else {
			// Print the response status and error information.
			return rv, fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e.Error.Type,
				e.Error.Reason,
			)
		}
	}

	var rp ESResult[T]
	if err := json.NewDecoder(res.Body).Decode(&rp); err != nil {
		return nil, fmt.Errorf("Error parsing the response body: %s", err)
	}

	for _, hit := range rp.Hits.Hits {
		rv = append(rv, hit.Source)
	}

	return rv, nil
}
