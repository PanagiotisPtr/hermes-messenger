package mongo

import (
	"context"
	"testing"

	"github.com/panagiotisptr/hermes-messenger/libs/utils/mongoutils"
	"github.com/panagiotisptr/hermes-messenger/user/config"
	"github.com/panagiotisptr/hermes-messenger/user/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func getDatabase(
	t *testing.T,
	logger *zap.Logger,
	lifecycle fx.Lifecycle,
) *mongo.Database {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	config, err := config.ProvideTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	mongoClient, err := mongoutils.ProvideMongoClient(
		lifecycle,
		logger,
		&config.Mongo,
	)
	if err != nil {
		t.Fatal(err)
	}

	mongoDb := mongoutils.ProvideMongoDatabase(
		mongoClient,
		&config.Mongo,
	)

	return mongoDb
}

func TestCreate(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	lifecycle := fxtest.NewLifecycle(t)
	mongoDb := getDatabase(t, logger, lifecycle)
	defer mongoDb.Drop(context.Background())

	repo := ProvideUserRepository(
		lifecycle,
		logger,
		mongoDb,
	)

	user := model.UserDetails{}
	u, err := repo.Create(
		context.Background(),
		user,
	)
}
