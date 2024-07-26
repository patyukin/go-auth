package commands

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/bogatyr285/auth-go/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/cobra"
)

type Response struct {
	Error string
	Data  interface{}
}

func NewServeCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "serve",
		Aliases: []string{"s"},
		Short:   "Start API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := signal.NotifyContext(cmd.Context(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
			defer cancel()

			router := chi.NewRouter()
			router.Use(middleware.RequestID)
			router.Use(middleware.Recoverer)
			router.Use(middleware.Logger)

			s, err := storage.New("./users.sql")
			if err != nil {
				return err
			}

			router.Post("/register", RegisterHandler(&s))
			router.Post("/login", LoginHandler(&s))

			httpServer := http.Server{
				Addr:         "localhost:8080",
				ReadTimeout:  time.Second,
				WriteTimeout: time.Second,
				Handler:      router,
			}
			// TODO change logger to slog
			go func() {
				if err := httpServer.ListenAndServe(); err != nil {
					log.Println("ListenAndServe", err)
				}
			}()
			log.Println("server listening: 8080")
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
	return c
}

type UserRepository interface {
	RegisterUser(ctx context.Context, u storage.UserAccount) (storage.UserAccount, error)
	Login(ctx context.Context, username, password string) (storage.UserAccount, error)
}

func RegisterHandler(ur UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
		// parse body
	}
}

func LoginHandler(ur UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
