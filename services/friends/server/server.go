package server

import (
	"context"
	"log"

	"github.com/panagiotisptr/hermes-messenger/friends/server/friends"
	"github.com/panagiotisptr/hermes-messenger/protos"

	"github.com/google/uuid"
)

type FriendsServer struct {
	logger  *log.Logger
	service *friends.Service
	protos.UnimplementedFriendsServer
}

func NewFriendsServer(logger *log.Logger) (*FriendsServer, error) {
	service := friends.NewService(logger)

	return &FriendsServer{
		logger:  logger,
		service: service,
	}, nil
}

func (as *FriendsServer) AddFriend(
	ctx context.Context,
	request *protos.AddFriendRequest,
) (*protos.AddFriendResponse, error) {
	response := &protos.AddFriendResponse{
		Sent: false,
	}
	userUuid, err := uuid.Parse(request.UserUuid)
	if err != nil {
		return response, err
	}

	friendUuid, err := uuid.Parse(request.FriendUuid)
	if err != nil {
		return response, err
	}

	err = as.service.AddFriend(userUuid, friendUuid)
	if err != nil {
		return response, err
	}
	response.Sent = true

	return response, nil
}

func (as *FriendsServer) RemoveFriend(
	ctx context.Context,
	request *protos.RemoveFriendRequest,
) (*protos.RemoveFriendResponse, error) {
	response := &protos.RemoveFriendResponse{
		Success: false,
	}
	userUuid, err := uuid.Parse(request.UserUuid)
	if err != nil {
		return response, err
	}

	friendUuid, err := uuid.Parse(request.FriendUuid)
	if err != nil {
		return response, err
	}

	err = as.service.RemoveFriend(userUuid, friendUuid)
	if err != nil {
		return response, err
	}
	response.Success = true

	return response, nil
}

func (as *FriendsServer) GetFriends(
	ctx context.Context,
	request *protos.GetFriendsRequest,
) (*protos.GetFriendsResponse, error) {
	response := &protos.GetFriendsResponse{
		Friends: make([]*protos.Friend, 0),
	}
	userUuid, err := uuid.Parse(request.UserUuid)
	if err != nil {
		return response, err
	}

	friends, err := as.service.GetFriends(userUuid)
	if err != nil {
		return response, err
	}

	for _, friend := range friends {
		response.Friends = append(response.Friends, &protos.Friend{
			UserUuid: friend.UserUuid.String(),
			Status:   friend.Status,
		})
	}

	return response, nil
}
