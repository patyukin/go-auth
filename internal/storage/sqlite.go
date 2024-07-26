package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// TODO to models package
type UserAccount struct {
	Username string
	Password string
}

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
		password text not null);
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

func (s *SQLLiteStorage) RegisterUser(ctx context.Context, u UserAccount) (UserAccount, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return UserAccount{}, err
	}

	stmt, err := s.db.PrepareContext(ctx, `INSERT INTO users(username, password) VALUES(?,?)`)
	if err != nil {
		return UserAccount{}, err
	}

	if _, err := stmt.Exec(u.Username, hashedPassword); err != nil {
		return UserAccount{}, err
	}

	return UserAccount{}, nil
}

func (s *SQLLiteStorage) Login(ctx context.Context, username, password string) (UserAccount, error) {
	stmt, err := s.db.PrepareContext(ctx, `SELECT password FROM users WHERE username = ?`)
	if err != nil {
		return UserAccount{}, err
	}

	pswdFromDB := ""

	if err := stmt.QueryRow(username).Scan(&pswdFromDB); err != nil {
		return UserAccount{}, err
	}

	log.Println(pswdFromDB)

	return UserAccount{}, nil
}
