package datasource_test

import (
	"database/sql"
	"marketplace/internal/app"
	"marketplace/internal/config"
	"marketplace/internal/datasource"
	"testing"
	_ "github.com/mattn/go-sqlite3"
	"github.com/google/uuid"
)

func setupTestDB(t *testing.T) *sql.DB {
	cfg := &config.Config{Db: ":memory:"}
	db, err := datasource.NewStorage(cfg)
	if err != nil {
		t.Fatalf("failed to setup test DB: %v", err)
	}
	return db
}

func TestUserRepo_SaveAndFind(t *testing.T) {
	db := setupTestDB(t)
	repo := datasource.NewUserRepo(db)

	user := app.User{
		UUID:     uuid.New(),
		Login:    "testuser",
		Password: "hashedPassword",
	}

	err := repo.SaveNewUser(user)
	if err != nil {
		t.Fatalf("failed to save user: %v", err)
	}

	found, err := repo.FindByLogin("TestUser")
	if err != nil {
		t.Fatalf("failed to find user by login: %v", err)
	}
	if found.UUID != user.UUID {
		t.Errorf("expected UUID %s, got %s", user.UUID, found.UUID)
	}

	found2, err := repo.FindByUUID(user.UUID.String())
	if err != nil {
		t.Fatalf("failed to find user by uuid: %v", err)
	}
	if found2.Login != user.Login {
		t.Errorf("expected login %s, got %s", user.Login, found2.Login)
	}
}