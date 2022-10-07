package mongo

import (
	"testing"

	"github.com/panagiotisptr/hermes-messenger/libs/utils/mongoutils"
	"github.com/panagiotisptr/hermes-messenger/user/config"
	"github.com/panagiotisptr/hermes-messenger/user/repository"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func getRepository(t *testing.T) repository.Repository {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	config := config.ProvideConfig()
	lifecycle := fxtest.NewLifecycle(t)

	mongoClient, err := mongoutils.ProvideMongoClient(
		lifecycle,
		logger,
		&config.MongoConfig,
	)
	if err != nil {
		t.Fatal(err)
	}

	mongoDb := mongoutils.ProvideMongoDatabase(
		mongoClient,
		&config.MongoConfig,
	)

	return ProvideUserRepository(
		lifecycle,
		logger,
		mongoDb,
	)
}
