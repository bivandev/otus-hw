package sqlstorage

import (
	"context"
	"time"

	"github.com/devv4n/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Pool *pgxpool.Pool
}

func New(connString string) (*Storage, error) {
	Pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	return &Storage{Pool: Pool}, nil
}

func (s *Storage) Close() {
	s.Pool.Close()
}

// CreateEvent creates a new event in the database.
func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) (string, error) {
	query := `
		INSERT INTO events (title, event_datetime, duration, description, user_id, notify_before)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	var id string
	err := s.Pool.QueryRow(ctx, query,
		event.Title,
		event.EventTime,
		event.Duration,
		event.Description,
		event.UserID,
		event.NotifyBefore,
	).Scan(&id)
	return id, err
}

// UpdateEvent updates an existing event in the database.
func (s *Storage) UpdateEvent(ctx context.Context, eventID string, event storage.Event) error {
	query := `
		UPDATE events
		SET title = $1, event_datetime = $2, duration = $3, description = $4, 
		    user_id = $5, notify_before = $6, updated_at = NOW()
		WHERE id = $7`
	_, err := s.Pool.Exec(ctx, query,
		event.Title,
		event.EventTime,
		event.Duration,
		event.Description,
		event.UserID,
		event.NotifyBefore,
		eventID,
	)
	return err
}

// DeleteEvent deletes an event from the database by its ID.
func (s *Storage) DeleteEvent(ctx context.Context, eventID string) error {
	query := `DELETE FROM events WHERE id = $1`
	_, err := s.Pool.Exec(ctx, query, eventID)
	return err
}

// ListEventsForDay retrieves events for a specific day.
func (s *Storage) ListEventsForDay(ctx context.Context, date time.Time) ([]storage.Event, error) {
	query := `
		SELECT id, title, event_datetime, duration, description, user_id, notify_before
		FROM events
		WHERE DATE(event_datetime) = DATE($1)`
	return s.listEvents(ctx, query, date)
}

// ListEventsForWeek retrieves events for a specific week.
func (s *Storage) ListEventsForWeek(ctx context.Context, startOfWeek time.Time) ([]storage.Event, error) {
	query := `
		SELECT id, title, event_datetime, duration, description, user_id, notify_before
		FROM events
		WHERE event_datetime >= $1 AND event_datetime < $1 + INTERVAL '7 days'`
	return s.listEvents(ctx, query, startOfWeek)
}

// ListEventsForMonth retrieves events for a specific month.
func (s *Storage) ListEventsForMonth(ctx context.Context, startOfMonth time.Time) ([]storage.Event, error) {
	query := `
		SELECT id, title, event_datetime, duration, description, user_id, notify_before
		FROM events
		WHERE event_datetime >= $1 AND event_datetime < $1 + INTERVAL '1 month'`
	return s.listEvents(ctx, query, startOfMonth)
}

// Helper method for listing events.
func (s *Storage) listEvents(ctx context.Context, query string, param time.Time) ([]storage.Event, error) {
	rows, err := s.Pool.Query(ctx, query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []storage.Event
	for rows.Next() {
		var event storage.Event
		if err = rows.Scan(
			&event.ID,
			&event.Title,
			&event.EventTime,
			&event.Duration,
			&event.Description,
			&event.UserID,
			&event.NotifyBefore,
		); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, rows.Err()
}
