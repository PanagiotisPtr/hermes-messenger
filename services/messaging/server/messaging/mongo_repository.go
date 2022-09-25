package messaging

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	MessagingCollectionName = "messaging"
	MessagingTimestampIndex = "messaging_timestamp_index"
	MessagingFromToIndex    = "messaging_from_to_index"
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
		coll:   database.Collection(MessagingCollectionName),
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
					SetName(MessagingFromToIndex),
			},
			{
				Keys: bson.D{
					{Key: "Timestamp", Value: 1},
				},
				Options: options.Index().
					SetName(MessagingTimestampIndex),
			},
		},
	)

	return err
}

func (r *MongoRepository) SaveMessage(
	ctx context.Context,
	from uuid.UUID,
	to uuid.UUID,
	content string,
) error {
	m := &Message{
		ID:        uuid.New(),
		From:      from,
		To:        to,
		Timestamp: time.Now().Unix(),
		Content:   content,
	}
	_, err := r.coll.InsertOne(ctx, m)

	return err
}

func (r *MongoRepository) GetMessages(
	ctx context.Context,
	from uuid.UUID,
	to uuid.UUID,
	size int64,
	offset int64,
) ([]*Message, error) {
	var ms []*Message
	opts := options.
		Find().
		SetSort(bson.D{{Key: "Timestamp", Value: -1}}).
		SetSkip(offset).
		SetLimit(size)
	cursor, err := r.coll.Find(
		ctx,
		bson.D{{Key: "$or", Value: []interface{}{
			bson.D{{Key: "$and", Value: []interface{}{
				bson.D{{Key: "From", Value: from}},
				bson.D{{Key: "To", Value: to}},
			}}},
			bson.D{{Key: "$and", Value: []interface{}{
				bson.D{{Key: "From", Value: to}},
				bson.D{{Key: "To", Value: from}},
			}}},
		}}},
		opts,
	)
	if err == mongo.ErrNoDocuments {
		return ms, nil
	}
	err = cursor.All(ctx, &ms)

	return ms, err
}
