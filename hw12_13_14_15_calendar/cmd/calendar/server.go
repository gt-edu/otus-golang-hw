package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
)

func startServer(configFile string) error {
	appConfig, err := config.NewConfig(configFile)
	if err != nil {
		return err
	}

	logg, err := logger.New(appConfig.Logger.Level, appConfig.Logger.Preset)
	if err != nil {
		return err
	}
	defer func(logg logger.Logger) {
		err := logg.Sync()
		if err != nil && !strings.Contains(err.Error(), "sync /dev/stderr: invalid argument") {
			_, _ = fmt.Fprintf(os.Stderr, "Error during logger syncing: %v", err)
		}
	}(logg)

	appStorage, err := storage.NewEventStorage(appConfig)
	if err != nil {
		return err
	}

	calendar := app.New(logg, appStorage)

	server := internalhttp.NewServer(logg, calendar, appConfig)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("Calendar server is running with config " + configFile + " ...")

	if err := server.Start(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
			return err
		} else {
			logg.Info("Server was stopped.")
		}
	}

	return nil
}
