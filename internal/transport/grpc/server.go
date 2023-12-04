package grpc

import (
	"context"
	"net"

	"github.com/chirikova/go-anti-brute-force/internal/app"
	"github.com/chirikova/go-anti-brute-force/internal/config"
	"github.com/chirikova/go-anti-brute-force/internal/logger"
	"github.com/chirikova/go-anti-brute-force/internal/transport/grpc/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	config config.GRPC
	logger logger.Logger
	server *grpc.Server
	ctx    context.Context
}

func NewServer(ctx context.Context, config config.GRPC, logger logger.Logger, app app.Application) *Server {
	apiService := NewService(app)
	server := grpc.NewServer(grpc.UnaryInterceptor(loggingMiddleware(logger)))
	api.RegisterApiServiceServer(server, apiService)
	reflection.Register(server)

	return &Server{
		config: config,
		logger: logger,
		server: server,
		ctx:    ctx,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", net.JoinHostPort(s.config.Host, s.config.Port))
	if err != nil {
		return err
	}

	if err := s.server.Serve(listener); err != nil {
		return err
	}
	s.logger.Info("starting grpc server")

	return nil
}

func (s *Server) Stop() error {
	s.server.GracefulStop()

	s.logger.Info("stopping grpc server")

	return nil
}
