package utils

import (
	"errors"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	jwtSecret   []byte
	jwtSecretMu sync.RWMutex
)

type Claims struct {
	UserID  int    `json:"user_id"`
	Phone   string `json:"phone"`
	IsAdmin bool   `json:"is_admin"`
	jwt.StandardClaims
}

// GenerateToken creates a new JWT token for regular users
func GenerateToken(userID int, phone string) (string, error) {
	return GenerateTokenWithAdmin(userID, phone, false)
}

// GenerateTokenWithAdmin creates a new JWT token with admin status
func GenerateTokenWithAdmin(userID int, phone string, isAdmin bool) (string, error) {
	secret, err := currentJWTSecret()
	if err != nil {
		return "", err
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID:  userID,
		Phone:   phone,
		IsAdmin: isAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "funkey-grab-bite",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// GenerateAdminToken creates a token specifically for admin users
func GenerateAdminToken(userID int, phone string) (string, error) {
	return GenerateTokenWithAdmin(userID, phone, true)
}

func ValidateToken(tokenString string) (*Claims, error) {
	secret, err := currentJWTSecret()
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// GetSecret returns the JWT secret (for testing or special cases)
func GetSecret() []byte {
	jwtSecretMu.RLock()
	defer jwtSecretMu.RUnlock()

	return append([]byte(nil), jwtSecret...)
}

// SetSecret allows changing the JWT secret (useful for testing)
func SetSecret(secret string) {
	jwtSecretMu.Lock()
	defer jwtSecretMu.Unlock()

	jwtSecret = []byte(secret)
	if len(jwtSecret) == 0 {
		jwtSecret = nil
	}
}

// ConfigureJWTSecret validates and applies the runtime JWT secret.
func ConfigureJWTSecret(secret string) error {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return errors.New("JWT_SECRET is required")
	}

	SetSecret(secret)
	return nil
}

// ConfigureJWTSecretFromEnv loads the JWT secret from the environment.
func ConfigureJWTSecretFromEnv() error {
	return ConfigureJWTSecret(os.Getenv("JWT_SECRET"))
}

// GenerateAdminJWT creates a JWT token specifically for admin dashboard access
func GenerateAdminJWT(adminID int, email, role string) (string, error) {
	secret, err := currentJWTSecret()
	if err != nil {
		return "", err
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID:  adminID,
		Phone:   email, // Using email as phone field for admin
		IsAdmin: true,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "funkey-admin-dashboard",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func currentJWTSecret() ([]byte, error) {
	jwtSecretMu.RLock()
	defer jwtSecretMu.RUnlock()

	if len(jwtSecret) == 0 {
		return nil, errors.New("jwt secret is not configured")
	}

	return append([]byte(nil), jwtSecret...), nil
}
