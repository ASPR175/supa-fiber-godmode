package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID string) (string, error) {
	secret := os.Getenv("SUPABASE_JWT_SECRET")
	if secret == "" {
		return "", jwt.ErrTokenInvalidClaims // or a custom error
	}

	claims := jwt.MapClaims{
		"sub": userID,                                // subject: the user id
		"exp": time.Now().Add(24 * time.Hour).Unix(), // expiry
		"iat": time.Now().Unix(),                     // issued at
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
