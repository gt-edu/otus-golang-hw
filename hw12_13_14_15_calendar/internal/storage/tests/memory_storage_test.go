package tests

import (
	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	"testing"

	memorystorage "github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/memory"
)

func TestMemoryStorage(t *testing.T) {
	factory := func() storage.EventStorage {
		return memorystorage.New()
	}

	EventsCommonTest(t, factory)

	EventsCommonConcurrencyTest(t, factory)
}
