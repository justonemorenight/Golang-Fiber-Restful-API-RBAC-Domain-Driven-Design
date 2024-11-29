package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"backend-fiber/internal/pkg/config"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID int32  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessExp     time.Duration
	RefreshExp    time.Duration
}

var jwtConfig *JWTConfig

func InitJWTConfig(cfg *config.Config) {
	jwtConfig = &JWTConfig{
		AccessSecret:  cfg.JWTSecret,
		RefreshSecret: cfg.JWTRefreshSecret,
		AccessExp:     cfg.AccessTokenExp,
		RefreshExp:    cfg.RefreshTokenExp,
	}
}

func GenerateToken(userID int32, email string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtConfig.AccessExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtConfig.AccessSecret))
}

func GenerateRefreshToken(userID int32, email string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtConfig.RefreshExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtConfig.RefreshSecret))
}

func GetJWTConfig() *JWTConfig {
	return jwtConfig
}

func ExtractBearerToken(auth string) string {
	if auth == "" {
		return ""
	}
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	return parts[1]
}

func ValidateToken(token string, secret string) (*JWTClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(*JWTClaims); ok && parsedToken.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

type Claims struct {
	UserID int32
	Email  string
}
