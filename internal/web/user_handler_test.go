package web

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "bytes"
    "go.uber.org/zap"
    "marketplace/internal/app"
    "marketplace/internal/config"
    "encoding/json"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func TestUserHandler_Register_Success(t *testing.T) {
    repo := &app.MockUserRepo{Users: make(map[string]app.User)}
    service := app.NewUserService(repo)
    cfg := &config.Config{Username: config.Username{MinLength: 3, MaxLength: 20, AllowedCharacters: "A-Za-z0-9_-"}, Password: config.Password{MinLength: 8, MaxLength: 64}}
    logger := zap.NewNop()
    jwt := app.NewJwtProvider(cfg)
    handler := NewUserHandler(service, cfg, jwt, logger)

    body := bytes.NewBufferString(`{"login":"TestUser","password":"Password1"}`)
    req := httptest.NewRequest("POST", "/register", body)
    w := httptest.NewRecorder()
    handler.Register(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("expected 200, got %d", w.Code)
    }
}

func TestUserHandler_Register_Fail(t *testing.T) {
    repo := &app.MockUserRepo{Users: make(map[string]app.User)}
    service := app.NewUserService(repo)
    cfg := &config.Config{Username: config.Username{MinLength: 3, MaxLength: 20, AllowedCharacters: "A-Za-z0-9_-"}, Password: config.Password{MinLength: 8, MaxLength: 64}}
    logger := zap.NewNop()
    jwt := app.NewJwtProvider(cfg)
    handler := NewUserHandler(service, cfg, jwt, logger)

    body := bytes.NewBufferString(`{"login":"Test User","password":"Pass"}`)
    req := httptest.NewRequest("POST", "/register", body)
    w := httptest.NewRecorder()
    handler.Register(w, req)

    if w.Code != http.StatusBadRequest {
        t.Errorf("expected 400, got %d", w.Code)
    }
}
func TestUserHandler_Login(t *testing.T) {
	repo := &app.MockUserRepo{Users: make(map[string]app.User)}
	service := app.NewUserService(repo)
	cfg := &config.Config{}
	logger := zap.NewNop()
	jwt := &app.JwtProvider{}
	handler := NewUserHandler(service, cfg, jwt, logger)

	password := "secret"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := app.User{
		UUID:     uuid.New(),
		Login:    "test",
		Password: string(hashedPassword),
	}
	repo.SaveNewUser(user)

	validReq := app.JwtRequest{Login: "test", Password: password}
	body, _ := json.Marshal(validReq)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler.Login(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	req = httptest.NewRequest("POST", "/login", bytes.NewBufferString("{invalid json"))
	w = httptest.NewRecorder()
	handler.Login(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	invalidReq := app.JwtRequest{Login: "test", Password: "wrongpass"}
	body, _ = json.Marshal(invalidReq)
	req = httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	w = httptest.NewRecorder()
	handler.Login(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401 for wrong password, got %d", w.Code)
	}
}

func TestUserHandler_RefreshAccessToken(t *testing.T) {
	repo := &app.MockUserRepo{Users: make(map[string]app.User)}
    service := app.NewUserService(repo)
	cfg := &config.Config{JWT_ACCESS_SECRET: "secret", JWT_EXP_ACCESS_TOKEN: 15, JWT_REFRESH_SECRET: "secret", JWT_EXP_REFRESH_TOKEN: 24}
	logger := zap.NewNop()
	jwt := &app.JwtProvider{}
	handler := NewUserHandler(service, cfg, jwt, logger)

	user := app.User{
		UUID:     uuid.New(),
		Login:    "refreshUser",
		Password: "somepass",
	}
	repo.SaveNewUser(user)

	refreshToken, err := jwt.GenerateRefreshToken(user, cfg)
	if err != nil {
		t.Fatalf("failed to generate refresh token: %v", err)
	}

	validReq := app.RefreshJwtRequest{RefreshToken: refreshToken}
	body, _ := json.Marshal(validReq)
	req := httptest.NewRequest("POST", "/refresh", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler.RefreshAccessToken(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	req = httptest.NewRequest("POST", "/refresh", bytes.NewBufferString("{bad json"))
	w = httptest.NewRecorder()
	handler.RefreshAccessToken(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	invalidReq := app.RefreshJwtRequest{RefreshToken: "invalidtoken"}
	body, _ = json.Marshal(invalidReq)
	req = httptest.NewRequest("POST", "/refresh", bytes.NewReader(body))
	w = httptest.NewRecorder()
	handler.RefreshAccessToken(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}