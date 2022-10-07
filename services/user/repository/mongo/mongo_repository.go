package mongo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/user/model"
	"github.com/panagiotisptr/hermes-messenger/user/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	UserCollectionName = "users"
	UserEmailIndex     = "usres_email_index"
)

type MongoRepository struct {
	coll   *mongo.Collection
	logger *zap.Logger
}

func ProvideUserRepository(
	lc fx.Lifecycle,
	logger *zap.Logger,
	database *mongo.Database,
) repository.Repository {
	repo := &MongoRepository{
		logger: logger,
		coll:   database.Collection(UserCollectionName),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Sugar().Info("initialising mongobd indexes for user repository")
			return repo.initIndexes(ctx)
		},
	})

	return repo
}

func (r *MongoRepository) initIndexes(
	ctx context.Context,
) error {
	_, err := r.coll.Indexes().CreateMany(ctx,
		[]mongo.IndexModel{
			{
				Keys: bson.D{{Key: "Email", Value: 1}},
				Options: options.Index().
					SetUnique(true).
					SetName(UserEmailIndex),
			},
		},
	)

	return err
}

func (r *MongoRepository) Create(
	ctx context.Context,
	args model.UserDetails,
) (*model.User, error) {
	if args.Email == "" {
		return nil, fmt.Errorf("email address is empty")
	}

	u := &model.User{
		ID:          uuid.New(),
		UserDetails: args,
	}
	_, err := r.coll.InsertOne(ctx, u)

	return u, err
}

func (r *MongoRepository) Get(
	ctx context.Context,
	id uuid.UUID,
) (*model.User, error) {
	var u model.User
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	return &u, err
}

func (r *MongoRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*model.User, error) {
	var u model.User
	err := r.coll.FindOne(ctx, bson.M{"Email": email}).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	return &u, nil
}
