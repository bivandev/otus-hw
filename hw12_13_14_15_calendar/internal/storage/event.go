package storage

import (
	"database/sql"
	"time"
)

type Event struct {
	ID           string         `json:"id"`
	Title        string         `json:"title"`
	EventTime    time.Time      `json:"eventTime"`
	Duration     sql.NullInt64  `json:"duration"`
	Description  sql.NullString `json:"description"`
	UserID       string         `json:"userId"`
	NotifyBefore sql.NullInt64  `json:"notifyBefore"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
}
