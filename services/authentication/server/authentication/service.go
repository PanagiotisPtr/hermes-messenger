package authentication

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/panagiotisptr/hermes-messenger/protos"
	"github.com/panagiotisptr/hermes-messenger/services/authentication/server/config"
	"github.com/panagiotisptr/hermes-messenger/services/authentication/server/credentials"
	"github.com/panagiotisptr/hermes-messenger/services/authentication/server/keys"
	"github.com/panagiotisptr/hermes-messenger/services/authentication/server/token"
	"go.uber.org/zap"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	selfUuid              uuid.UUID
	logger                *zap.Logger
	refreshKid            uuid.UUID
	accessKid             uuid.UUID
	tokenRepository       token.Repository
	keysRepository        keys.Repository
	credentialsRepository credentials.Repository
	refreshTokenDuration  time.Duration
	accessTokenDuration   time.Duration
	userClient            protos.UserClient
}

func ProvideAuthenticationService(
	cfg *config.Config,
	logger *zap.Logger,
	tokenRepository token.Repository,
	keysRepository keys.Repository,
	credentialsRepository credentials.Repository,
	userClient protos.UserClient,
) (*Service, error) {
	service := &Service{
		selfUuid:              cfg.UUID,
		logger:                logger,
		tokenRepository:       tokenRepository,
		keysRepository:        keysRepository,
		credentialsRepository: credentialsRepository,
		refreshTokenDuration:  cfg.RefreshTokenDuration,
		accessTokenDuration:   cfg.AccessTokenDuration,
		userClient:            userClient,
	}

	return service, nil
}

func (s *Service) GenerateKeyPair(
	ctx context.Context,
	keyPairGenerationInterval time.Duration,
) {
	refreshTokenKeyPair, err := keys.GenerateRSAKeyPair(keys.RefreshKeyType)
	if err != nil {
		s.logger.Sugar().Error(err)
		return
	}

	accessTokenKeyPair, err := keys.GenerateRSAKeyPair(keys.AccessKeyType)
	if err != nil {
		s.logger.Sugar().Error(err)
		return
	}

	err = s.keysRepository.StoreKeyPair(ctx, refreshTokenKeyPair, 2*keyPairGenerationInterval)
	if err != nil {
		s.logger.Sugar().Error(err)
		return
	}

	err = s.keysRepository.StoreKeyPair(ctx, accessTokenKeyPair, 2*keyPairGenerationInterval)
	if err != nil {
		s.logger.Sugar().Error(err)
		return
	}

	s.refreshKid = refreshTokenKeyPair.Uuid
	s.accessKid = accessTokenKeyPair.Uuid

	s.logger.Sugar().Infof(
		"Updated access and refresh keys for authentication service with UUID: " +
			s.selfUuid.String() +
			". Using refresh token key with UUID:" +
			refreshTokenKeyPair.Uuid.String() +
			". Using access token key with UUID:" +
			accessTokenKeyPair.Uuid.String(),
	)
}

func (s *Service) Register(
	ctx context.Context,
	email string,
	password string,
) error {
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	resp, err := s.userClient.RegisterUser(
		ctx,
		&protos.RegisterUserRequest{
			Email: email,
		},
	)
	if err != nil {
		return err
	}
	uid, err := uuid.Parse(resp.User.Uuid)
	if err != nil {
		return err
	}

	err = s.credentialsRepository.CreateCredentials(
		ctx,
		uid,
		string(passwordHash),
	)
	if err != nil {
		return err
	}

	return err
}

func (s *Service) generateToken(
	ctx context.Context,
	keyType string,
	data string,
	kid uuid.UUID,
	ttl time.Duration,
) (string, error) {
	tokenPrivateKey, err := s.keysRepository.GetPrivateKey(
		ctx,
		kid,
		keyType,
	)
	if err != nil {
		s.logger.Sugar().Errorf("[ERROR] Failed to get private key: %v", err)
		return "", fmt.Errorf("Error when generating token")
	}
	token, err := s.tokenRepository.SignTokenWithClaims(
		data,
		kid,
		ttl,
		tokenPrivateKey,
	)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) Authenticate(
	ctx context.Context,
	email string,
	password string,
) (refreshToken string, accessToken string, err error) {
	resp, err := s.userClient.RegisterUser(
		ctx,
		&protos.RegisterUserRequest{
			Email: email,
		},
	)
	if err != nil {
		return refreshToken, accessToken, err
	}
	uid, err := uuid.Parse(resp.User.Uuid)
	if err != nil {
		return refreshToken, accessToken, err
	}

	credentials, err := s.credentialsRepository.GetCredentials(
		ctx,
		uid,
	)
	if err != nil {
		return refreshToken, accessToken, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(credentials.Password), []byte(password))
	if err != nil {
		return refreshToken, accessToken, err
	}

	refreshToken, err = s.generateToken(
		ctx,
		keys.RefreshKeyType,
		uid.String(),
		s.refreshKid,
		s.refreshTokenDuration,
	)
	if err != nil {
		return refreshToken, accessToken, err
	}

	accessToken, err = s.generateToken(
		ctx,
		keys.AccessKeyType,
		uid.String(),
		s.accessKid,
		s.accessTokenDuration,
	)
	if err != nil {
		return refreshToken, accessToken, err
	}

	return refreshToken, accessToken, nil
}

func (s *Service) Refresh(
	ctx context.Context,
	refreshToken string,
) (string, error) {
	publicKey, err := s.keysRepository.GetPublicKey(
		ctx,
		s.refreshKid,
		keys.RefreshKeyType,
	)
	if err != nil {
		s.logger.Sugar().Errorf("[ERROR] Failed to get public keys: %v", err)
		return "", fmt.Errorf("Could not validate refresh token")
	}

	data, err := s.tokenRepository.ValidateToken(refreshToken, publicKey)
	if err != nil {
		return "", fmt.Errorf("Could not validate refresh token")
	}
	accessToken, err := s.generateToken(
		ctx,
		keys.AccessKeyType,
		data,
		s.accessKid,
		s.accessTokenDuration,
	)
	if err != nil {
		return "", err
	}

	return accessToken, nil

}

func (s *Service) GetPublicKeys(ctx context.Context) (map[uuid.UUID]string, error) {
	publicKeyStrings := make(map[uuid.UUID]string, 0)
	// Only allow access to the public key for the access tokens
	publicKeys, err := s.keysRepository.GetAllPublicKeys(
		ctx,
		keys.AccessKeyType,
	)
	if err != nil {
		s.logger.Sugar().Warnf("Failed to fetch public keys. Error: %v", err)
		return publicKeyStrings, err
	}

	for kid, publicKey := range publicKeys {
		publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
		if err != nil {
			s.logger.Sugar().Warnf("Failed to marshal public key. Error: %v", err)
			continue
		}
		publicKeyBlock := &pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyBytes,
		}
		var publicKeyBuffer bytes.Buffer
		err = pem.Encode(&publicKeyBuffer, publicKeyBlock)
		if err != nil {
			s.logger.Sugar().Warnf("Failed to encode public key to string. Error: %v", err)
			continue
		}

		publicKeyStrings[kid] = publicKeyBuffer.String()
	}

	return publicKeyStrings, nil
}
