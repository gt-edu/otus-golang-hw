package storage

type Event struct {
	ID                 int
	Title              string
	OwnerID            int `db:"owner_id"`
	Descr              string
	StartDate          string `db:"start_date"`
	StartTime          string `db:"start_time"`
	EndDate            string `db:"end_date"`
	EndTime            string `db:"end_time"`
	NotificationPeriod string `db:"notification_period"`
}
