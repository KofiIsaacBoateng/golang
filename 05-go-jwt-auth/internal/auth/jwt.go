package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims

	Role string `json:"role"`
}

func CreateToken(jwtSecret string, userId string, role string) (string, error) {
	// timing
	now := time.Now().UTC();
	exp := now.Add(7 * 24 * time.Hour);

	// define claims
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: userId, // unique identifier
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},

		Role: role,
	}

	// new token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign or encode token
	signed, err := token.SignedString([]byte(jwtSecret))
	if(err != nil) {
		return "", fmt.Errorf("Signing token failed: %w", err)
	}

	return signed, nil;
}