package utils

import (
	"errors"
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func getJWTSecret() ([]byte, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET is empty")
	}
	return []byte(secret), nil
}

func GenerateJWT(userID uuid.UUID, role string) (string, error) {
	secret, err := getJWTSecret()
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"role": role,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}