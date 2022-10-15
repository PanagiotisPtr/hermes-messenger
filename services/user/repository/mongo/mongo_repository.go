package mongo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/entityutils"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/entityutils/filter"
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
	coll *mongo.Collection
	entityutils.RepoHelper
	logger *zap.Logger
}

func ProvideUserRepository(
	lc fx.Lifecycle,
	logger *zap.Logger,
	database *mongo.Database,
) repository.UserRepository {
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

func (r *MongoRepository) Find(
	ctx context.Context,
	f filter.Filter,
) ([]*model.User, error) {
	users := []*model.User{}
	cur, err := r.coll.Find(ctx, f.ToBSON())
	if err = cur.All(ctx, &users); err != nil {
		return users, err
	}

	return users, err
}

func (r *MongoRepository) FindOne(
	ctx context.Context,
	f filter.Filter,
) (*model.User, error) {
	u := model.User{}
	err := r.coll.FindOne(ctx, f.ToBSON()).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	return &u, err
}

func (r *MongoRepository) Create(
	ctx context.Context,
	args *model.User,
) (*model.User, error) {
	if args.Email == "" {
		return nil, fmt.Errorf("email address is empty")
	}

	u := args
	id := uuid.New()
	u.ID = &id
	u.Meta = entityutils.Meta{}
	r.RepoHelper.UpdateMeta(
		ctx,
		&u.Meta,
		entityutils.CreateOp,
	)

	_, err := r.coll.InsertOne(ctx, u)

	return u, err
}

func (r *MongoRepository) Update(
	ctx context.Context,
	f filter.Filter,
	args *model.User,
) ([]*model.User, error) {
	users := []*model.User{}
	b, err := bson.Marshal(args)
	if err != nil {
		return users, err
	}
	data := bson.M{}
	err = bson.Unmarshal(b, &data)
	if err != nil {
		return users, err
	}
	_, err = r.coll.UpdateMany(
		ctx,
		f.ToBSON(),
		bson.M{"$set": data},
	)
	if err != nil {
		return users, err
	}

	return r.Find(ctx, f)
}

func (r *MongoRepository) Delete(
	ctx context.Context,
	f filter.Filter,
) ([]*model.User, error) {
	users, err := r.Find(ctx, f)
	if err != nil {
		return users, err
	}
	_, err = r.coll.DeleteMany(ctx, f.ToBSON())

	return users, err
}
