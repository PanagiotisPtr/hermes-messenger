package token

import (
	"crypto/rsa"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

type Repository struct {
	logger *log.Logger
}

func NewRepository(logger *log.Logger) *Repository {
	return &Repository{
		logger: logger,
	}
}

func (sr *Repository) SignTokenWithClaims(
	data string,
	ttl time.Duration,
	privateKey *rsa.PrivateKey,
) (string, error) {
	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["dat"] = data
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("Failed to sign token %v", err)
	}

	return token, nil
}

func (sr *Repository) ValidateToken(
	tokenString string,
	publicKey *rsa.PublicKey,
) (string, error) {
	token, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return publicKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("Failed to validate token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("Invalid token claims")
	}

	return claims["dat"].(string), nil
}
