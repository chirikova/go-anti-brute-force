package app

import (
	"context"

	"antibruteforce/internal/logger"
)

type Application interface {
	Verify() bool
}

type App struct {
	logger logger.Logger
}

func NewApp(_ context.Context, logger logger.Logger) Application {
	return App{
		logger: logger,
	}
}

func (a App) Verify() bool {
	// TODO
	return false
}
