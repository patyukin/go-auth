package auth

import (
	"context"
	"log/slog"
	"net"
	"net/http"

	authpb "github.com/bogatyr285/auth-go/pkg/server/grpc/auth"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

// Gateway - proxy which tranforms HTTP requests to GRPC
type Gateway struct {
	mux        http.Handler
	httpGwAddr string
	logger     *slog.Logger
}

func NewGateway(ctx context.Context, grpcAddr, httpGwAddr string, logger *slog.Logger) (*Gateway, error) {
	gwMux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard,
		&runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
		}),
	)

	err := authpb.RegisterAuthServiceHandlerFromEndpoint(
		context.Background(),
		gwMux,
		grpcAddr,
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			// grpc.WithStatsHandler(new(ocgrpc.ClientHandler)),
		},
	)
	if err != nil {
		return nil, err
	}

	return &Gateway{
		mux:        gwMux,
		httpGwAddr: httpGwAddr,
		logger:     logger.With("module", "grpc/http-gateway"),
	}, nil
}

func (g *Gateway) Start() (func() error, error) {
	hserver := http.Server{
		Handler: g.mux,
	}

	g.logger.Info("starting", slog.String("addr", g.httpGwAddr))
	l, err := net.Listen("tcp", g.httpGwAddr)
	if err != nil {
		return nil, err
	}

	go func() {
		err = hserver.Serve(l)
		if err != nil {
			g.logger.Error("http/grpc gateway server", slog.Any("err", err))
		}
	}()

	return func() error {
		g.logger.Info("shutting down")
		return hserver.Close()
	}, nil
}

type Closer func() error
