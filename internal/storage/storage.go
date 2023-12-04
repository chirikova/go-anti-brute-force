package storage

import (
	"context"

	"github.com/chirikova/go-anti-brute-force/internal/config"
)

type Connector interface {
	Connect(ctx context.Context, db config.DB) error
	Close() error
}
