package datasource

import (
	"database/sql"
	"fmt"
	"marketplace/internal/config"
)


func New(config config.Config) (*Storage, error){
	db, err := sql.Open("sqlite3", config.Db)
	if err != nil {
		return nil, fmt.Errorf("init error DB:%w", err)
	}

	_, err = db.Prepare(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid TEXT NOT NULL,
		login TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`)
	
	if err != nil {
		return nil, fmt.Errorf("prepare error DB:%w", err)
	}		
	return &Storage{
		db: db,
	}, nil
}