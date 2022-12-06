package tests

import (
	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
	"testing"
)

func EventsCommonTest(t *testing.T, s storage.EventStorage) {
	t.Helper()

	testTitles := []string{"Event title 1", "Event title 2", "Event title 3", "Event title 4"}

	t.Run("add events", func(t *testing.T) {
		for ind, title := range testTitles {
			e := storage.Event{Title: title}
			id, err := s.Add(e)
			require.NoError(t, err)
			require.Equal(t, ind+1, id)
		}

		allEvents, err := s.GetAll()
		require.NoError(t, err)
		require.Len(t, allEvents, len(testTitles))
	})

	//t.Run("other tests", func(t *testing.T) {
	//	for ind, title := range testTitles {
	//		savedEvent, err := s.Get(ind + 1)
	//		require.NoError(t, err)
	//		require.NotNil(t, savedEvent)
	//		require.Equal(t, savedEvent.ID, ind+1)
	//		require.Equal(t, savedEvent.Title, title)
	//	}
	//
	//	err := s.Delete(1)
	//	require.NoError(t, err)
	//
	//	allEvents, err := s.GetAll()
	//	require.NoError(t, err)
	//	require.Len(t, allEvents, len(testTitles)-1)
	//
	//	nilEvent, err := s.Get(1)
	//	require.Nil(t, nilEvent)
	//	require.ErrorIs(t, err, storage.ErrEventNotFound)
	//
	//	_, err = s.Add(storage.Event{Title: "Event title 5"})
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
	//})

}
