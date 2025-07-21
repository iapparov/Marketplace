package app

import (
    "testing"
    "marketplace/internal/config"
    "github.com/google/uuid"
    "strings"
)

func TestNewAd_Success(t *testing.T) {
    marketRepo := &MockMarketRepo{}
    userRepo := &MockUserRepo{Users: make(map[string]User)}
    service := NewMarketService(marketRepo, userRepo)
    cfg := config.Config{
        Ad: config.Ad{
            MinLengthTitle: 3, MaxLengthTitle: 100,
            MinLengthDescription: 10, MaxLengthDescription: 1000,
            ImgType: []string{".jpg", ".png"},
            AllowedImgTypesMap: map[string]bool{".jpg": true, ".png": true},
            PriceMin: 1,
        },
    }
    user := User{UUID: uuid.New(), Login: "user", Password: "pass"}
    userRepo.SaveNewUser(user)
    ad := Ad{
        Title: "Test Ad",
        Description: "Test Description for Ad",
        ImageURL: "image.jpg",
        Price: 10,
    }
    created, err := service.NewAd(ad, cfg, user.UUID)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if created.Title != ad.Title {
        t.Errorf("expected title %v, got %v", ad.Title, created.Title)
    }
    if created.Description != ad.Description {
        t.Errorf("expected description %v, got %v", ad.Description, created.Description)
    }
    if created.ImageURL != ad.ImageURL {
        t.Errorf("expected image URL %v, got %v", ad.ImageURL, created.ImageURL)
    }
    if created.Price != ad.Price {
        t.Errorf("expected price %v, got %v", ad.Price, created.Price)
    }
    if created.UserID != user.UUID {
        t.Errorf("expected user ID %v, got %v", user.UUID, created.UserID)
    }
    if created.Username != user.Login {
        t.Errorf("expected username %v, got %v", user.Login, created.Username)
    }
    if created.UUID == uuid.Nil {
        t.Error("expected non-nil UUID for ad")
    }
    if created.CreatedAt.IsZero() {
        t.Error("expected non-zero CreatedAt for ad")
    }
}

func TestNewAd_Fail(t *testing.T) {
	cfg := config.Config{
		Ad: config.Ad{
			MinLengthTitle:       3,
			MaxLengthTitle:       100,
			MinLengthDescription: 10,
			MaxLengthDescription: 1000,
			ImgType:              []string{".jpg", ".jpeg", ".png", ".webm"},
			AllowedImgTypesMap: map[string]bool{
				".jpg":  true,
				".jpeg": true,
				".png":  true,
				".webm": true,
			},
			PriceMin: 0.01,
		},
	}
	user := User{UUID: uuid.New(), Login: "user", Password: "pass"}
	userRepo := &MockUserRepo{Users: make(map[string]User)}
	userRepo.SaveNewUser(user)
	service := NewMarketService(&MockMarketRepo{}, userRepo)

	tests := []struct {
		name string
		ad   Ad
	}{
		{
			name: "too short title",
			ad: Ad{
				Title:       "Hi",
				Description: "Valid description here",
				ImageURL:    "image.jpg",
				Price:       10,
			},
		},
		{
			name: "too long title",
			ad: Ad{
				Title:       strings.Repeat("a", 101),
				Description: "Valid description here",
				ImageURL:    "image.jpg",
				Price:       10,
			},
		},
		{
			name: "too short description",
			ad: Ad{
				Title:       "Valid title",
				Description: "short",
				ImageURL:    "image.jpg",
				Price:       10,
			},
		},
		{
			name: "too long description",
			ad: Ad{
				Title:       "Valid title",
				Description: strings.Repeat("a", 1001),
				ImageURL:    "image.jpg",
				Price:       10,
			},
		},
		{
			name: "invalid image type",
			ad: Ad{
				Title:       "Valid title",
				Description: "Valid description here",
				ImageURL:    "image.zip",
				Price:       10,
			},
		},
		{
			name: "price below minimum",
			ad: Ad{
				Title:       "Valid title",
				Description: "Valid description here",
				ImageURL:    "image.jpg",
				Price:       0.0,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.NewAd(tc.ad, cfg, user.UUID)
			if err == nil {
				t.Fatalf("expected validation error for case: %s", tc.name)
			}
		})
	}
}

func TestAdsList_Empty(t *testing.T) {
    marketRepo := &MockMarketRepo{}
    userRepo := &MockUserRepo{Users: make(map[string]User)}
    service := NewMarketService(marketRepo, userRepo)
    params := AdsListParams{Page: 1, Limit: 10}
    ads, err := service.AdsList(params, uuid.Nil)
    if err != nil && err.Error() != "list is empty" {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(ads) != len(marketRepo.Ads) {
        t.Errorf("expected %d ads, got %d", len(marketRepo.Ads), len(ads))
    }
}

func TestAdsList_Success(t *testing.T) {
    marketRepo := &MockMarketRepo{}
    userRepo := &MockUserRepo{Users: make(map[string]User)}
    service := NewMarketService(marketRepo, userRepo)
    cfg := config.Config{
        Ad: config.Ad{
            MinLengthTitle: 3, MaxLengthTitle: 100,
            MinLengthDescription: 10, MaxLengthDescription: 1000,
            ImgType: []string{".jpg", ".png"},
            AllowedImgTypesMap: map[string]bool{".jpg": true, ".png": true},
            PriceMin: 1,
        },
    }
    user := User{UUID: uuid.New(), Login: "user", Password: "pass"}
    userRepo.SaveNewUser(user)
    ad := Ad{
        Title: "Test Ad",
        Description: "Test Description for Ad",
        ImageURL: "image.jpg",
        Price: 10,
    }
    _, err := service.NewAd(ad, cfg, user.UUID)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    params := AdsListParams{Page: 1, Limit: 10}
    ads, err := service.AdsList(params, uuid.Nil)
    if err != nil && err.Error() != "list is empty" {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(ads) != len(marketRepo.AdsResponse) {
        t.Errorf("expected %d ads, got %d", len(marketRepo.Ads), len(ads))
    }
}