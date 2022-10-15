package memory

import (
	"testing"

	"github.com/panagiotisptr/hermes-messenger/user/repository"
	"github.com/panagiotisptr/hermes-messenger/user/repository/testcases"
	"go.uber.org/zap"
)

func getRepository() repository.UserRepository {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	return ProvideUserRepository(logger)
}

func TestCreate(t *testing.T) {
	repo := getRepository()
	testcases := testcases.CreateTestcases(repo)

	for _, tc := range testcases {
		t.Run(tc.Name, tc.RunFunc())
	}
}
