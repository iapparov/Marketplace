package app

import (
	"marketplace/internal/config"
	"testing"
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
)

func TestRegisterUser_Success(t *testing.T) {
	repo := &MockUserRepo{Users: make(map[string]User)}
	service := NewUserService(repo)
	cfg := &config.Config{
		Username: config.Username{MinLength: 3, MaxLength: 20, AllowedCharacters: "A-Za-z0-9_-"},
		Password: config.Password{MinLength: 8, MaxLength: 64, RequireUpper: true, RequireLower: true, RequireDigit: true},
	}
	req := SignUpRequest{Login: "TestUser", Password: "Password1"}
	user, err := service.RegisterUser(req, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Login == "" || user.UUID == uuid.Nil {
		t.Errorf("invalid user returned")
	}
}

func TestRegisterUser_InvalidPassword(t *testing.T) {
	repo := &MockUserRepo{Users: make(map[string]User)}
	service := NewUserService(repo)
	cfg := &config.Config{
		Username: config.Username{MinLength: 3, MaxLength: 20, AllowedCharacters: "A-Za-z0-9_-"},
		Password: config.Password{
			MinLength:    8,
			MaxLength:    64,
			RequireUpper: true,
			RequireLower: true,
			RequireDigit: true,
		},
	}

	cases := []struct {
		name     string
		password string
	}{
        {"empty", ""},
        {"symbol", "ḍ̇"},
		{"no uppercase", "password1"},
        {"no lowercase", "PASSWORD1"},
        {"no digit", "Password"},
        {"with space", "Pass word"},
        {"too short", "Pass1"},
        {"too long", "ThisPasswordIsWayTooLongAndShouldFail1234567890ThisPasswordIsWayTooLongAndShouldFail1234567890ThisPasswordIsWayTooLongAndShouldFail1234567890"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := SignUpRequest{Login: "TestUser", Password: tc.password}
			user, err := service.RegisterUser(req, cfg)

			if err == nil {
				t.Fatal("expected error, got nil case: " + tc.name)
			}
			if user.Login != "" || user.UUID != uuid.Nil {
				t.Errorf("expected empty user, got %+v", user)
			}
		})
	}
}


func TestRegisterUser_InvalidLogin(t *testing.T) {
	repo := &MockUserRepo{Users: make(map[string]User)}
	service := NewUserService(repo)
	cfg := &config.Config{
		Username: config.Username{MinLength: 3, MaxLength: 20, AllowedCharacters: "A-Za-z0-9_-"},
		Password: config.Password{
			MinLength:    8,
			MaxLength:    64,
			RequireUpper: true,
			RequireLower: true,
			RequireDigit: true,
		},
	}

	cases := []struct {
		name     string
		login    string
	}{
        {"empty", ""},
		{"symbols", "***"},
        {"too short", "lo"},
        {"too long", "ThisLoginIsWayTooLongAndShouldFail1234567890ThisLoginIsWayTooLongAndShouldFail1234567890ThisLoginIsWayTooLongAndShouldFail1234567890"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := SignUpRequest{Login: tc.login, Password: "Password1"}
			user, err := service.RegisterUser(req, cfg)

			if err == nil {
				t.Fatal("expected error, got nil case: " + tc.name)
			}
			if user.Login != "" || user.UUID != uuid.Nil {
				t.Errorf("expected empty user, got %+v", user)
			}
		})
	}
}

func TestRegisterUser_Duplicate(t *testing.T) {
	repo := &MockUserRepo{Users: make(map[string]User)}
	service := NewUserService(repo)
	cfg := &config.Config{
		Username: config.Username{MinLength: 3, MaxLength: 20, AllowedCharacters: "A-Za-z0-9_-"},
		Password: config.Password{MinLength: 8, MaxLength: 64, RequireUpper: false, RequireLower: false, RequireDigit: false},
	}
	req := SignUpRequest{Login: "TestUser", Password: "Password1"}
	_, _ = service.RegisterUser(req, cfg)
	_, err := service.RegisterUser(req, cfg)
	if err == nil {
		t.Errorf("expected error for duplicate user")
	}
}

func TestUserService_LoginJwt(t *testing.T) {
	password := "StrongPass1"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := User{
		UUID:     uuid.New(),
		Login:    "testuser",
		Password: string(hashed),
	}
	repo := &MockUserRepo{
		Users: map[string]User{"testuser": user},
	}
	cfg := &config.Config{JWT_ACCESS_SECRET: "secret", JWT_EXP_ACCESS_TOKEN: 15}
	jwtProvider := NewJwtProvider(cfg)
	service := NewUserService(repo)

	t.Run("success", func(t *testing.T) {
		req := JwtRequest{Login: "testuser", Password: password}
		resp, err := service.LoginJwt(req, jwtProvider, cfg)
		if err != nil {
			t.Fatalf("expected success, got error: %v", err)
		}
		if resp.AccessToken == "" || resp.RefreshToken == "" {
			t.Error("expected tokens to be generated")
		}
	})

	t.Run("user not found", func(t *testing.T) {
		req := JwtRequest{Login: "unknown", Password: password}
		_, err := service.LoginJwt(req, jwtProvider, cfg)
		if err == nil || err.Error() != "user not found" {
			t.Errorf("expected user not found error, got: %v", err)
		}
	})

	t.Run("unauthorized (wrong password)", func(t *testing.T) {
		req := JwtRequest{Login: "testuser", Password: "wrong"}
		_, err := service.LoginJwt(req, jwtProvider, cfg)
		if err == nil || err.Error() != "unauthorized" {
			t.Errorf("expected unauthorized error, got: %v", err)
		}
	})
}

func TestUserService_RefreshAccessToken(t *testing.T) {
	password := "StrongPass1"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := User{
		UUID:     uuid.New(),
		Login:    "testuser",
		Password: string(hashed),
	}
	cfg := &config.Config{JWT_REFRESH_SECRET: "secret", JWT_EXP_REFRESH_TOKEN: 24}
	jwtProvider := NewJwtProvider(cfg)

	repo := &MockUserRepo{
		Users: map[string]User{user.UUID.String(): user},
	}
	service := NewUserService(repo)

	refresh, _ := jwtProvider.GenerateRefreshToken(user, cfg)

	t.Run("success", func(t *testing.T) {
		req := RefreshJwtRequest{RefreshToken: refresh}
		resp, err := service.RefreshAccessToken(req, jwtProvider, cfg)
		if err != nil {
			t.Fatalf("expected success, got error: %v", err)
		}
		if resp.AccessToken == "" || resp.RefreshToken == "" {
			t.Error("expected tokens to be generated")
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		req := RefreshJwtRequest{RefreshToken: "invalid.token.here"}
		_, err := service.RefreshAccessToken(req, jwtProvider, cfg)
		if err == nil || err.Error() != "invalid refresh token" {
			t.Errorf("expected invalid token error, got: %v", err)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		// подменяем репо на пустое
		emptyRepo := &MockUserRepo{}
		service := NewUserService(emptyRepo)

		req := RefreshJwtRequest{RefreshToken: refresh}
		_, err := service.RefreshAccessToken(req, jwtProvider, cfg)
		if err == nil || err.Error() != "user not found" {
			t.Errorf("expected user not found error, got: %v", err)
		}
	})
}