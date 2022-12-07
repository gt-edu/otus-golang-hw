package memorystorage

import (
	"sync"

	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
)

type MemoryStorage struct {
	mu        sync.RWMutex
	eventsMap map[int]*storage.Event
	lastID    int
}

func New() *MemoryStorage {
	return &MemoryStorage{
		eventsMap: make(map[int]*storage.Event),
	}
}

func (s *MemoryStorage) Add(e storage.Event) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lastID++
	e.ID = s.lastID
	s.eventsMap[e.ID] = &e

	return e.ID, nil
}

func (s *MemoryStorage) Update(e storage.Event) error {
	_, err := s.Get(e.ID)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.eventsMap[e.ID] = &e

	return nil
}

func (s *MemoryStorage) Get(id int) (*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	e, ok := s.eventsMap[id]
	if !ok {
		return nil, storage.ErrEventNotFound
	}

	return e, nil
}

func (s *MemoryStorage) GetAll() ([]*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	eventsList := make([]*storage.Event, len(s.eventsMap))
	ind := 0
	for _, evt := range s.eventsMap {
		eventsList[ind] = evt
		ind++
	}
	return eventsList, nil
}

func (s *MemoryStorage) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.eventsMap, id)

	return nil
}
