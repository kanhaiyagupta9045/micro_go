package helpers

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hashedPassword)
}

func VerifyPassword(providedpassword, userpassword string) (bool, error) {

	check := false
	err := bcrypt.CompareHashAndPassword([]byte(userpassword), []byte(providedpassword))
	if err == nil {
		check = true
	}

	return check, err

}

type SignedDetails struct {
	Id string
	jwt.StandardClaims
}

func GenerateAccessToken(id int) (string, error) {
	ACCESS_TOKEN_SECRET := os.Getenv("ACCESS_TOKEN_SECRET")
	if ACCESS_TOKEN_SECRET == "" {
		return "", errors.New("access token secret is not set")
	}
	accessClaims := &SignedDetails{
		Id: strconv.Itoa(id),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(ACCESS_TOKEN_SECRET))
	if err != nil {
		return "", err
	}

	return token, nil
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
