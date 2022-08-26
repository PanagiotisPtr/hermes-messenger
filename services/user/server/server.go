package server

import (
	"context"

	"github.com/panagiotisptr/hermes-messenger/protos"
	"github.com/panagiotisptr/hermes-messenger/user/server/user"
	"go.uber.org/zap"

	"github.com/google/uuid"
)

type UserServer struct {
	logger  *zap.Logger
	service *user.Service
	protos.UnimplementedUserServer
}

func ProvideUserServer(
	logger *zap.Logger,
	service *user.Service,
) (*UserServer, error) {
	return &UserServer{
		logger:  logger,
		service: service,
	}, nil
}

func userToEntity(u *user.User) *protos.UserEntity {
	if u == nil {
		return nil
	}

	return &protos.UserEntity{
		Uuid:  u.Uuid.String(),
		Email: u.Email,
	}
}

func (us *UserServer) RegisterUser(
	ctx context.Context,
	request *protos.RegisterUserRequest,
) (*protos.RegisterUserResponse, error) {
	u, err := us.service.RegisterUser(ctx, request.Email)

	return &protos.RegisterUserResponse{
		User: userToEntity(u),
	}, err
}

func (us *UserServer) GetUser(
	ctx context.Context,
	request *protos.GetUserRequest,
) (*protos.GetUserResponse, error) {
	response := &protos.GetUserResponse{
		User: nil,
	}
	id, err := uuid.Parse(request.Uuid)
	if err != nil {
		return response, err
	}

	u, err := us.service.GetUser(ctx, id)
	if err != nil {
		return response, err
	}
	response.User = userToEntity(u)

	return response, nil
}

func (us *UserServer) GetUserByEmail(
	ctx context.Context,
	request *protos.GetUserByEmailRequest,
) (*protos.GetUserByEmailResponse, error) {
	response := &protos.GetUserByEmailResponse{
		User: nil,
	}

	u, err := us.service.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return response, err
	}
	response.User = userToEntity(u)

	return response, nil
}
