package grpc

import (
	"context"

	"antibruteforce/internal/transport/grpc/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(ctx context.Context, address string) (api.ApiServiceClient, error) {
	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := api.NewApiServiceClient(conn)

	return client, nil
}
