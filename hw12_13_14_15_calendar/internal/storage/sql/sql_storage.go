package sqlstorage

import (
	"context"
	"database/sql"
	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	"time"
)

type SqlStorage struct {
	DataSourceName string

	db  *sql.DB
	ctx context.Context
	log *logger.Logger
}

func New() *SqlStorage {
	return &SqlStorage{}
}

func (s *SqlStorage) Connect(ctx context.Context) error {
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

func (s *SqlStorage) Close(ctx context.Context) error {
	if s.db != nil {
		return s.db.Close()
	}

	s.db = nil
	return nil
}

func (s *SqlStorage) Add(e storage.Event) (int, error) {
	query := `insert into events(title, start_date, end_date) values($1, $2, $3) returning id`
	var eventId int
	err := s.db.QueryRowContext(s.ctx, query, e.Title, "2019-12-31", "2019-12-31").Scan(&eventId)
	if err != nil {
		return 0, err
	}

	return eventId, nil
}

func (s *SqlStorage) Update(e storage.Event) error {
	query := `update events set title = $1 where id = $2`
	_, err := s.db.ExecContext(s.ctx, query, e.Title, e.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *SqlStorage) Get(id int) (*storage.Event, error) {
	query := `select title from events where id = $1`
	row := s.db.QueryRowContext(s.ctx, query, id)
	var title string
	err := row.Scan(&title)
	if err == sql.ErrNoRows {
		return nil, storage.ErrEventNotFound
	} else if err != nil {
		return nil, err
	}
	return &storage.Event{ID: id, Title: title}, nil
}

func (s *SqlStorage) GetAll() ([]*storage.Event, error) {
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

func (s *SqlStorage) Delete(id int) error {
	query := `delete from events where id = $1`
	_, err := s.db.ExecContext(s.ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
