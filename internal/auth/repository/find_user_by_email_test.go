package repository

import (
	"context"
	"database/sql"
	"github.com/bogatyr285/auth-go/internal/auth/entity"
	"reflect"
	"testing"
)

func setupTestDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			password TEXT NOT NULL
		);
	`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		INSERT INTO users (username, password) VALUES ('tuser', 'tpassword');
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestSQLLiteStorage_FindUserByEmail(t *testing.T) {
	type fields struct {
		db *sql.DB
	}

	type args struct {
		ctx      context.Context
		username string
	}

	dbSQLite, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}

	defer dbSQLite.Close()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entity.UserAccount
		wantErr bool
	}{
		{
			name: "User not found",
			fields: fields{
				db: dbSQLite,
			},
			args: args{
				ctx:      context.Background(),
				username: "testuser",
			},
			want:    entity.UserAccount{},
			wantErr: true,
		},
		{
			name: "User found",
			fields: fields{
				db: dbSQLite,
			},
			args: args{
				ctx:      context.Background(),
				username: "tuser",
			},
			want: entity.UserAccount{
				ID:       1,
				Username: "tuser",
				Password: "tpassword",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SQLLiteStorage{
				db: tt.fields.db,
			}
			got, err := s.FindUserByEmail(tt.args.ctx, tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindUserByEmail() got = %v, want %v", got, tt.want)
			}
		})
	}
}
