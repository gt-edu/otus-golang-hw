package app

import (
	"context"

	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  *logger.Logger
	storage storage.EventStorage
}

func New(logger *logger.Logger, storage storage.EventStorage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}
