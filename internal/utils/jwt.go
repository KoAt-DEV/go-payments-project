package utils

import (
	"errors"
	"time"

	"go-payments-portfolio-project/internal/config"
	"go-payments-portfolio-project/internal/domain/user"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateTokens(user *user.User, cfg *config.Config) (accessToken, refreshToken string, err error) {
	now := time.Now().UTC()

	accessClaims := Claims{
		UserID: user.ID.String(),
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(cfg.JWT.AccesExpire)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = access.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", "", err
	}

	refreshClaims := jwt.RegisteredClaims{
		Subject:   user.ID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(cfg.JWT.RefreshExpire)),
	}

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	refreshToken, err = refresh.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil

}

func ParseAccessToken(tokenStr string, cfg config.Config) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
