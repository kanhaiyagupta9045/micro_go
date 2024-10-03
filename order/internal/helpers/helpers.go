package helpers

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type SignedDetails struct {
	Id string
	jwt.StandardClaims
}

func ValidateToken(signedToken string) (*SignedDetails, error) {
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		secretKey := os.Getenv("ACCESS_TOKEN_SECRET")
		if secretKey == "" {
			return nil, errors.New("SECRET_KEY not set in environment")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims or token is not valid")
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
