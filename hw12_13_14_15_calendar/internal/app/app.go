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
	UpdateEvent(ctx context.Context, eventID string, event storage.Event) error

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

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	_, err := a.storage.CreateEvent(ctx, storage.Event{ID: id, Title: title})
	if err != nil {
		return err
	}

	return nil
}
