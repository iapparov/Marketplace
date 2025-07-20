package app

import (
	"marketplace/internal/config"
	"testing"

	"github.com/google/uuid"
)

func TestJwtProvider_AccessToken(t *testing.T) {
    cfg := &config.Config{JWT_ACCESS_SECRET: "secret", JWT_EXP_ACCESS_TOKEN: 15}
    jwtProvider := NewJwtProvider(cfg)
    user := User{UUID: uuid.New(), Login: "user"}
    token, err := jwtProvider.GenerateAccessToken(user, cfg)
    if err != nil {
        t.Fatalf("failed to generate token: %v", err)
    }
    claims, err := jwtProvider.ValidateAccessToken(token)
    if err != nil {
        t.Fatalf("failed to validate token: %v", err)
    }
    if claims["uuid"] != user.UUID.String() {
        t.Errorf("expected uuid %v, got %v", user.UUID.String(), claims["uuid"])
    }
}

func TestJwtProvider_RefreshToken(t *testing.T) {
    cfg := &config.Config{JWT_REFRESH_SECRET: "secret", JWT_EXP_REFRESH_TOKEN: 24}
    jwtProvider := NewJwtProvider(cfg)
    user := User{UUID: uuid.New(), Login: "user"}
    token, err := jwtProvider.GenerateRefreshToken(user, cfg)
    if err != nil {
        t.Fatalf("failed to generate token: %v", err)
    }
    claims, err := jwtProvider.ValidateRefreshToken(token)
    if err != nil {
        t.Fatalf("failed to validate token: %v", err)
    }
    if claims["uuid"] != user.UUID.String() {
        t.Errorf("expected uuid %v, got %v", user.UUID.String(), claims["uuid"])
    }
}