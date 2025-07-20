package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"marketplace/internal/app"
	"marketplace/internal/config"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"context"
)

type MockMarketService struct {
	NewAdFunc   func(ad app.Ad, cfg config.Config, userID uuid.UUID) (app.Ad, error)
	AdsListFunc func(params app.AdsListParams, userID uuid.UUID) ([]app.Ad, error)
}

func (m *MockMarketService) NewAd(ad app.Ad, cfg config.Config, userID uuid.UUID) (app.Ad, error) {
	return m.NewAdFunc(ad, cfg, userID)
}

func (m *MockMarketService) AdsList(params app.AdsListParams, userID uuid.UUID) ([]app.Ad, error) {
	return m.AdsListFunc(params, userID)
}

func TestMarketHandler_NewAd_Success(t *testing.T) {
	mockService := &MockMarketService{
		NewAdFunc: func(ad app.Ad, cfg config.Config, userID uuid.UUID) (app.Ad, error) {
			ad.UUID = uuid.New()
			ad.UserID = userID
			return ad, nil
		},
	}
	handler := NewMarketHandler(mockService, &config.Config{}, zap.NewNop())

	ad := app.Ad{
		Title:       "Test Ad",
		Description: "Desc",
		ImageURL:    "img.png",
		Price:       100,
	}
	body, _ := json.Marshal(ad)
	req := httptest.NewRequest("POST", "/ads", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), UserIDKey, uuid.New().String())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.NewAd(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}

	var result app.Ad
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	if result.Title != ad.Title {
		t.Errorf("expected title %s, got %s", ad.Title, result.Title)
	}
}

func TestMarketHandler_NewAd_InvalidBody(t *testing.T) {
	handler := NewMarketHandler(nil, &config.Config{}, zap.NewNop())
	req := httptest.NewRequest("POST", "/ads", bytes.NewBufferString("invalid json"))
	ctx := context.WithValue(req.Context(), UserIDKey, uuid.New().String())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.NewAd(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestMarketHandler_NewAd_InvalidUUID(t *testing.T) {
	handler := NewMarketHandler(nil, &config.Config{}, zap.NewNop())

	ad := app.Ad{Title: "Ad"}
	body, _ := json.Marshal(ad)
	req := httptest.NewRequest("POST", "/ads", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), UserIDKey, "not-a-uuid")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.NewAd(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestMarketHandler_AdsList_Success(t *testing.T) {
	mockService := &MockMarketService{
		AdsListFunc: func(params app.AdsListParams, userID uuid.UUID) ([]app.Ad, error) {
			return []app.Ad{
				{Title: "Ad1", UUID: uuid.New(), CreatedAt: time.Now()},
			}, nil
		},
	}
	handler := NewMarketHandler(mockService, &config.Config{}, zap.NewNop())

	req := httptest.NewRequest("GET", "/ads?page=1&limit=5&sort_by=price&order=desc", nil)
	ctx := context.WithValue(req.Context(), UserIDKey, uuid.New().String())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.AdsList(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}
	var ads []app.Ad
	if err := json.NewDecoder(w.Body).Decode(&ads); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}
	if len(ads) != 1 || ads[0].Title != "Ad1" {
		t.Errorf("unexpected ads response: %+v", ads)
	}
}

func TestMarketHandler_AdsList_ServiceError(t *testing.T) {
	mockService := &MockMarketService{
		AdsListFunc: func(params app.AdsListParams, userID uuid.UUID) ([]app.Ad, error) {
			return nil, errors.New("failed to fetch")
		},
	}
	handler := NewMarketHandler(mockService, &config.Config{}, zap.NewNop())

	req := httptest.NewRequest("GET", "/ads?page=1&limit=5", nil)
	w := httptest.NewRecorder()

	handler.AdsList(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}