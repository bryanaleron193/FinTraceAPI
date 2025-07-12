package utils

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

// CustomClaims represents the JWT claims
type CustomClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a JWT for a user
func GenerateJWT(userID uint) (string, error) {
	// Define JWT expiration time (e.g., 72 hours)
	expirationTime := time.Now().Add(72 * time.Hour)

	// Create the JWT claims, which includes the user ID and expiry time
	claims := &CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create the token with the claims and sign it with the secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("SECRET_KEY")

	// Sign the token
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseJWT validates the JWT and returns the claims
func ParseJWT(tokenString string) (*CustomClaims, error) {
	secret := os.Getenv("SECRET_KEY")

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// Extract the claims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}