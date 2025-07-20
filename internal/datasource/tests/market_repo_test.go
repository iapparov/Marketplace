package datasource_test

import (
	"marketplace/internal/app"
	"marketplace/internal/datasource"
	"testing"
	"time"
	_ "github.com/mattn/go-sqlite3"
	"github.com/google/uuid"
)


func TestMarketRepo_SaveAndGetAds(t *testing.T) {
	db := setupTestDB(t)
	userRepo := datasource.NewUserRepo(db)
	adRepo := datasource.NewMarketRepo(db)

	user := app.User{
		UUID:     uuid.New(),
		Login:    "testuser",
		Password: "secret",
	}
	err := userRepo.SaveNewUser(user)
	if err != nil {
		t.Fatalf("failed to save user: %v", err)
	}

	ad := app.Ad{
		ID:          0,
		UUID:        uuid.New(),
		Title:       "Test Ad",
		Description: "This is a test ad",
		Price:       9.99,
		ImageURL:    "img.jpg",
		UserID:      user.UUID,
		CreatedAt:   time.Now(),
	}

	saved, err := adRepo.SaveAd(ad)
	if err != nil {
		t.Fatalf("failed to save ad: %v", err)
	}

	ads, err := adRepo.GetAdsList(app.AdsListParams{
		MinPrice: 0,
		MaxPrice: 100,
		Page:     1,
		Limit:    10,
		SortBy:   "created_at",
		Order:    "asc",
	}, user.UUID.String())
	if err != nil {
		t.Fatalf("failed to get ads list: %v", err)
	}

	if len(ads) != 1 {
		t.Errorf("expected 1 ad, got %d", len(ads))
	}

	if ads[0].UUID != saved.UUID {
		t.Errorf("expected ad UUID %s, got %s", saved.UUID, ads[0].UUID)
	}
}