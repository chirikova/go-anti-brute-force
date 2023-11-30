package sqlstorage

import (
	"context"
	"fmt"

	"github.com/chirikova/go-anti-brute-force/internal/config"
	"github.com/chirikova/go-anti-brute-force/internal/storage"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	ctx context.Context
	db  *sqlx.DB
}

func New(ctx context.Context) (*Storage, error) {
	return &Storage{ctx: ctx}, nil
}

func (s *Storage) Connect(ctx context.Context, cfg config.DB) error {
	db, err := sqlx.ConnectContext(ctx, "pgx", cfg.BuildDSN())
	if err != nil {
		err = fmt.Errorf("%w: %w", storage.ErrDBConnect, err)
	}

	s.db = db

	return err
}

func (s *Storage) Close() error {
	err := s.db.Close()
	if err != nil {
		err = fmt.Errorf("%w: %w", storage.ErrDBClose, err)
	}

	return err
}
