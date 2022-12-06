package sqlstorage

import (
	"context"
	"database/sql"
	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	"time"
)

type Storage struct {
	DataSourceName string

	db  *sql.DB
	ctx context.Context
	log *logger.Logger
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sql.Open("pgx", s.DataSourceName)
	if err != nil {
		return err
	}

	s.db = db
	s.ctx = ctx

	for i := 0; i < 10; i++ {
		err = s.db.PingContext(s.ctx)
		if err != nil {
			time.Sleep(time.Duration(500) * time.Millisecond)
		} else {
			break
		}
	}
	return err
}

func (s *Storage) Close(ctx context.Context) error {
	if s.db != nil {
		return s.db.Close()
	}

	s.db = nil
	return nil
}

func (s *Storage) Add(e storage.Event) (int, error) {
	query := `insert into events(title, start_date, end_date) values($1, $2, $3) returning id`
	var eventId int
	err := s.db.QueryRowContext(s.ctx, query, e.Title, "2019-12-31", "2019-12-31").Scan(&eventId)
	if err != nil {
		return 0, err
	}

	return eventId, nil
}

func (s *Storage) Update(e storage.Event) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) Get(id int) (*storage.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetAll() ([]*storage.Event, error) {
	query := `select id, title from events`
	rows, err := s.db.QueryContext(s.ctx, query)
	if err != nil {
		return nil, err
	}
	// ошибка при выполнении запроса
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			s.log.Error("error during closing row set" + err.Error())
		}
	}(rows)

	var events []*storage.Event
	for rows.Next() {
		var id int
		var title string
		if err := rows.Scan(&id, &title); err != nil {
			return nil, err
		}

		events = append(events, &storage.Event{ID: id, Title: title})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Storage) Delete(id int) error {
	//TODO implement me
	panic("implement me")
}
