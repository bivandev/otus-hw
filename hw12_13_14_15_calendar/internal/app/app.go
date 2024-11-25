package app

import (
	"context"
	"time"

	"github.com/devv4n/otus-hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	storage Storage
}

type Storage interface {
	Close()

	// CreateEvent creates a new event.
	CreateEvent(ctx context.Context, event storage.Event) (string, error)

	// UpdateEvent updates an existing event by its ID.
	UpdateEvent(ctx context.Context, event storage.Event) error

	// DeleteEvent deletes an event by its ID.
	DeleteEvent(ctx context.Context, eventID string) error

	// ListEventsForDay returns a list of events for a specific day.
	ListEventsForDay(ctx context.Context, date time.Time) ([]storage.Event, error)

	// ListEventsForWeek returns a list of events for a specific week.
	ListEventsForWeek(ctx context.Context, startOfWeek time.Time) ([]storage.Event, error)

	// ListEventsForMonth returns a list of events for a specific month.
	ListEventsForMonth(ctx context.Context, startOfMonth time.Time) ([]storage.Event, error)
}

func New(storage Storage) *App {
	return &App{
		storage: storage,
	}
}

// CreateEvent creates a new event.
func (a *App) CreateEvent(ctx context.Context, event storage.Event) (string, error) {
	return a.storage.CreateEvent(ctx, event)
}

// UpdateEvent updates an existing event by its ID.
func (a *App) UpdateEvent(ctx context.Context, event storage.Event) error {
	return a.storage.UpdateEvent(ctx, event)
}

// DeleteEvent deletes an event by its ID.
func (a *App) DeleteEvent(ctx context.Context, eventID string) error {
	return a.storage.DeleteEvent(ctx, eventID)
}

// ListEventsForDay returns a list of events for a specific day.
func (a *App) ListEventsForDay(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.storage.ListEventsForDay(ctx, date)
}

// ListEventsForWeek returns a list of events for a specific week.
func (a *App) ListEventsForWeek(ctx context.Context, startOfWeek time.Time) ([]storage.Event, error) {
	return a.storage.ListEventsForWeek(ctx, startOfWeek)
}

// ListEventsForMonth returns a list of events for a specific month.
func (a *App) ListEventsForMonth(ctx context.Context, startOfMonth time.Time) ([]storage.Event, error) {
	return a.storage.ListEventsForMonth(ctx, startOfMonth)
}
