package sqlstorage

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

//go:embed migrations
var migrations embed.FS

func TestStrorage_TestContainers(t *testing.T) {
	_, db := SetupTestDatabase(t)
	var err error
	for i := 0; i < 10; i++ {
		err := db.PingContext(context.Background())
		if err != nil {
			time.Sleep(time.Duration(500) * time.Millisecond)
		} else {
			break
		}
	}

	require.NoError(t, err)

	driver, err := pgx.WithInstance(db, &pgx.Config{})
	require.NoError(t, err)

	source, err := iofs.New(migrations, "migrations")
	require.NoError(t, err)

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	require.NoError(t, err)
	m.Log = &MigrateLog{}

	err = m.Up()
	require.NoError(t, err)
}

// func TestStorage_Postgresql(t *testing.T) {
//	testTitles := []string{"Event title 1", "Event title 2", "Event title 3", "Event title 4"}
//
//	s := memorystorage.New()
//	for _, title := range testTitles {
//		e := storage.Event{Title: title}
//		err := s.Add(e)
//		require.NoError(t, err)
//	}
//
//	allEvents, err := s.GetAll()
//	require.NoError(t, err)
//	require.Len(t, allEvents, len(testTitles))
//
//	for ind, title := range testTitles {
//		savedEvent, err := s.Get(ind + 1)
//		require.NoError(t, err)
//		require.NotNil(t, savedEvent)
//		require.Equal(t, savedEvent.ID, ind+1)
//		require.Equal(t, savedEvent.Title, title)
//	}
//
//	err = s.Delete(1)
//	require.NoError(t, err)
//
//	allEvents, err = s.GetAll()
//	require.NoError(t, err)
//	require.Len(t, allEvents, len(testTitles)-1)
//
//	nilEvent, err := s.Get(1)
//	require.Nil(t, nilEvent)
//	require.ErrorIs(t, err, storage.ErrEventNotFound)
//
//	err = s.Add(storage.Event{Title: "Event title 5"})
//	require.NoError(t, err)
//
//	allEvents, err = s.GetAll()
//	require.NoError(t, err)
//	require.Len(t, allEvents, len(testTitles))
//
//	event, err := s.Get(5)
//	require.NoError(t, err)
//	require.NotNil(t, event)
//	require.Equal(t, 5, event.ID)
//
//	err = s.Update(storage.Event{ID: 5, Title: "Event title 55"})
//	require.NoError(t, err)
//
//	event, err = s.Get(5)
//	require.NoError(t, err)
//	require.NotNil(t, event)
//	require.Equal(t, 5, event.ID)
//	require.Equal(t, "Event title 55", event.Title)
//
//	err = s.Update(storage.Event{ID: 6, Title: "Event title 66"})
//	require.ErrorIs(t, err, storage.ErrEventNotFound)
//}

func SetupTestDatabase(t *testing.T) (*testcontainers.Container, *sql.DB) {
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
	db, err := sql.Open("pgx", dbURI)
	require.NoError(t, err)

	return &dbContainer, db
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
