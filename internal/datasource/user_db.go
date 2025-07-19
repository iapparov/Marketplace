package datasource

import (
	"marketplace/internal/app"
	"fmt"
	"database/sql"
	"strings"
)

type UserRepo struct{
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (s *UserRepo) SaveNewUser(user app.User) error {
	stmt, err := s.db.Prepare(`INSERT INTO users (uuid, login, password) VALUES (?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("prepare error DB:%w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.UUID, strings.ToLower(user.Login), user.Password)
	if err != nil {
		return fmt.Errorf("exec error DB:%w", err)
	}

	return nil
}

func (s *UserRepo) FindByLogin(login string) (app.User, error) {
	var user app.User
	row := s.db.QueryRow(`SELECT uuid, login, password FROM users WHERE login = ?`, strings.ToLower(login))
	err := row.Scan(&user.UUID, &user.Login, &user.Password)
	if err != nil {
		return app.User{}, fmt.Errorf("scan error DB:%w", err)
	}
	return user, nil
}

func (s *UserRepo) FindByUUID(uuid string) (app.User, error) {
	var user app.User
	row := s.db.QueryRow(`SELECT uuid, login, password FROM users WHERE uuid = ?`, uuid)
	err := row.Scan(&user.UUID, &user.Login, &user.Password)
	if err != nil {
		return app.User{}, fmt.Errorf("scan error DB:%w", err)
	}
	return user, nil
}
