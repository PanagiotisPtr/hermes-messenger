package server

import (
	"context"

	friendApp "github.com/panagiotisptr/hermes-messenger/friends/app"
	"github.com/panagiotisptr/hermes-messenger/friends/model"
	"github.com/panagiotisptr/hermes-messenger/protos"
	"go.uber.org/zap"

	"github.com/google/uuid"
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

func (fs *FriendServer) AddFriend(
	ctx context.Context,
	request *protos.AddFriendRequest,
) (*protos.AddFriendResponse, error) {
	response := &protos.AddFriendResponse{}
	fromUuid, err := uuid.Parse(request.UserUuid)
	if err != nil {
		return response, err
	}
	toUuid, err := uuid.Parse(request.FriendUuid)
	if err != nil {
		return response, err
	}
	err = fs.service.AddFriend(ctx, fromUuid, toUuid)

	return response, err
}

func (fs *FriendServer) GetFriends(
	ctx context.Context,
	request *protos.GetFriendsRequest,
) (*protos.GetFriendsResponse, error) {
	response := &protos.GetFriendsResponse{
		Friends: make([]*protos.Friend, 0),
	}
	id, err := uuid.Parse(request.UserUuid)
	if err != nil {
		return response, err
	}

	frs, err := fs.service.GetFriends(ctx, id)
	if err != nil {
		return response, err
	}
	for _, f := range frs {
		response.Friends = append(response.Friends, friendToEntity(f))
	}

	return response, err
}

func (fs *FriendServer) RemoveFriend(
	ctx context.Context,
	request *protos.RemoveFriendRequest,
) (*protos.RemoveFriendResponse, error) {
	response := &protos.RemoveFriendResponse{}
	fromUuid, err := uuid.Parse(request.UserUuid)
	if err != nil {
		return response, err
	}
	toUuid, err := uuid.Parse(request.FriendUuid)
	if err != nil {
		return response, err
	}

	err = fs.service.RemoveFriend(ctx, fromUuid, toUuid)

	return response, err
}
