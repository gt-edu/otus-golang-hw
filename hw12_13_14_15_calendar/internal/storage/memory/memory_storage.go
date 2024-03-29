package memorystorage

import (
	"sync"

	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/dto"
)

type MemoryStorage struct {
	mu        sync.RWMutex
	eventsMap map[int]*dto.Event
	lastID    int
}

func New() *MemoryStorage {
	return &MemoryStorage{
		eventsMap: make(map[int]*dto.Event),
	}
}

func (s *MemoryStorage) Add(e dto.Event) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lastID++
	e.ID = s.lastID
	s.eventsMap[e.ID] = &e

	return e.ID, nil
}

func (s *MemoryStorage) Update(e dto.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.getWithoutLock(e.ID)
	if err != nil {
		return err
	}

	s.eventsMap[e.ID] = &e

	return nil
}

func (s *MemoryStorage) Get(id int) (*dto.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.getWithoutLock(id)
}

func (s *MemoryStorage) getWithoutLock(id int) (*dto.Event, error) {
	e, ok := s.eventsMap[id]
	if !ok {
		return nil, dto.ErrEventNotFound
	}
	return e, nil
}

func (s *MemoryStorage) GetAll() ([]*dto.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	eventsList := make([]*dto.Event, len(s.eventsMap))
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
