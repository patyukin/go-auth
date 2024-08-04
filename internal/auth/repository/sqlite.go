package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bogatyr285/auth-go/internal/auth/entity"
	_ "github.com/mattn/go-sqlite3"
)

type SQLLiteStorage struct {
	db *sql.DB
}

func New(dbPath string) (SQLLiteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return SQLLiteStorage{}, err
	}
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		username text not null,
		password text not null,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);
	create index if not exists idx_username ON users(username);
	`)
	if err != nil {
		return SQLLiteStorage{}, fmt.Errorf("db schema init err: %s", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return SQLLiteStorage{}, err
	}
	return SQLLiteStorage{db: db}, nil
}

func (s *SQLLiteStorage) Close() error {
	return s.db.Close()
}

func (s *SQLLiteStorage) RegisterUser(ctx context.Context, u entity.UserAccount) error {
	stmt, err := s.db.PrepareContext(ctx, `INSERT INTO users(username, password) VALUES(?,?)`)
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(u.Username, u.Password); err != nil {
		return err
	}

	return nil
}

func (s *SQLLiteStorage) FindUserByEmail(ctx context.Context, username string) (entity.UserAccount, error) {
	stmt, err := s.db.PrepareContext(ctx, `SELECT password FROM users WHERE username = ?`)
	if err != nil {
		return entity.UserAccount{}, err
	}

	var pswdFromDB string

	if err := stmt.QueryRow(username).Scan(&pswdFromDB); err != nil {
		return entity.UserAccount{}, err
	}

	return entity.UserAccount{
		Username: username,
		Password: pswdFromDB,
	}, nil
}
