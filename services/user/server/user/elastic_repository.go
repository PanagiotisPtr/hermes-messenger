package user

import (
	"context"
	"fmt"
	"strings"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/esutils"
	"go.uber.org/zap"
)

const (
	UsersIndex = "users"
)

type ESRepository struct {
	logger *zap.Logger
	es     *elasticsearch.Client
}

func ProvideESRepository(
	logger *zap.Logger,
	es *elasticsearch.Client,
) Repository {
	return &ESRepository{
		logger: logger,
		es:     es,
	}
}

func (r *ESRepository) AddUser(
	ctx context.Context,
	email string,
) (*User, error) {
	if email == "" {
		return nil, fmt.Errorf("email address is empty")
	}

	u, err := r.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if u != nil {
		return nil, fmt.Errorf("A user with this email address already exists")
	}

	docId := uuid.New()
	u = &User{
		Uuid:  docId,
		Email: email,
	}
	err = esutils.StoreDocument(
		ctx,
		r.es,
		UsersIndex,
		docId.String(),
		u,
		true,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *ESRepository) GetUser(
	ctx context.Context,
	userUuid uuid.UUID,
) (*User, error) {
	q := fmt.Sprintf(`{
			"from": 0,
			"size": 1,
			"query": {
				"bool": {
					"must": [
						{ "term": { "Uuid.keyword": "%s" } }
					]
				}
			}
		}`, userUuid.String())
	users, err := esutils.GetResults[*User](
		r.es,
		r.es.Search.WithContext(ctx),
		r.es.Search.WithIndex(UsersIndex),
		r.es.Search.WithBody(strings.NewReader(q)),
	)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, fmt.Errorf(
			"could not find user with UUID %s",
			userUuid.String(),
		)
	}

	return users[0], nil
}

func (r *ESRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*User, error) {
	q := fmt.Sprintf(`{
			"from": 0,
			"size": 1,
			"query": {
				"bool": {
					"must": [
						{ "term": { "Email.keyword": "%s" } }
					]
				}
			}
		}`, email)
	users, err := esutils.GetResults[*User](
		r.es,
		r.es.Search.WithContext(ctx),
		r.es.Search.WithIndex(UsersIndex),
		r.es.Search.WithBody(strings.NewReader(q)),
	)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}

	return users[0], nil
}
