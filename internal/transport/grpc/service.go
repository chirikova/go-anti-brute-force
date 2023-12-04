package grpc

import (
	"context"
	"errors"
	"net"

	"github.com/chirikova/go-anti-brute-force/internal/app"
	"github.com/chirikova/go-anti-brute-force/internal/storage"
	"github.com/chirikova/go-anti-brute-force/internal/transport"
	"github.com/chirikova/go-anti-brute-force/internal/transport/grpc/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	app app.Application
	api.UnimplementedApiServiceServer
}

func NewService(app app.Application) *Service {
	return &Service{
		app: app,
	}
}

func (s *Service) Auth(_ context.Context, r *api.AuthRequest) (*api.AuthResponse, error) {
	ipAddress := net.ParseIP(r.Ip)
	if ipAddress.IsUnspecified() {
		return &api.AuthResponse{Ok: false},
			status.Errorf(codes.InvalidArgument, transport.ErrInvalidIP.Error())
	}

	ok, err := s.app.Verify(r.Login, r.Password, &ipAddress)

	return &api.AuthResponse{Ok: ok}, err
}

func (s *Service) Reset(_ context.Context, r *api.ResetRequest) (*api.ResetResponse, error) {
	ipAddress := net.ParseIP(r.Ip)
	if ipAddress.IsUnspecified() {
		return &api.ResetResponse{Ok: false},
			status.Errorf(codes.InvalidArgument, transport.ErrInvalidIP.Error())
	}

	ok := s.app.Reset(r.Login, &ipAddress)

	return &api.ResetResponse{Ok: ok}, nil
}

func (s *Service) WhitelistAdd(_ context.Context, r *api.WhitelistAddRequest) (*api.WhitelistAddResponse, error) {
	response := &api.WhitelistAddResponse{}

	ipNet, err := toIPNet(r.SubNet)
	if err != nil {
		return response, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = s.app.AddToWhiteList(ipNet)

	if errors.Is(err, storage.ErrAlreadyExist) {
		return response, status.Errorf(codes.AlreadyExists, err.Error())
	}

	return response, nil
}

func (s *Service) WhitelistRemove(_ context.Context, r *api.WhitelistRemoveRequest) (*api.WhitelistRemoveResponse, error) { //nolint:all
	response := &api.WhitelistRemoveResponse{}

	ipNet, err := toIPNet(r.SubNet)
	if err != nil {
		return response, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = s.app.RemoveFromWhiteList(ipNet)

	if errors.Is(err, storage.ErrNotFound) {
		return response, status.Errorf(codes.NotFound, err.Error())
	}

	return response, nil
}

func (s *Service) BlacklistAdd(_ context.Context, r *api.BlacklistAddRequest) (*api.BlacklistAddResponse, error) {
	response := &api.BlacklistAddResponse{}

	ipNet, err := toIPNet(r.SubNet)
	if err != nil {
		return response, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = s.app.AddToBlackList(ipNet)

	if errors.Is(err, storage.ErrAlreadyExist) {
		return response, status.Errorf(codes.AlreadyExists, err.Error())
	}

	return response, nil
}

func (s *Service) BlacklistRemove(_ context.Context, r *api.BlacklistRemoveRequest) (*api.BlacklistRemoveResponse, error) { //nolint:all
	response := &api.BlacklistRemoveResponse{}

	ipNet, err := toIPNet(r.SubNet)
	if err != nil {
		return response, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = s.app.RemoveFromBlackList(ipNet)

	if errors.Is(err, storage.ErrNotFound) {
		return response, status.Errorf(codes.NotFound, err.Error())
	}

	return response, nil
}

func toIPNet(subNet *api.SubNet) (*net.IPNet, error) {
	_, ipNet, err := net.ParseCIDR(subNet.Ip + "/" + subNet.GetMask())

	return ipNet, err
}
