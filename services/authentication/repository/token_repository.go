package repository

import (
	"context"

	"github.com/golang-jwt/jwt"
)

type TokenRepository interface {
	// SignTokenWithClaims uses a private RSA key to sign a token with claims
	SignTokenWithClaims(
		ctx context.Context,
		claims jwt.Claims,
	) (string, error)

	// ValidateToken validates a token against a list of public keys
	ValidateToken(
		ctx context.Context,
		token string,
	) (jwt.Claims, error)
}
