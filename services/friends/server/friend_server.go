package server

import (
	friendApp "github.com/panagiotisptr/hermes-messenger/friends/app"
	"github.com/panagiotisptr/hermes-messenger/friends/model"
	"github.com/panagiotisptr/hermes-messenger/protos"
	"go.uber.org/zap"
)

type FriendServer struct {
	logger *zap.Logger
	app    *friendApp.App
	protos.UnimplementedFriendServiceServer
}

func ProvideFriendsServer(
	logger *zap.Logger,
	app *friendApp.App,
) (*FriendServer, error) {
	return &FriendServer{
		logger: logger.With(
			zap.String("server", "FriendServer"),
		),
		app: app,
	}, nil
}

func friendToEntity(f *model.Friend) *protos.Friend {
	if f == nil {
		return nil
	}

	return &protos.Friend{
		Id:       f.ID.String(),
		UserId:   f.UserId.String(),
		FriendId: f.FriendId.String(),
		Status:   protos.Status(f.Status),
	}
}
