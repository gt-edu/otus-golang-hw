package tests

import (
	"github.com/stretchr/testify/require"
	"strconv"
	"sync"
	"testing"

	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/dto"
)

func EventsCommonConcurrencyTest(t *testing.T, storageFactory func() storage.EventStorage) {
	t.Helper()

	testEvents := []dto.Event{
		{
			Title: "Event title 1", OwnerID: 2, Descr: "Description 1", StartDate: "2022-12-06T00:00:00Z", StartTime: "17:00:00",
			EndDate: "2022-12-06T00:00:00Z", EndTime: "18:00:00", NotificationPeriod: "1",
		},
	}

	t.Run("test concurrent adding and changing", func(t *testing.T) {
		s := storageFactory()
		n := 10
		wg := &sync.WaitGroup{}
		for i := 0; i < n; i++ {
			wg.Add(1)

			go func(i int) {
				defer wg.Done()
				eventId, err := s.Add(testEvents[0])
				require.NoError(t, err)

				chDone := make(chan struct{})
				go func(eventId int) {
					changedEvent, err := s.Get(eventId)
					require.NoError(t, err)
					changedEvent.Title = "Event title M " + strconv.Itoa(changedEvent.ID)

					err = s.Update(*changedEvent)
					require.NoError(t, err)

					chDone <- struct{}{}
				}(eventId)
				<-chDone
			}(i)
		}
		wg.Wait()

		allEvents, err := s.GetAll()
		require.NoError(t, err)
		require.Equal(t, n, len(allEvents))

		for _, savedEvent := range allEvents {
			require.Equal(t, "Event title M "+strconv.Itoa(savedEvent.ID), savedEvent.Title)
		}
	})

	t.Run("test concurrent adding and deleting", func(t *testing.T) {
		s := storageFactory()
		n := 10
		wg := &sync.WaitGroup{}
		for i := 0; i < n; i++ {
			wg.Add(1)

			go func(i int) {
				defer wg.Done()
				eventId, err := s.Add(testEvents[0])
				require.NoError(t, err)

				chDone := make(chan struct{})
				go func(eventId int) {
					err := s.Delete(eventId)
					require.NoError(t, err)
					chDone <- struct{}{}
				}(eventId)
				<-chDone
			}(i)
		}
		wg.Wait()

		allEvents, err := s.GetAll()
		require.NoError(t, err)
		require.Equal(t, 0, len(allEvents))
	})

}
