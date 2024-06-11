package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenUseCase interface {
	GenerateAccessToken(claims JwtCustomClaims) (string, time.Time, error)
}

type tokenUseCase struct {
	secretKey          string
	expirationDuration time.Duration
}

type JwtCustomClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func NewTokenUseCase(secretKey string, expirationDuration time.Duration) *tokenUseCase {
	return &tokenUseCase{
		secretKey:          secretKey,
		expirationDuration: expirationDuration,
	}
}

func (t *tokenUseCase) GenerateAccessToken(claims JwtCustomClaims) (string string, expiredAt time.Time, err error) {
	// Set default expiration time if not provided
	if claims.ExpiresAt == nil {
		expirationTime := time.Now().Add(t.expirationDuration)
		claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	}
	plainToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	encodedToken, err := plainToken.SignedString([]byte(t.secretKey))

	if err != nil {
		return "", time.Time{}, err
	}
	return encodedToken, claims.ExpiresAt.Time, nil
}
