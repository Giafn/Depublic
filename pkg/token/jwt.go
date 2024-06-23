package token

import (
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenUseCase interface {
	GenerateAccessToken(claims JwtCustomClaims) (string, time.Time, error)
	CreateClaims(userId string, email string, role string) JwtCustomClaims
	GetClaimsFromToken(tokenString string) (JwtCustomClaims, error)
	InvalidateToken(tokenString string) error
	IsTokenBlacklisted(tokenString string) bool
}

type tokenUseCase struct {
	secretKey          string
	expirationDuration time.Duration
	blacklist          map[string]bool
	mu                 sync.Mutex
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
		blacklist:          make(map[string]bool),
	}
}

func (t *tokenUseCase) GenerateAccessToken(claims JwtCustomClaims) (string, time.Time, error) {
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

func (t *tokenUseCase) CreateClaims(userId string, email string, role string) JwtCustomClaims {
	return JwtCustomClaims{
		ID:    userId,
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "Depublic-App",
		},
	}
}

func (t *tokenUseCase) GetClaimsFromToken(tokenString string) (JwtCustomClaims, error) {
	claims := JwtCustomClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.secretKey), nil
	})
	if err != nil {
		return JwtCustomClaims{}, err
	}
	return claims, nil
}

func (t *tokenUseCase) InvalidateToken(tokenString string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	_, err := t.GetClaimsFromToken(tokenString)
	if err != nil {
		return err
	}

	t.blacklist[tokenString] = true
	return nil
}

func (t *tokenUseCase) IsTokenBlacklisted(tokenString string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.blacklist[tokenString]
}
