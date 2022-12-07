package tests

import (
	"context"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

func TestSqlStorage(t *testing.T) {
	s := setupTestdbAndRunMigrations(t)
	require.NotNil(t, s)

	EventsCommonTest(t, s)
}

func setupTestdbAndRunMigrations(t *testing.T) *sqlstorage.SqlStorage {
	_, dbURI := SetupTestcontainersDatabase(t)
	eventStorage := sqlstorage.New()
	eventStorage.DataSourceName = dbURI

	err := sqlstorage.RunMigrationsUp(t, eventStorage)
	require.NoError(t, err)

	return eventStorage
}

func SetupTestcontainersDatabase(t *testing.T) (*testcontainers.Container, string) {
	t.Helper()

	// 1. Create PostgreSQL container request
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:15.1",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_USER":     "postgres",
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

	// 3.2 Create db connection string and connect
	dbURI := fmt.Sprintf("postgres://postgres:postgres@%v:%v/testdb", host, port.Port())

	return &dbContainer, dbURI
}
