package tests

import (
	"testing"

	memorystorage "github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/memory"
)

func TestMemoryStorage(t *testing.T) {
	s := memorystorage.New()
	EventsCommonTest(t, s)
}
