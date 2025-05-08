package utils

import (
	config "TaskManagmentApis/configs"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims defines the payload structure for the JWT token
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT with custom and registered claims
func GenerateToken(userID, email string, duration time.Duration) (string, error) {
	now := time.Now()

	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "TaskManagerApis",
			Subject:   userID,
			Audience:  jwt.ClaimStrings{"taskmanager-client"},
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(config.Config.JWTSecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return signedToken, nil
}

// GenerateAccessToken generates a short-lived token for authentication
func GenerateAccessToken(userID, email string) (string, error) {
	expiration := time.Duration(config.Config.AccessTokenExpireMinutes) * time.Minute
	return GenerateToken(userID, email, expiration)
}

// GenerateRefreshToken generates a long-lived token for re-authentication
func GenerateRefreshToken(userID, email string) (string, error) {
	expiration := time.Duration(config.Config.RefreshTokenExpireHours) * time.Hour
	return GenerateToken(userID, email, expiration)
}

// ValidateToken parses and validates a JWT string and returns its claims
func ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Config.JWTSecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid or expired token")
	}

	return claims, nil
}
