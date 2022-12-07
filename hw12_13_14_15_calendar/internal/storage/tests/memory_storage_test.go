package tests

import (
	memorystorage "github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/memory"
	"testing"
)

func TestMemoryStorage(t *testing.T) {
	s := memorystorage.New()
	EventsCommonTest(t, s)
}
