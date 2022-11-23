package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/friends/model"
	"github.com/panagiotisptr/hermes-messenger/friends/service"
	"github.com/panagiotisptr/hermes-messenger/protos"
	"go.uber.org/zap"
)

type App struct {
	logger        *zap.Logger
	friendService *service.FriendService
	userClient    protos.UserServiceClient
}

func ProvideApp(
	logger *zap.Logger,
	friendService *service.FriendService,
	userClient protos.UserServiceClient,
) *App {
	return &App{
		logger: logger.With(
			zap.String("app", "FriendApp"),
		),
		friendService: friendService,
		userClient:    userClient,
	}
}

// Returns error if the user with the specified uuid doesn't exist
func (a *App) userExists(ctx context.Context, id uuid.UUID) error {
	userResp, err := a.userClient.GetUser(ctx, &protos.GetUserRequest{
		Id: id.String(),
	})
	if err != nil {
		return err
	}
	if userResp.User == nil {
		return fmt.Errorf("could not find user with UUID: %s", id.String())
	}

	return nil
}

func (a *App) FriendCreate(
	ctx context.Context,
	friendId uuid.UUID,
) (*model.Friend, error) {
	err := a.userExists(ctx, friendId)
	if err != nil {
		return nil, err
	}

	return a.friendService.Create(ctx, &model.Friend{
		FriendId: friendId,
	})
}

func (a *App) FriendDelete(
	ctx context.Context,
	friendId uuid.UUID,
) error {
	err := a.userExists(ctx, friendId)
	if err != nil {
		return err
	}

	return a.friendService.Delete(ctx, &model.Friend{
		FriendId: friendId,
	})
}

func (a *App) FriendFind(
	ctx context.Context,
) (<-chan *model.Friend, error) {
	return a.friendService.Find(ctx)
}
