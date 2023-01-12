package tests

import (
	"context"
	"testing"

	sqlstorage "github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/sql"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestSqlStorage(t *testing.T) {
	factory := func() storage.EventStorage {
		t.Helper()
		s := setupTestdbAndRunMigrations(t)
		require.NotNil(t, s)
		return s
	}

	EventsCommonTest(t, factory)

	EventsCommonConcurrencyTest(t, factory)
}

func setupTestdbAndRunMigrations(t *testing.T) *sqlstorage.SQLStorage {
	t.Helper()

	_, storageConf := SetupTestcontainersDatabase(t)
	eventStorage := sqlstorage.New(storageConf)

	err := sqlstorage.RunMigrationsUp(eventStorage)
	require.NoError(t, err)

	return eventStorage
}

func SetupTestcontainersDatabase(t *testing.T) (*testcontainers.Container, config.StorageConfig) {
	t.Helper()

	storageConf := config.StorageConfig{
		Type:     "sql",
		Dbname:   "testdb",
		Username: "postgres",
		Password: "postgres",
	}

	// 1. Create PostgreSQL container request
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:15.1",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       storageConf.Dbname,
			"POSTGRES_PASSWORD": storageConf.Password,
			"POSTGRES_USER":     storageConf.Username,
		},
		Labels: map[string]string{
			"app-calendar": "postgresql",
		},
	}

	// 2. Start PostgreSQL container
	dbContainer, err := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})
	require.NoError(t, err)

	// 3.1 Get host and port of PostgreSQL container
	host, err := dbContainer.Host(context.Background())
	require.NoError(t, err)
	port, err := dbContainer.MappedPort(context.Background(), "5432")
	require.NoError(t, err)

	storageConf.Port = port.Port()
	storageConf.Hostname = host

	return &dbContainer, storageConf
}
