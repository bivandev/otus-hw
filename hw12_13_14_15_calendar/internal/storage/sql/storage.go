package sqlstorage

import (
	"context"
	"time"

	"github.com/devv4n/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
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
	if event.ID == "" {
		event.ID = uuid.NewString()
	}

	query := `
		INSERT INTO events (id, title, event_datetime, duration, description, user_id, notify_before)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`
	var id string
	err := s.Pool.QueryRow(ctx, query,
		event.ID,
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
func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
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
		event.ID,
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
	sqlDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	query := `
        SELECT id, title, event_datetime, duration, description, user_id, notify_before
        FROM events
        WHERE event_datetime::date = $1::date`
	return s.listEvents(ctx, query, sqlDate)
}

// ListEventsForWeek retrieves events for a specific week.
func (s *Storage) ListEventsForWeek(ctx context.Context, startOfWeek time.Time) ([]storage.Event, error) {
	sqlDate := time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())

	query := `
		SELECT id, title, event_datetime, duration, description, user_id, notify_before
		FROM events
		WHERE event_datetime::date >= $1::date AND event_datetime::date < $1::date + INTERVAL '7 days'`
	return s.listEvents(ctx, query, sqlDate)
}

// ListEventsForMonth retrieves events for a specific month.
func (s *Storage) ListEventsForMonth(ctx context.Context, startOfMonth time.Time) ([]storage.Event, error) {
	sqlDate := time.Date(
		startOfMonth.Year(),
		startOfMonth.Month(),
		startOfMonth.Day(),
		0,
		0,
		0,
		0,
		startOfMonth.Location(),
	)

	query := `
		SELECT id, title, event_datetime, duration, description, user_id, notify_before
		FROM events
		WHERE event_datetime::date >= $1::date AND event_datetime::date < $1::date + INTERVAL '1 month'`
	return s.listEvents(ctx, query, sqlDate)
}

// Helper method for listing events.
func (s *Storage) listEvents(ctx context.Context, query string, args ...any) ([]storage.Event, error) {
	rows, err := s.Pool.Query(ctx, query, args...)
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

func (s *Storage) GetEventsForNotification(ctx context.Context) ([]storage.Event, error) {
	const query = `
		SELECT id, title, event_datetime, user_id
		FROM events
		WHERE event_datetime <= NOW() + INTERVAL '1 hour'
	`

	rows, err := s.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []storage.Event
	for rows.Next() {
		var event storage.Event
		if err = rows.Scan(&event.ID, &event.Title, &event.EventTime, &event.UserID); err != nil {
			return nil, err
		}

		events = append(events, event)
	}
	return events, nil
}

func (s *Storage) CleanOldEvents(ctx context.Context) error {
	const query = `
		DELETE FROM events
		WHERE event_datetime < NOW() - INTERVAL '1 year'
	`

	_, err := s.Pool.Exec(ctx, query)

	return err
}
