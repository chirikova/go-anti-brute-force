package grpc

import (
	"antibruteforce/internal/app"
	"context"

	"antibruteforce/internal/transport/grpc/api"
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

func (s *Service) Auth(ctx context.Context, request *api.AuthRequest) (*api.AuthResponse, error) {
	return &api.AuthResponse{}, nil
}

func (s *Service) Reset(ctx context.Context, request *api.ResetRequest) (*api.ResetResponse, error) {

	return &api.ResetResponse{}, nil
}

func (s *Service) WhitelistAdd(ctx context.Context, request *api.WhitelistAddRequest) (*api.WhitelistAddResponse, error) {

	return &api.WhitelistAddResponse{}, nil
}

func (s *Service) WhitelistRemove(ctx context.Context, request *api.WhitelistRemoveRequest) (*api.WhitelistRemoveResponse, error) {

	return &api.WhitelistRemoveResponse{}, nil
}

func (s *Service) BlacklistAdd(ctx context.Context, request *api.BlacklistAddRequest) (*api.BlacklistAddResponse, error) {

	return &api.BlacklistAddResponse{}, nil
}

func (s *Service) BlacklistRemove(ctx context.Context, request *api.BlacklistRemoveRequest) (*api.BlacklistRemoveResponse, error) {

	return &api.BlacklistRemoveResponse{}, nil
}
