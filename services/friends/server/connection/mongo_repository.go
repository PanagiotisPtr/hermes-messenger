package connection

import (
	"context"

	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/friends/server/connection/status"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	FriendsCollectionName = "friends"
	FriendsFromToIndex    = "friends_from_to_index"
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
		coll:   database.Collection(FriendsCollectionName),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Sugar().Info("Initialising mongobd indexes for friends repository")
			return repo.InitIndexes(ctx)
		},
	})

	return repo
}

func (r *MongoRepository) InitIndexes(
	ctx context.Context,
) error {
	_, err := r.coll.Indexes().CreateMany(ctx,
		[]mongo.IndexModel{
			{
				Keys: bson.D{
					{Key: "From", Value: 1},
					{Key: "To", Value: 1},
				},
				Options: options.Index().
					SetUnique(true).
					SetName(FriendsFromToIndex),
			},
		},
	)

	return err
}

func (r *MongoRepository) AddConnection(
	ctx context.Context,
	from uuid.UUID,
	to uuid.UUID,
) error {
	c, err := r.getConnection(ctx, to, from)
	if err != nil {
		return err
	}
	if c == nil {
		c = &Connection{
			ID:     uuid.New(),
			From:   from,
			To:     to,
			Status: status.Pending,
		}
		_, err = r.coll.InsertOne(ctx, c)
	} else {
		_, err = r.coll.UpdateOne(
			ctx,
			bson.M{"_id": c.ID},
			bson.D{{
				Key: "$set",
				Value: bson.D{
					{Key: "Status", Value: status.Accepted},
				},
			}},
		)
	}

	return err
}

func (r *MongoRepository) getConnection(
	ctx context.Context,
	from uuid.UUID,
	to uuid.UUID,
) (*Connection, error) {
	var c Connection
	err := r.coll.FindOne(
		ctx,
		bson.M{"From": from, "To": to},
	).Decode(&c)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	return &c, err
}

func (r *MongoRepository) RemoveConnection(
	ctx context.Context,
	from uuid.UUID,
	to uuid.UUID,
) error {
	_, err := r.coll.DeleteOne(
		ctx,
		bson.M{"From": from, "To": to},
	)

	return err
}

func (r *MongoRepository) GetConnections(
	ctx context.Context,
	userUuid uuid.UUID,
) ([]*Connection, error) {
	var cs []*Connection
	cursor, err := r.coll.Find(
		ctx,
		bson.D{{Key: "$or", Value: []interface{}{
			bson.D{{Key: "From", Value: userUuid}},
			bson.D{{Key: "To", Value: userUuid}},
		}}},
	)
	if err == mongo.ErrNoDocuments {
		return cs, nil
	}
	err = cursor.All(ctx, &cs)

	return cs, err
}
