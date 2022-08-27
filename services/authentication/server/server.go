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
	err := as.service.Register(
		ctx,
		request.Email,
		request.Password,
	)

	return &protos.RegisterResponse{}, err
}

func (as *AuthenticationServer) Authenticate(
	ctx context.Context,
	request *protos.AuthenticateRequest,
) (*protos.AuthenticateResponse, error) {
	refreshToken, accessToken, err := as.service.Authenticate(
		ctx,
		request.Email,
		request.Password,
	)

	return &protos.AuthenticateResponse{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, err
}

func (as *AuthenticationServer) Refresh(
	ctx context.Context,
	request *protos.RefreshRequest,
) (*protos.RefreshResponse, error) {
	accessToken, err := as.service.Refresh(
		ctx,
		request.RefreshToken,
	)

	return &protos.RefreshResponse{
		AccessToken: accessToken,
	}, err
}

func (as *AuthenticationServer) GetPublicKeys(
	ctx context.Context,
	request *protos.GetPublicKeysRequest,
) (*protos.GetPublicKeysResponse, error) {
	publicKeyEntities := make([]*protos.PublicKey, 0)
	publicKeys, err := as.service.GetPublicKeys(ctx)
	if err != nil {
		return &protos.GetPublicKeysResponse{
			PublicKeys: publicKeyEntities,
		}, err
	}

	for kid, keyValue := range publicKeys {
		publicKeyEntities = append(
			publicKeyEntities,
			&protos.PublicKey{
				Uuid:  kid.String(),
				Value: keyValue,
			},
		)
	}
	return &protos.GetPublicKeysResponse{
		PublicKeys: publicKeyEntities,
	}, nil
}
