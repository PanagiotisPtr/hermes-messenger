package server

import (
	"context"

	"github.com/panagiotisptr/hermes-messenger/protos"
	"github.com/panagiotisptr/hermes-messenger/services/authentication/server/authentication"
	"go.uber.org/zap"
)

type AuthenticationServer struct {
	logger  *zap.Logger
	service *authentication.Service
	protos.UnimplementedAuthenticationServer
}

func ProvideAuthenticationServer(logger *zap.Logger, service *authentication.Service) (*AuthenticationServer, error) {
	return &AuthenticationServer{
		logger:  logger,
		service: service,
	}, nil
}

func (as *AuthenticationServer) Register(
	ctx context.Context,
	request *protos.RegisterRequest,
) (*protos.RegisterResponse, error) {
	success, err := as.service.Register(request.Email, request.Password)

	return &protos.RegisterResponse{
		Success: success,
	}, err
}

func (as *AuthenticationServer) Authenticate(
	ctx context.Context,
	request *protos.AuthenticateRequest,
) (*protos.AuthenticateResponse, error) {
	refreshToken, accessToken, err := as.service.Authenticate(request.Email, request.Password)

	return &protos.AuthenticateResponse{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, err
}

func (as *AuthenticationServer) Refresh(
	ctx context.Context,
	request *protos.RefreshRequest,
) (*protos.RefreshResponse, error) {
	accessToken, err := as.service.Refresh(request.RefreshToken)

	return &protos.RefreshResponse{
		AccessToken: accessToken,
	}, err
}

func (as *AuthenticationServer) GetPublicKey(
	ctx context.Context,
	request *protos.GetPublicKeyRequest,
) (*protos.GetPublicKeyResponse, error) {
	publicKey, err := as.service.GetPublicKey()

	return &protos.GetPublicKeyResponse{
		PublicKey: publicKey,
	}, err
}
