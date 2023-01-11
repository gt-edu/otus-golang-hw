package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/dto"
	"github.com/jmoiron/sqlx"
)

type SQLStorage struct {
	DataSourceName string

	db  *sqlx.DB
	ctx context.Context
	log logger.Logger
}

func New(conf config.StorageConfig) *SQLStorage {
	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", conf.Username, conf.Password, conf.Hostname, conf.Port, conf.Dbname)
	return &SQLStorage{
		DataSourceName: dbURI,
	}
}

func (s *SQLStorage) Connect(ctx context.Context) error {
	db, err := sqlx.Open("pgx", s.DataSourceName) // *sqlx.DB
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

func (s *SQLStorage) Close(ctx context.Context) error {
	if s.db != nil {
		return s.db.Close()
	}

	s.db = nil
	return nil
}

func (s *SQLStorage) Add(e dto.Event) (int, error) {
	query := `insert into events(owner_id, title, descr, start_date, start_time, end_date, end_time, 
                                 notification_period) 
                          values(:owner_id, :title, :descr, :start_date, :start_time, :end_date, :end_time, 
                                 :notification_period) returning id`

	var evventID int
	rows, err := s.db.NamedQueryContext(s.ctx, query, e)
	if err != nil {
		return 0, err
	}

	if rows.Next() {
		err := rows.Scan(&evventID)
		if err != nil {
			return 0, err
		}
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}

	return evventID, nil
}

func (s *SQLStorage) Update(e dto.Event) error {
	_, err := s.Get(e.ID)
	if err != nil {
		return err
	}

	query := `update events set 
                  owner_id = :owner_id, 
                  title = :title, 
                  descr = :descr, 
                  start_date = :start_date, 
                  start_time = :start_time, 
                  end_date = :end_date, 
                  end_time = :end_time,
                  notification_period = :notification_period 
              where id = :id`
	_, err = s.db.NamedExecContext(s.ctx, query, e)
	if err != nil {
		return err
	}

	return nil
}

func (s *SQLStorage) Get(id int) (*dto.Event, error) {
	query := `select * from events where id = :id`
	rows, err := s.db.NamedQueryContext(s.ctx, query, map[string]interface{}{"id": id})
	if err != nil {
		return nil, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			s.log.Error("error during closing row set" + err.Error())
		}
	}(rows)

	if rows.Next() {
		var event dto.Event
		err := rows.StructScan(&event)
		if err != nil {
			return nil, err
		}
		return &event, nil
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nil, dto.ErrEventNotFound
}

func (s *SQLStorage) GetAll() ([]*dto.Event, error) {
	query := `select * from events`
	rows, err := s.db.NamedQueryContext(s.ctx, query, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	// ошибка при выполнении запроса
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			s.log.Error("error during closing row set" + err.Error())
		}
	}(rows)

	var events []*dto.Event
	for rows.Next() {
		var event dto.Event
		if err := rows.StructScan(&event); err != nil {
			return nil, err
		}

		events = append(events, &event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (s *SQLStorage) Delete(id int) error {
	query := `delete from events where id = $1`
	_, err := s.db.ExecContext(s.ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
