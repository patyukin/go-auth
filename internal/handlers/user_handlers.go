package handlers

import (
	"context"
	"net/http"

	"github.com/bogatyr285/auth-go/internal/storage"
)

type UserRepository interface {
	RegisterUser(ctx context.Context, u storage.UserAccount) (storage.UserAccount, error)
	Login(ctx context.Context, username, password string) (storage.UserAccount, error)
}

func RegisterHandler(ur UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Transport layer tasks:
		// parse query/body
		// validate
		// call app layer
		// render response
	}
}

func LoginHandler(ur UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
