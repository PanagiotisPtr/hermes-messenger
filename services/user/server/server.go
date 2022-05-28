package server

import (
	"context"
	"log"

	"github.com/panagiotisptr/hermes-messenger/user/protos"
	"github.com/panagiotisptr/hermes-messenger/user/server/user"

	"github.com/google/uuid"
)

type UserServer struct {
	logger  *log.Logger
	service *user.Service
	protos.UnimplementedUserServer
}

func NewUserService(logger *log.Logger) (*UserServer, error) {
	service := user.NewService(logger)

	return &UserServer{
		logger:  logger,
		service: service,
	}, nil
}

func (us *UserServer) RegisterUser(
	ctx context.Context,
	request *protos.RegisterUserRequest,
) (*protos.RegisterUserResponse, error) {
	response := &protos.RegisterUserResponse{
		Success: false,
	}
	userUuid, err := uuid.Parse(request.Uuid)
	if err != nil {
		return response, err
	}
	err = us.service.RegisterUser(userUuid, request.Email)
	response.Success = err == nil

	return response, err
}

func (us *UserServer) GetUser(
	ctx context.Context,
	request *protos.GetUserRequest,
) (*protos.GetUserResponse, error) {
	response := &protos.GetUserResponse{
		User: &protos.UserEntity{
			Uuid:  "",
			Email: "",
		},
	}
	userUuid, err := uuid.Parse(request.Uuid)
	if err != nil {
		return response, err
	}

	user, err := us.service.GetUser(userUuid)
	if err != nil {
		return response, err
	}
	response.User.Uuid = user.Uuid.String()
	response.User.Email = user.Email

	return response, err
}
