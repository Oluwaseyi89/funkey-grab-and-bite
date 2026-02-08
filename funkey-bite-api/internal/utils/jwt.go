package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("your-secret-key-change-in-production")

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
	return token.SignedString(jwtSecret)
}

// GenerateAdminToken creates a token specifically for admin users
func GenerateAdminToken(userID int, phone string) (string, error) {
	return GenerateTokenWithAdmin(userID, phone, true)
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
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
	return jwtSecret
}

// SetSecret allows changing the JWT secret (useful for testing)
func SetSecret(secret string) {
	jwtSecret = []byte(secret)
}

// package utils

// import (
// 	"time"

// 	"github.com/dgrijalva/jwt-go"
// )

// var jwtSecret = []byte("your-secret-key-change-in-production")

// type Claims struct {
// 	UserID int    `json:"user_id"`
// 	Phone  string `json:"phone"`
// 	jwt.StandardClaims
// }

// func GenerateToken(userID int, phone string) (string, error) {
// 	expirationTime := time.Now().Add(24 * time.Hour)

// 	claims := &Claims{
// 		UserID: userID,
// 		Phone:  phone,
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: expirationTime.Unix(),
// 			IssuedAt:  time.Now().Unix(),
// 			Issuer:    "funkey-grab-bite",
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString(jwtSecret)
// }

// func ValidateToken(tokenString string) (*Claims, error) {
// 	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
// 		return jwtSecret, nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
// 		return claims, nil
// 	}

// 	return nil, jwt.ErrSignatureInvalid
// }
