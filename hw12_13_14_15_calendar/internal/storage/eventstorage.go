package storage

import (
	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/dto"
	memorystorage "github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

type EventStorage interface {
	Add(e dto.Event) (int, error)
	Update(e dto.Event) error
	Get(id int) (*dto.Event, error)
	GetAll() ([]*dto.Event, error)
	Delete(id int) error
}

func NewEventStorage(appConfig *config.Config, logg logger.Logger) (EventStorage, error) {
	switch appConfig.Storage.Type {
	case "memory":
		logg.Info("Using memory storage")
		return memorystorage.New(), nil
	case "sql":
		eventStorage := sqlstorage.New(appConfig.Storage)
		return eventStorage, nil
	default:
		return nil, dto.ErrStorageTypeIsNotCorrect
	}
}
