package authentication

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/panagiotisptr/hermes-messenger/services/authentication/server/config"
	"github.com/panagiotisptr/hermes-messenger/services/authentication/server/keys"
	"github.com/panagiotisptr/hermes-messenger/services/authentication/server/token"
	"github.com/panagiotisptr/hermes-messenger/services/authentication/server/user"
	"go.uber.org/zap"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	refreshKeyPairNamePrefix = "refreshToken"
	accessKeyPairNamePrefix  = "accessToken"
)

type Service struct {
	selfUuid                  uuid.UUID
	logger                    *zap.Logger
	refreshTokenKeyName       string
	accessTokenKeyName        string
	tokenRepository           token.Repository
	keysRepository            keys.Repository
	userRepository            user.Repository
	refreshTokenDuration      time.Duration
	accessTokenDuration       time.Duration
	keyPairGenerationInterval time.Duration
}

func ProvideAuthenticationService(
	cfg *config.Config,
	logger *zap.Logger,
	tokenRepository token.Repository,
	keysRepository keys.Repository,
	userRepository user.Repository,
) (*Service, error) {
	service := &Service{
		selfUuid:                  cfg.UUID,
		logger:                    logger,
		refreshTokenKeyName:       "",
		accessTokenKeyName:        "",
		tokenRepository:           tokenRepository,
		keysRepository:            keysRepository,
		userRepository:            userRepository,
		refreshTokenDuration:      cfg.RefreshTokenDuration,
		accessTokenDuration:       cfg.AccessTokenDuration,
		keyPairGenerationInterval: cfg.KeyPairGenerationInterval,
	}

	service.generateKeyPair()
	go func() {
		for range time.Tick(service.keyPairGenerationInterval) {
			service.generateKeyPair()
		}
	}()

	return service, nil
}

func (s *Service) generateKeyPair() {
	keyUuid := uuid.New().String()
	refreshTokenKeyName := refreshKeyPairNamePrefix + ":" + keyUuid
	accessTokenKeyName := accessKeyPairNamePrefix + ":" + keyUuid

	refreshTokenKeyPair, err := keys.GenerateRSAKeyPair()
	if err != nil {
		s.logger.Sugar().Error(err)
		return
	}

	accessTokenKeyPair, err := keys.GenerateRSAKeyPair()
	if err != nil {
		s.logger.Sugar().Error(err)
		return
	}

	err = s.keysRepository.StoreKeyPair(refreshTokenKeyName, refreshTokenKeyPair, 2*s.keyPairGenerationInterval)
	if err != nil {
		s.logger.Sugar().Error(err)
		return
	}

	err = s.keysRepository.StoreKeyPair(accessTokenKeyName, accessTokenKeyPair, 2*s.keyPairGenerationInterval)
	if err != nil {
		s.logger.Sugar().Error(err)
		return
	}

	s.refreshTokenKeyName = refreshTokenKeyName
	s.accessTokenKeyName = accessTokenKeyName

	s.logger.Sugar().Infof(
		"Updated access and refresh keys for authentication service with UUID: " +
			s.selfUuid.String() +
			". Using keys with UUID:" +
			keyUuid,
	)
}

func (s *Service) Register(
	email string,
	password string,
) (bool, error) {
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return false, err
	}

	err = s.userRepository.CreateUser(user.User{
		Uuid:     uuid.New(),
		Email:    email,
		Password: string(passwordHash),
	})

	return err == nil, err
}

func (s *Service) generateToken(
	data string,
	ttl time.Duration,
	tokenKeyName string,
) (string, error) {
	tokenPrivateKey, err := s.keysRepository.GetPrivateKey(
		tokenKeyName,
	)
	if err != nil {
		s.logger.Sugar().Errorf("[ERROR] Failed to get private key: %v", err)
		return "", fmt.Errorf("Error when generating token")
	}
	token, err := s.tokenRepository.SignTokenWithClaims(
		data,
		time.Hour*24,
		tokenPrivateKey,
	)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) Authenticate(
	email string,
	password string,
) (string, string, error) {
	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.generateToken(
		user.Uuid.String(),
		s.refreshTokenDuration,
		s.refreshTokenKeyName,
	)
	if err != nil {
		return "", "", err
	}

	accessToken, err := s.generateToken(
		user.Uuid.String(),
		s.accessTokenDuration,
		s.accessTokenKeyName,
	)
	if err != nil {
		return refreshToken, "", err
	}

	return refreshToken, accessToken, nil
}

func (s *Service) Refresh(refreshToken string) (string, error) {
	publicKeys, err := s.keysRepository.GetAllPublicKeys(
		s.refreshTokenKeyName,
	)
	if err != nil {
		s.logger.Sugar().Errorf("[ERROR] Failed to get public keys: %v", err)
		return "", fmt.Errorf("Could not validate refresh token")
	}

	for _, publicKey := range publicKeys {
		data, err := s.tokenRepository.ValidateToken(refreshToken, publicKey)
		if err == nil {
			accessToken, err := s.generateToken(
				data,
				time.Hour,
				s.accessTokenKeyName,
			)
			if err != nil {
				return "", err
			}

			return accessToken, nil
		}
	}

	return "", fmt.Errorf("Could not validate refresh token")
}

func (s *Service) GetPublicKeys() ([]string, error) {
	publicKeyStrings := make([]string, 0)
	// Only allow access to the public key for the access tokens
	publicKeys, err := s.keysRepository.GetAllPublicKeys(
		accessKeyPairNamePrefix,
	)
	if err != nil {
		s.logger.Sugar().Warnf("Failed to fetch public keys. Error: %v", err)
		return publicKeyStrings, err
	}

	for _, publicKey := range publicKeys {
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

		publicKeyStrings = append(publicKeyStrings, publicKeyBuffer.String())
	}

	return publicKeyStrings, nil
}
