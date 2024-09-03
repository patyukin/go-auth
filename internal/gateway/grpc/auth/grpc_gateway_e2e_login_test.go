package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"testing"

	"github.com/bogatyr285/auth-go/config"
	"github.com/bogatyr285/auth-go/internal/auth/repository"
	"github.com/bogatyr285/auth-go/internal/buildinfo"
	"github.com/bogatyr285/auth-go/internal/gateway/grpc/auth"
	"github.com/bogatyr285/auth-go/pkg/crypto"
	"github.com/bogatyr285/auth-go/pkg/jwt"
	authpb "github.com/bogatyr285/auth-go/pkg/server/grpc/auth"
	"github.com/bogatyr285/auth-go/playground"
	"github.com/stretchr/testify/assert"
)

func TestGRPCGatewayLogin(t *testing.T) {
	// configure
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	cfg, err := config.Parse("../../../../config.yaml")
	assert.NoError(t, err)

	storage, err := repository.New(cfg.Storage.SQLitePath)
	assert.NoError(t, err)

	passwordHasher := crypto.NewPasswordHasher()
	jwtManager, err := jwt.NewJWTManager(
		cfg.JWT.Issuer,
		cfg.JWT.ExpiresIn,
		[]byte(cfg.JWT.PublicKey),
		[]byte(cfg.JWT.PrivateKey))
	assert.NoError(t, err)

	grpcAddress := ":9090"
	httpGwAddress := ":9091"
	authGRPCHandlers := auth.NewAuthHandlers(&storage, passwordHasher, jwtManager, buildinfo.New())
	grpcServer, err := auth.NewGRPCServer(grpcAddress, authGRPCHandlers, log)
	assert.NoError(t, err)

	grpcCloser, err := grpcServer.Run()
	assert.NoError(t, err)
	defer grpcCloser()

	grpcGw, err := auth.NewGateway(context.Background(), grpcAddress, httpGwAddress, log)
	assert.NoError(t, err)

	grpcGwCloser, err := grpcGw.Start()
	assert.NoError(t, err)
	defer grpcGwCloser()

	// prepare
	loginUserReq := &authpb.LoginUserRequest{
		LoginMethod: &authpb.LoginUserRequest_Email{Email: "user1"},
		Password:    "rLy_5tr0nG!",
	}

	loginUserReqBytes, _ := playground.ProtobufToJSON(loginUserReq)

	// act
	res, err := http.Post(fmt.Sprintf("http://localhost%s/api/v1/login", httpGwAddress),
		"application/json",
		bytes.NewReader(loginUserReqBytes))
	assert.NoError(t, err)

	assert.Equal(t, res.StatusCode, http.StatusOK)

	// check
	loginResModel := &authpb.LoginUserResponse{}
	bodyRes, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	json.Unmarshal(bodyRes, loginResModel)

	assert.NotEmpty(t, loginResModel.Token)
}
