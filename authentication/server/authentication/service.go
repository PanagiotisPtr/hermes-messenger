package authentication

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"panagiotisptr/authentication/server/secret"
	"panagiotisptr/authentication/server/token"
	"panagiotisptr/authentication/server/user"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	logger              *log.Logger
	refreshTokenKeyName string
	accessTokenKeyName  string
	tokenRepository     token.Repository
	secretRepository    secret.Repository
	userRepository      user.Repository
}

func NewService(logger *log.Logger) (*Service, error) {
	tokenRepository := token.NewRepository(logger)
	secretRepository := secret.NewMemoryRepository(logger)
	userRepository := user.NewMemoryRepository(logger)

	refreshTokenKeyName := "refreshTokenKeyPair"
	accessTokenKeyName := "accessTokenKeyPair"

	refreshTokenKeyPair, err := secretRepository.GenerateRSAKeyPair()
	if err != nil {
		return nil, err
	}

	accessTokenKeyPair, err := secretRepository.GenerateRSAKeyPair()
	if err != nil {
		return nil, err
	}

	err = secretRepository.StoreKeyPair(refreshTokenKeyName, refreshTokenKeyPair)
	if err != nil {
		return nil, err
	}

	err = secretRepository.StoreKeyPair(accessTokenKeyName, accessTokenKeyPair)
	if err != nil {
		return nil, err
	}

	return &Service{
		logger:              logger,
		refreshTokenKeyName: refreshTokenKeyName,
		accessTokenKeyName:  accessTokenKeyName,
		tokenRepository:     *tokenRepository,
		secretRepository:    secretRepository,
		userRepository:      userRepository,
	}, nil
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
	tokenPrivateKey, err := s.secretRepository.GetPrivateKey(
		tokenKeyName,
	)
	if err != nil {
		s.logger.Printf("[ERROR] Failed to get private key: %v", err)
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
		time.Hour*24,
		s.refreshTokenKeyName,
	)
	if err != nil {
		return "", "", err
	}

	accessToken, err := s.generateToken(
		user.Uuid.String(),
		time.Hour*24,
		s.accessTokenKeyName,
	)
	if err != nil {
		return refreshToken, "", err
	}

	return refreshToken, accessToken, nil
}

func (s *Service) Refresh(refreshToken string) (string, error) {
	publicKey, err := s.secretRepository.GetPublicKey(
		s.refreshTokenKeyName,
	)
	if err != nil {
		s.logger.Printf("[ERROR] Failed to get public key: %v", err)
		return "", fmt.Errorf("Error when generating token")
	}
	data, err := s.tokenRepository.ValidateToken(refreshToken, publicKey)
	if err != nil {
		return "", err
	}

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

func (s *Service) GetPublicKey() (string, error) {
	// Only allow access to the public key for the access tokens
	publicKey, err := s.secretRepository.GetPublicKey(
		s.accessTokenKeyName,
	)
	if err != nil {
		return "", err
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	publicKeyBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	var publicKeyBuffer bytes.Buffer
	err = pem.Encode(&publicKeyBuffer, publicKeyBlock)

	return publicKeyBuffer.String(), nil
}
