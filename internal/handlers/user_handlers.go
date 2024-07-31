package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/bogatyr285/auth-go/internal/storage"
)

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=1"`
}

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
		req := &RegisterRequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			http.Error(w, "parsing err", http.StatusBadRequest)
			return
		}
		log.Println("req", req)
	}
}

func UserProfileHandler(ur UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
