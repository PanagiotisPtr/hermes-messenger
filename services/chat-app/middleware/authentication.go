package middleware

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"

	"github.com/panagiotisptr/hermes-messenger/chat-app/protos"

	"github.com/golang-jwt/jwt"
)

func validateToken(
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

func Authentication(
	authenticationClient protos.AuthenticationClient,
	handler func(http.ResponseWriter, *http.Request),
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := ""
		for _, c := range r.Cookies() {
			if c.Name == "AccessToken" {
				token = c.Value
				break
			}
		}

		publicKeyResponse, err := authenticationClient.GetPublicKey(
			context.Background(),
			&protos.GetPublicKeyRequest{},
		)
		if err != nil {
			fmt.Println(err)
		}

		block, _ := pem.Decode([]byte(publicKeyResponse.PublicKey))
		key, _ := x509.ParsePKIXPublicKey(block.Bytes)

		userUuid, err := validateToken(token, key.(*rsa.PublicKey))

		if err != nil {
			fmt.Println("ERROR: ", err)
		}

		handler(w, r.Clone(context.WithValue(
			r.Context(),
			"UserUuid",
			userUuid,
		)))
	}
}
