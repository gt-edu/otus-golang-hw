package tests

import (
	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/dto"
	"testing"

	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func EventsCommonTest(t *testing.T, s storage.EventStorage) {
	t.Helper()

	testEvents := []dto.Event{
		{
			Title: "Event title 1", OwnerID: 2, Descr: "Description 1", StartDate: "2022-12-06T00:00:00Z", StartTime: "17:00:00",
			EndDate: "2022-12-06T00:00:00Z", EndTime: "18:00:00", NotificationPeriod: "1",
		},
		{
			Title: "Event title 2", OwnerID: 2, Descr: "Description 2", StartDate: "2022-12-06T00:00:00Z", StartTime: "17:00:00",
			EndDate: "2022-12-06T00:00:00Z", EndTime: "18:00:00", NotificationPeriod: "1",
		},
		{
			Title: "Event title 3", OwnerID: 2, Descr: "Description 3", StartDate: "2022-12-06T00:00:00Z", StartTime: "17:00:00",
			EndDate: "2022-12-06T00:00:00Z", EndTime: "18:00:00", NotificationPeriod: "1",
		},
		{
			Title: "Event title 4", OwnerID: 2, Descr: "Description 4", StartDate: "2022-12-06T00:00:00Z", StartTime: "17:00:00",
			EndDate: "2022-12-06T00:00:00Z", EndTime: "18:00:00", NotificationPeriod: "1",
		},
	}

	t.Run("add events", func(t *testing.T) {
		for ind, tEvent := range testEvents {
			id, err := s.Add(tEvent)
			require.NoError(t, err)
			require.Equal(t, ind+1, id)
		}

		allEvents, err := s.GetAll()
		require.NoError(t, err)
		require.Len(t, allEvents, len(testEvents))

		for _, savedEvent := range allEvents {
			compareSavedEventWithExpected(t, savedEvent.ID, testEvents[savedEvent.ID-1], savedEvent)
		}
	})

	t.Run("other tests", func(t *testing.T) {
		for ind, testEvent := range testEvents {
			savedEvent, err := s.Get(ind + 1)
			require.NoError(t, err)
			require.NotNil(t, savedEvent)
			expectedEventID := ind + 1
			compareSavedEventWithExpected(t, expectedEventID, testEvent, savedEvent)
		}

		err := s.Delete(1)
		require.NoError(t, err)

		allEvents, err := s.GetAll()
		require.NoError(t, err)
		require.Len(t, allEvents, len(testEvents)-1)

		nilEvent, err := s.Get(1)
		require.Nil(t, nilEvent)
		require.ErrorIs(t, err, dto.ErrEventNotFound)

		newEvent := dto.Event{
			Title: "Event title 5", OwnerID: 2, Descr: "Description 5", StartDate: "2022-12-06T00:00:00Z", StartTime: "17:00:00",
			EndDate: "2022-12-06T00:00:00Z", EndTime: "18:00:00", NotificationPeriod: "1",
		}
		newEventID, err := s.Add(newEvent)
		require.NoError(t, err)
		require.Equal(t, 5, newEventID)

		allEvents, err = s.GetAll()
		require.NoError(t, err)
		require.Len(t, allEvents, len(testEvents))

		event, err := s.Get(5)
		require.NoError(t, err)
		require.NotNil(t, event)
		compareSavedEventWithExpected(t, 5, newEvent, event)

		event.Title = "Event title 55"
		event.OwnerID = 2
		event.Descr = "Description 55"
		event.StartDate = "2023-12-06T00:00:00Z"
		event.StartTime = "17:55:00"
		event.EndDate = "2023-12-06T00:00:00Z"
		event.EndTime = "18:55:00"
		event.NotificationPeriod = "1"
		err = s.Update(*event)
		require.NoError(t, err)

		event2, err := s.Get(5)
		require.NoError(t, err)
		require.NotNil(t, event2)
		compareSavedEventWithExpected(t, 5, *event, event2)

		err = s.Update(dto.Event{ID: 6, Title: "Event title 66"})
		require.ErrorIs(t, err, dto.ErrEventNotFound)
	})
}

func compareSavedEventWithExpected(t *testing.T, expectedEventID int, expectedEvent dto.Event,
	savedEvent *dto.Event,
) {
	t.Helper()
	require.Equal(t, expectedEventID, savedEvent.ID)
	require.Equal(t, expectedEvent.Title, savedEvent.Title)
	require.Equal(t, expectedEvent.OwnerID, savedEvent.OwnerID)
	require.Equal(t, expectedEvent.Descr, savedEvent.Descr)
	require.Equal(t, expectedEvent.StartDate, savedEvent.StartDate)
	require.Equal(t, expectedEvent.StartTime, savedEvent.StartTime)
	require.Equal(t, expectedEvent.EndDate, savedEvent.EndDate)
	require.Equal(t, expectedEvent.EndTime, savedEvent.EndTime)
	require.Equal(t, expectedEvent.NotificationPeriod, savedEvent.NotificationPeriod)
}
