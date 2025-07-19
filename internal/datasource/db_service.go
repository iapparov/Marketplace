package datasource

import (
	"database/sql"
	"fmt"
	"marketplace/internal/config"
	_ "github.com/mattn/go-sqlite3"
)



func NewStorage(config *config.Config) (*sql.DB, error){
	db, err := sql.Open("sqlite3", config.Db)
	if err != nil {
		return nil, fmt.Errorf("init error DB:%w", err)
	}

    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        uuid TEXT NOT NULL,
        login TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL
    );`)
    if err != nil {
        return nil, fmt.Errorf("create users table error: %w", err)
    }

    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS ads (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        uuid TEXT NOT NULL,
        title TEXT NOT NULL,
        description TEXT NOT NULL,
        price REAL NOT NULL,
        img TEXT NOT NULL,
        user_uuid TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_uuid) REFERENCES users(uuid)
    );`)
    if err != nil {
        return nil, fmt.Errorf("create ads table error: %w", err)
	}

	return db, nil
}