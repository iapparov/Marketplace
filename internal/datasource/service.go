package datasource

import (
	"database/sql"
	"fmt"
	"marketplace/internal/config"
)


func New(config config.Config) (*Storage, error){
	db, err := sql.Open("sqlite3", config.Db)
	if err != nil {
		return nil, fmt.Errorf("Init Error DB:%w", err)
	}

	db.Prepare(`CRATE TABLE IF NOT EXISTS `)
	
	return &Storage{
		db:db,
	}, nil
}