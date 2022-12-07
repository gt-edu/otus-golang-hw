package sqlstorage

import (
	"context"
	"embed"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"log"
	"testing"
)

//go:embed migrations
var migrations embed.FS

func RunMigrationsUp(t *testing.T, eventStorage *SqlStorage) error {
	err := eventStorage.Connect(context.Background())
	if err != nil {
		return err
	}

	driver, err := pgx.WithInstance(eventStorage.db, &pgx.Config{})
	if err != nil {
		return err
	}

	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return err
	}

	m.Log = &MigrateLog{}

	err = m.Up()
	if err != nil {
		return err
	}

	return nil
}

// MigrateLog represents the logger.
type MigrateLog struct{}

// Printf prints out formatted string into a log.
func (l *MigrateLog) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// Println prints out args into a log.
func (l *MigrateLog) Println(args ...interface{}) {
	log.Println(args...)
}

// Verbose shows if verbose print enabled.
func (l *MigrateLog) Verbose() bool {
	return true
}
