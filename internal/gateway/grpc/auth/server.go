package auth

import (
	"log/slog"
	"net"
	"time"

	authpb "github.com/bogatyr285/auth-go/pkg/server/grpc/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const (
	// GRPCDefaultGracefulStopTimeout - period to wait result of grpc.GracefulStop
	// after call grpc.Stop
	GRPCDefaultGracefulStopTimeout = 5 * time.Second
)

// GRPC - structure describes gRPC props
type Server struct {
	grpcAddr            string
	grpcSrv             *grpc.Server
	listener            net.Listener
	gracefulStopTimeout time.Duration

	logger *slog.Logger
}

func NewGRPCServer(
	grpcAddr string,
	authHadndlers authpb.AuthServiceServer,
	logger *slog.Logger,
) (*Server, error) {
	logger = logger.With("module", "grpc-server")
	netListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return nil, err
	}
	opts := []grpc.ServerOption{}

	grpcSrv := grpc.NewServer(opts...)
	authpb.RegisterAuthServiceServer(grpcSrv, authHadndlers)

	// register health check service
	healthService := NewHealthChecker(logger)
	grpc_health_v1.RegisterHealthServer(grpcSrv, healthService)

	// Register reflection service on gRPC server. can be a flag
	reflection.Register(grpcSrv)

	server := &Server{
		grpcAddr:            grpcAddr,
		listener:            netListener,
		grpcSrv:             grpcSrv,
		gracefulStopTimeout: GRPCDefaultGracefulStopTimeout,
		logger:              logger,
	}

	return server, nil
}

func (s *Server) Run() (func() error, error) {
	s.logger.Info("starting", slog.String("grpcAddr", s.grpcAddr))

	go func() {
		err := s.grpcSrv.Serve(s.listener)
		if err == grpc.ErrServerStopped {
			s.logger.Error("grpc server", slog.Any("err", err))
		}
	}()

	return s.close, nil
}

// stop - gracefully stop server & listeners
func (s *Server) close() error {
	s.logger.Info("gracefully stopping....", slog.String("grpcAddr", s.grpcAddr))

	stopped := make(chan struct{})
	go func() {
		s.grpcSrv.GracefulStop()
		close(stopped)
	}()

	t := time.NewTimer(s.gracefulStopTimeout)
	defer t.Stop()

	select {
	case <-t.C:
		s.logger.Info("ungracefully stopping....", slog.String("grpcAddr", s.grpcAddr))
		s.grpcSrv.Stop()
	case <-stopped:
		t.Stop()
	}
	s.logger.Info("stopped", slog.String("grpcAddr", s.grpcAddr))
	return nil
}
