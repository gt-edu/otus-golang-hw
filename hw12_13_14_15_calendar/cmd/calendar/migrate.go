package main

import (
	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	sqlstorage "github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

func migrateAppStorage(configFile string) error {
	appConfig, err := config.NewConfig(configFile)
	if err != nil {
		return err
	}

	logg, err := logger.New(appConfig.Logger.Level, appConfig.Logger.Preset)
	if err != nil {
		return err
	}
	defer logger.SafeLoggerSync(logg)

	appStorage := sqlstorage.New(appConfig.Storage)
	logg.Info("Run migrations ...")
	return sqlstorage.RunMigrationsUp(appStorage)
}
