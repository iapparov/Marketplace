package app

import (
	"marketplace/internal/config"
	"errors"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)



func NewJwtProvider(config *config.Config) *JwtProvider {
	return &JwtProvider{
		accessSecret: []byte(config.JWT_ACCESS_SECRET),
		refreshSecret: []byte(config.JWT_REFRESH_SECRET),
	}
}

func (j *JwtProvider) GenerateAccessToken(user User, config *config.Config) (string, error) {
	claims := jwt.MapClaims{
		"uuid":  user.UUID.String(),
        "exp": time.Now().Add(time.Minute * time.Duration(config.JWT_EXP_ACCESS_TOKEN)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.accessSecret)
}

func (j *JwtProvider) GenerateRefreshToken(user User, config *config.Config) (string, error) {
	claims := jwt.MapClaims{
		"uuid":  user.UUID.String(),
        "exp": time.Now().Add(time.Hour * time.Duration(config.JWT_EXP_REFRESH_TOKEN)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.refreshSecret)
}

func (j *JwtProvider) ValidateAccessToken(tokenStr string) (uuid.UUID, error) {
    
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
        return j.accessSecret, nil
    })
    if err != nil || !token.Valid {
        return uuid.Nil, errors.New("invalid access token")
    }
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return uuid.Nil, errors.New("invalid claims")
    }
    id, err := uuid.Parse(claims["uuid"].(string))
    if err != nil {
        return uuid.Nil, err
    }
    return id, nil
}

func (j *JwtProvider) ValidateRefreshToken(tokenStr string) (uuid.UUID, error) {
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
        return j.refreshSecret, nil
    })
    if err != nil || !token.Valid {
        return uuid.Nil, errors.New("invalid access token")
    }
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return uuid.Nil, errors.New("invalid claims")
    }
    id, err := uuid.Parse(claims["uuid"].(string))
    if err != nil {
        return uuid.Nil, err
    }
    return id, nil
}