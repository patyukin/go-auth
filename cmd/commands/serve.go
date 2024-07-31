package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bogatyr285/auth-go/config"
	"github.com/bogatyr285/auth-go/internal/buildinfo"
	"github.com/bogatyr285/auth-go/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
)

func NewServeCmd() *cobra.Command {
	var configPath string

	c := &cobra.Command{
		Use:     "serve",
		Aliases: []string{"s"},
		Short:   "Start API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

			ctx, cancel := signal.NotifyContext(cmd.Context(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
			defer cancel()

			router := chi.NewRouter()
			router.Use(middleware.Logger)
			router.Use(middleware.RequestID)
			router.Use(middleware.Recoverer)

			s, err := storage.New("./users.sql")
			if err != nil {
				return err
			}
			cfg, err := config.Parse(configPath)
			if err != nil {
				return err
			}

			slog.Info("loaded cfg", slog.Any("cfg", cfg))

			router.Post("/register", RegisterHandler(&s))
			router.Post("/login", LoginHandler(&s))
			router.Get("/build", buildinfo.BuildInfoHandler(buildinfo.New()))

			httpServer := http.Server{
				Addr:         cfg.HTTPServer.Address,
				ReadTimeout:  cfg.HTTPServer.Timeout,
				WriteTimeout: cfg.HTTPServer.Timeout,
				Handler:      router,
			}

			go func() {
				if err := httpServer.ListenAndServe(); err != nil {
					log.Error("ListenAndServe", slog.Any("err", err))
				}
			}()
			log.Info("server listening: 8080")
			<-ctx.Done()

			closeCtx, _ := context.WithTimeout(context.Background(), time.Second*5)
			if err := httpServer.Shutdown(closeCtx); err != nil {
				return fmt.Errorf("http closing err: %s", err)
			}

			// close db connection
			// etc

			return nil
		},
	}
	c.Flags().StringVar(&configPath, "config", "", "path to config")
	return c
}

type UserRepository interface {
	RegisterUser(ctx context.Context, u storage.UserAccount) error
	Login(ctx context.Context, username, password string) (storage.UserAccount, error)
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=1"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=1"`
}

type Response struct {
	Error string
	Data  interface{}
}

func RegisterHandler(ur UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &RegisterRequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			http.Error(w, "parsing err", http.StatusBadRequest)
			// render.JSON(w, r, Response{
			// 	Error: ,
			// })
			return
		}

		err := ur.RegisterUser(r.Context(), storage.UserAccount{
			Username: req.Email,
			Password: req.Password,
		})
		if err != nil {
			http.Error(w, "reg", http.StatusBadRequest)
			return
		}

		if err = validator.New().Struct(req); err != nil {
			http.Error(w, "validation err", http.StatusBadRequest)
			return
		}
	}
}

func LoginHandler(ur UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &LoginRequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			http.Error(w, "parsing err", http.StatusBadRequest)
			return
		}

		// handle different errors
		account, err := ur.Login(r.Context(), req.Email, req.Password)
		if err != nil {
			http.Error(w, "login err", http.StatusUnauthorized)
			return
		}

		// JWT
		log.Println("account", account)
	}
}
