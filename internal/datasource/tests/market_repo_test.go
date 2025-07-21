package datasource_test

import (
	"marketplace/internal/app"
	"marketplace/internal/datasource"
	"testing"
	"time"
	_ "github.com/mattn/go-sqlite3"
	"github.com/google/uuid"
	"database/sql"
)
func setupMarketTestDB(t *testing.T) *sql.DB {
	dsn := "file:testdb?mode=memory&cache=shared"
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		t.Fatalf("failed to open test DB: %v", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid TEXT NOT NULL UNIQUE,
		login TEXT NOT NULL,
		password TEXT NOT NULL
	);`)
	if err != nil {
		t.Fatalf("failed to create users table: %v", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS ads (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid TEXT NOT NULL UNIQUE,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		img TEXT,
		user_uuid TEXT NOT NULL,
		price REAL NOT NULL,
		created_at DATETIME NOT NULL,
		FOREIGN KEY(user_uuid) REFERENCES users(uuid)
	);`)
	if err != nil {
		t.Fatalf("failed to create ads table: %v", err)
	}

	return db
}


func TestMarketRepo_SaveAndGetAds(t *testing.T) {
	db := setupMarketTestDB(t)
	userRepo := datasource.NewUserRepo(db)
	adRepo := datasource.NewMarketRepo(db, userRepo)

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

	_, err = adRepo.SaveAd(ad)
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
}