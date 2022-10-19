package server

import (
	"context"
	"fmt"
	"time"

	"github.com/panagiotisptr/hermes-messenger/protos"
	"github.com/panagiotisptr/hermes-messenger/user/model"
	"github.com/panagiotisptr/hermes-messenger/user/service"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	"github.com/google/uuid"
)

// UserServer represents a server that can be used for the grpc user service
type UserServer struct {
	logger  *zap.Logger
	service *service.UserService
	protos.UnimplementedUserServiceServer
}

// ProvideUserServer provides a user service server
func ProvideUserServer(
	logger *zap.Logger,
	service *service.UserService,
) (*UserServer, error) {
	return &UserServer{
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

func metadataToContext(ctx context.Context, keys ...string) (context.Context, error) {
	m, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, fmt.Errorf("missing grpc metadata in request")
	}
	for _, k := range keys {
		values := m.Get(k)
		if len(values) == 0 {
			return ctx, fmt.Errorf("missing key '%s' in request metadata", k)
		}
		ctx = context.WithValue(ctx, k, values[0])
	}

	return ctx, nil
}

// CreateUser creates a new user
func (us *UserServer) CreateUser(
	ctx context.Context,
	request *protos.CreateUserRequest,
) (*protos.CreateUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	ctx, err := metadataToContext(ctx, "request-id")
	if err != nil {
		return nil, err
	}
	u, err := us.service.CreateUser(
		ctx,
		&model.User{
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
func (us *UserServer) GetUser(
	ctx context.Context,
	request *protos.GetUserRequest,
) (*protos.GetUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, err
	}
	ctx, err = metadataToContext(ctx, "user-id", "request-id")
	if err != nil {
		return nil, err
	}

	response := &protos.GetUserResponse{
		User: nil,
	}
	u, err := us.service.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	response.User = userToEntity(u)

	return response, nil
}

// GetUser finds a user by their email - returns nil if not found
func (us *UserServer) GetUserByEmail(
	ctx context.Context,
	request *protos.GetUserByEmailRequest,
) (*protos.GetUserByEmailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	ctx, err := metadataToContext(ctx, "user-id", "request-id")
	if err != nil {
		return nil, err
	}

	response := &protos.GetUserByEmailResponse{
		User: nil,
	}
	u, err := us.service.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	response.User = userToEntity(u)

	return response, nil
}
