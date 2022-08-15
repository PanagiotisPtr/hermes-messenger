package server

import (
	"context"

	"github.com/panagiotisptr/hermes-messenger/friends/server/friends"
	"github.com/panagiotisptr/hermes-messenger/protos"
	"go.uber.org/zap"

	"github.com/google/uuid"
)

type FriendsServer struct {
	logger  *zap.Logger
	service *friends.Service
	protos.UnimplementedFriendsServer
}

func ProvideFriendsServer(
	logger *zap.Logger,
	service *friends.Service,
) (*FriendsServer, error) {
	return &FriendsServer{
		logger:  logger,
		service: service,
	}, nil
}

func friendToEntity(f *friends.Friend) *protos.Friend {
	if f == nil {
		return nil
	}

	return &protos.Friend{
		FriendUuid: f.FriendUuid.String(),
		Status:     f.Status,
	}
}

func (fs *FriendsServer) AddFriend(
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

func (fs *FriendsServer) GetFriends(
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

func (fs *FriendsServer) RemoveFriend(
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
