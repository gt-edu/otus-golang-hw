package main

import (
	"context"
	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/memory"
	"os/signal"
	"syscall"
	"time"
)

func startServer(configFile string) error {

	config, err := NewConfig(configFile)
	if err != nil {
		return err
	}

	logg := logger.New(config.Logger.Level)

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar)

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

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		return err
	}

	return nil
}
