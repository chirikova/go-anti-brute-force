package grpc

import (
	"context"

	"github.com/chirikova/go-anti-brute-force/internal/app"
	"github.com/chirikova/go-anti-brute-force/internal/transport/grpc/api"
)

type Service struct {
	app *app.Application
	api.UnimplementedApiServiceServer
}

func NewService(app *app.Application) *Service {
	return &Service{
		app: app,
	}
}

func (s *Service) Auth(_ context.Context, _ *api.AuthRequest) (*api.AuthResponse, error) {
	return &api.AuthResponse{}, nil
}

func (s *Service) Reset(_ context.Context, _ *api.ResetRequest) (*api.ResetResponse, error) {
	return &api.ResetResponse{}, nil
}

func (s *Service) WhitelistAdd(_ context.Context, _ *api.WhitelistAddRequest) (*api.WhitelistAddResponse, error) {
	return &api.WhitelistAddResponse{}, nil
}

func (s *Service) WhitelistRemove(ctx context.Context, request *api.WhitelistRemoveRequest) (*api.WhitelistRemoveResponse, error) { //nolint:all
	return &api.WhitelistRemoveResponse{}, nil
}

func (s *Service) BlacklistAdd(_ context.Context, _ *api.BlacklistAddRequest) (*api.BlacklistAddResponse, error) {
	return &api.BlacklistAddResponse{}, nil
}

func (s *Service) BlacklistRemove(_ context.Context, _ *api.BlacklistRemoveRequest) (*api.BlacklistRemoveResponse, error) { //nolint:all
	return &api.BlacklistRemoveResponse{}, nil
}
