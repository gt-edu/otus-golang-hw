package storage

import "github.com/pkg/errors"

var ErrEventNotFound = errors.New("event not found")

type EventStorage interface {
	Add(e Event) (int, error)
	Update(e Event) error
	Get(id int) (*Event, error)
	GetAll() ([]*Event, error)
	Delete(id int) error
}
