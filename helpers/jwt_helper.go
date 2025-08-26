package helpers

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateJWT creates a signed JWT for a user
func GenerateJWT(userID, email string, isAdmin bool) (string, error) {
	claims := jwt.MapClaims{
		"ID":    userID,
		"email": email,
		"admin": isAdmin,
		"exp":   time.Now().Add(24 * time.Hour).Unix(), // token valid for 1 day
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("COOKIE_SECRET")
	return token.SignedString([]byte(secret))
}

// ParseJWT verifies a token's signature
func ParseJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("COOKIE_SECRET")), nil
	})
}

// ParseJWTClaims returns the claims map directly
func ParseJWTClaims(tokenString string) (map[string]interface{}, error) {
	token, err := ParseJWT(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
