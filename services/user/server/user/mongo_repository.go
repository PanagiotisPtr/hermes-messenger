package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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

func ProvideMongoRepository(
	lc fx.Lifecycle,
	logger *zap.Logger,
	database *mongo.Database,
) Repository {
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
	args UserDetails,
) (*User, error) {
	if args.Email == "" {
		return nil, fmt.Errorf("email address is empty")
	}

	u := &User{
		ID:          uuid.New(),
		UserDetails: args,
	}
	_, err := r.coll.InsertOne(ctx, u)

	return u, err
}

func (r *MongoRepository) Get(
	ctx context.Context,
	id uuid.UUID,
) (*User, error) {
	var u User
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	return &u, err
}

func (r *MongoRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*User, error) {
	var u User
	err := r.coll.FindOne(ctx, bson.M{"Email": email}).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	return &u, nil
}
