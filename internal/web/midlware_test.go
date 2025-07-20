package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"marketplace/internal/app"
	"marketplace/internal/config"
	"github.com/google/uuid"
)

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	jwt := app.NewJwtProvider(&config.Config{})
	handler := AuthMiddleware(jwt)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called on missing header")
	}))
	req := httptest.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.Code)
	}
}

func TestAuthMiddleware_InvalidHeader(t *testing.T) {
	jwt := app.NewJwtProvider(&config.Config{})
	handler := AuthMiddleware(jwt)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called on invalid header")
	}))
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "BadHeaderValue")
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	jwt := app.NewJwtProvider(&config.Config{})
	handler := AuthMiddleware(jwt)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called on invalid token")
	}))
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.value")
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.Code)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	cfg := &config.Config{
		JWT_ACCESS_SECRET: "testsecret",
		JWT_EXP_ACCESS_TOKEN: 15,
	}
	jwt := app.NewJwtProvider(cfg)
	user := app.User{
		Login:    "testuser",
		Password: "hashedPassword",
		UUID:    uuid.New(),
	}
	token, _ := jwt.GenerateAccessToken(user, cfg)

	called := false
	handler := AuthMiddleware(jwt)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		val := r.Context().Value(UserIDKey)
		if val == nil || val != user.UUID.String() {
			t.Errorf("expected uuid %v in context, got %v", user.UUID, val)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	if !called {
		t.Error("handler was not called")
	}
	if resp.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.Code)
	}
}

func TestOptionalAuthMiddleware_NoHeader(t *testing.T) {
	jwt := app.NewJwtProvider(&config.Config{})
	called := false

	handler := OptionalAuthMiddleware(jwt)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if r.Context().Value(UserIDKey) != nil {
			t.Error("expected no user in context")
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	if !called {
		t.Error("handler was not called")
	}
	if resp.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.Code)
	}
}

func TestOptionalAuthMiddleware_InvalidHeader(t *testing.T) {
	jwt := app.NewJwtProvider(&config.Config{})
	called := false

	handler := OptionalAuthMiddleware(jwt)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if r.Context().Value(UserIDKey) != nil {
			t.Error("expected no user in context")
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "InvalidHeaderValue")
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	if !called {
		t.Error("handler was not called")
	}
	if resp.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.Code)
	}
}

func TestOptionalAuthMiddleware_InvalidToken(t *testing.T) {
	jwt := app.NewJwtProvider(&config.Config{})
	called := false

	handler := OptionalAuthMiddleware(jwt)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if r.Context().Value(UserIDKey) != nil {
			t.Error("expected no user in context for invalid token")
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.value")
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	if !called {
		t.Error("handler was not called")
	}
	if resp.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.Code)
	}
}

func TestOptionalAuthMiddleware_ValidToken(t *testing.T) {
	cfg := &config.Config{
		JWT_ACCESS_SECRET:       "testsecret",
		JWT_EXP_ACCESS_TOKEN:    15,
	}
	jwt := app.NewJwtProvider(cfg)
	user := app.User{
		Login:    "testuser",
		Password: "hashedPassword",
		UUID:    uuid.New(),
	}
	token, _ := jwt.GenerateAccessToken(user, cfg)

	called := false
	handler := OptionalAuthMiddleware(jwt)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		val := r.Context().Value(UserIDKey)
		if val != user.UUID.String() {
			t.Errorf("expected uuid %v in context, got %v", user.UUID.String(), val)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	if !called {
		t.Error("handler was not called")
	}
	if resp.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.Code)
	}
}