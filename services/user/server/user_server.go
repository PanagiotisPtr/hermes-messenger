package server

import (
	"context"

	"github.com/panagiotisptr/hermes-messenger/protos"
	"github.com/panagiotisptr/hermes-messenger/user/model"
	"github.com/panagiotisptr/hermes-messenger/user/service"
	"go.uber.org/zap"

	"github.com/google/uuid"
)

// Server represents a server that can be used for the grpc user service
type Server struct {
	logger  *zap.Logger
	service *service.Service
	protos.UnimplementedUserServiceServer
}

// ProvideUserServer provides a user service server
func ProvideUserServer(
	logger *zap.Logger,
	service *service.Service,
) (*Server, error) {
	return &Server{
		logger:  logger,
		service: service,
	}, nil
}

// userToEntity converts an entity to a user object for grpc
func userToEntity(u *model.User) *protos.User {
	if u == nil {
		return nil
	}

	return &protos.User{
		Id:        u.ID.String(),
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

// CreateUser creates a new user
func (us *Server) CreateUser(
	ctx context.Context,
	request *protos.CreateUserRequest,
) (*protos.CreateUserResponse, error) {
	u, err := us.service.CreateUser(
		ctx,
		model.UserDetails{
			Email:     request.Email,
			FirstName: request.FirstName,
			LastName:  request.LastName,
		},
	)

	return &protos.CreateUserResponse{
		User: userToEntity(u),
	}, err
}

// GetUser finds a user by their (uu)id - returns nil if not found
func (us *Server) GetUser(
	ctx context.Context,
	request *protos.GetUserRequest,
) (*protos.GetUserResponse, error) {
	response := &protos.GetUserResponse{
		User: nil,
	}
	id, err := uuid.Parse(request.Id)
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

// GetUser finds a user by their email - returns nil if not found
func (us *Server) GetUserByEmail(
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
