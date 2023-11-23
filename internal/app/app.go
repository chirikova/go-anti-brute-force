package app

import (
	"antibruteforce/internal/logger"
	"context"
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
