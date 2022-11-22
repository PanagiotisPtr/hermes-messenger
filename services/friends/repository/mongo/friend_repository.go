package mongo

import (
	"context"

	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/friends/model"
	"github.com/panagiotisptr/hermes-messenger/friends/repository"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/entityutils"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/entityutils/filter"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/loggingutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	FriendCollectionName = "friends"
	UserIndex            = "user_index"
	FriendIndex          = "friend_index"
)

type MongoRepository struct {
	coll   *mongo.Collection
	logger *zap.Logger
}

func ProvideFriendRepository(
	lc fx.Lifecycle,
	logger *zap.Logger,
	database *mongo.Database,
) repository.FriendRepository {
	repo := &MongoRepository{
		coll: database.Collection(FriendCollectionName),
		logger: logger.With(
			zap.String("repository", "FriendRepository"),
			zap.String("type", "mongo"),
		),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			repo.logger.Sugar().Info("initialising mongobd indexes for friend repository")
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
				Keys: bson.D{{Key: "userId", Value: 1}},
				Options: options.Index().
					SetName(UserIndex),
			},
			{
				Keys: bson.D{{Key: "friendId", Value: 1}},
				Options: options.Index().
					SetName(FriendIndex),
			},
		},
	)

	return err
}

func (r *MongoRepository) Find(
	ctx context.Context,
	f filter.Filter,
) (<-chan *model.Friend, error) {
	l := loggingutils.
		LoggerWithRequestID(ctx, r.logger).
		With(
			zap.String("method", "Find"),
		).Sugar()
	ch := make(chan *model.Friend)
	go func() {
		close(ch)
		cur, err := r.coll.Find(ctx, f.ToBSON())
		for cur.Next(ctx) {
			var u *model.Friend
			if err = cur.Decode(&u); err != nil {
				l.Error("decoding bson to friend", err)
			}
			ch <- u
		}
	}()

	return ch, nil
}

func (r *MongoRepository) FindOne(
	ctx context.Context,
	f filter.Filter,
) (*model.Friend, error) {
	l := loggingutils.
		LoggerWithRequestID(ctx, r.logger).
		With(
			zap.String("method", "FindOne"),
		).Sugar()
	u := model.Friend{}
	err := r.coll.FindOne(ctx, f.ToBSON()).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		l.Error("finding friend", err)
	}

	return &u, err
}

func (r *MongoRepository) Create(
	ctx context.Context,
	args *model.Friend,
) (*model.Friend, error) {
	l := loggingutils.
		LoggerWithRequestID(ctx, r.logger).
		With(
			zap.String("method", "Create"),
		).Sugar()

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
		l.Error("creating friend", err)
	}

	return u, err
}

func (r *MongoRepository) Update(
	ctx context.Context,
	f filter.Filter,
	args *model.Friend,
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
		l.Error("updating friends", err)
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
		l.Error("deleting friends", err)
	}

	return res.DeletedCount, err
}
