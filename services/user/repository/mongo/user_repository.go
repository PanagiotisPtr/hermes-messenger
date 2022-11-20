package mongo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/entityutils"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/entityutils/filter"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/loggingutils"
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
	UserEmailIndex     = "users_email_index"
)

type MongoRepository struct {
	coll   *mongo.Collection
	logger *zap.Logger
}

func ProvideUserRepository(
	lc fx.Lifecycle,
	logger *zap.Logger,
	database *mongo.Database,
) repository.UserRepository {
	repo := &MongoRepository{
		coll: database.Collection(UserCollectionName),
		logger: logger.With(
			zap.String("repository", "UserRepository"),
			zap.String("type", "mongo"),
		),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			repo.logger.Sugar().Info("initialising mongobd indexes for user repository")
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
				Keys: bson.D{{Key: "email", Value: 1}},
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
) (<-chan *model.User, error) {
	l := loggingutils.
		LoggerWithRequestID(ctx, r.logger).
		With(
			zap.String("method", "Find"),
		).Sugar()
	ch := make(chan *model.User)
	go func() {
		close(ch)
		cur, err := r.coll.Find(ctx, f.ToBSON())
		for cur.Next(ctx) {
			var u *model.User
			if err = cur.Decode(&u); err != nil {
				l.Error("decoding bson to user", err)
			}
			ch <- u
		}
	}()

	return ch, nil
}

func (r *MongoRepository) FindOne(
	ctx context.Context,
	f filter.Filter,
) (*model.User, error) {
	l := loggingutils.
		LoggerWithRequestID(ctx, r.logger).
		With(
			zap.String("method", "FindOne"),
		).Sugar()
	u := model.User{}
	err := r.coll.FindOne(ctx, f.ToBSON()).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		l.Error("finding user", err)
	}

	return &u, err
}

func (r *MongoRepository) Create(
	ctx context.Context,
	args *model.User,
) (*model.User, error) {
	l := loggingutils.
		LoggerWithRequestID(ctx, r.logger).
		With(
			zap.String("method", "Create"),
		).Sugar()
	if args.Email == "" {
		return nil, fmt.Errorf("email address is empty")
	}

	u := args
	id := uuid.New()
	u.ID = &id
	u.Meta = entityutils.Meta{}
	u.UpdateMeta(
		ctx,
		entityutils.CreateOp,
	)

	_, err := r.coll.InsertOne(ctx, u)
	if err != nil {
		l.Error("creating user", err)
	}

	return u, err
}

func (r *MongoRepository) Update(
	ctx context.Context,
	f filter.Filter,
	args *model.User,
) (int64, error) {
	l := loggingutils.
		LoggerWithRequestID(ctx, r.logger).
		With(
			zap.String("method", "Update"),
		).Sugar()
	b, err := bson.Marshal(args)
	if err != nil {
		l.Error("unmarshalling args", err)
		return 0, err
	}
	data := bson.M{}
	err = bson.Unmarshal(b, &data)
	if err != nil {
		l.Error("marshalling data", err)
		return 0, err
	}
	res, err := r.coll.UpdateMany(
		ctx,
		f.ToBSON(),
		bson.M{"$set": data},
	)
	if err != nil {
		l.Error("updating users", err)
		return 0, err
	}

	return res.ModifiedCount, nil
}

func (r *MongoRepository) Delete(
	ctx context.Context,
	f filter.Filter,
) (int64, error) {
	l := loggingutils.
		LoggerWithRequestID(ctx, r.logger).
		With(
			zap.String("method", "Delete"),
		).Sugar()
	res, err := r.coll.DeleteMany(ctx, f.ToBSON())
	if err != nil {
		l.Error("deleting users", err)
	}

	return res.DeletedCount, err
}
