package storage

type Event struct {
	ID                 int
	Title              string
	OwnerId            int
	Descr              string
	StartDate          string
	StartTime          string
	EndDate            string
	EndTime            string
	NotificationPeriod string
}
