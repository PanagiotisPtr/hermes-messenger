package server

import (
	"context"
	"log"
	"panagiotisptr/authentication/protos"
	"panagiotisptr/authentication/server/authentication"
)

type AuthenticationServer struct {
	logger  *log.Logger
	service *authentication.Service
	protos.UnimplementedAuthenticationServer
}

func NewAuthenticationServer(logger *log.Logger) (*AuthenticationServer, error) {
	service, err := authentication.NewService(logger)
	if err != nil {
		return nil, err
	}

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
