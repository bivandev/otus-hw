package rabbitmq

import "time"

type Notification struct {
	EventID   string    `json:"eventId"`
	Title     string    `json:"title"`
	StartTime time.Time `json:"startTime"`
	OwnerID   string    `json:"ownerId"`
}
