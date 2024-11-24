package memorystorage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/devv4n/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type Storage struct {
	mu     sync.RWMutex
	events map[string]storage.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[string]storage.Event),
	}
}

func (s *Storage) Close() {}

// CreateEvent creates a new event in the in-memory store.
func (s *Storage) CreateEvent(_ context.Context, event storage.Event) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	event.ID = uuid.NewString()
	s.events[event.ID] = event
	return event.ID, nil
}

// UpdateEvent updates an existing event in the in-memory store.
func (s *Storage) UpdateEvent(_ context.Context, eventID string, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[eventID]; !exists {
		return errors.New("event not found")
	}
	event.ID = eventID
	s.events[eventID] = event
	return nil
}

// DeleteEvent deletes an event from the in-memory store by its ID.
func (s *Storage) DeleteEvent(_ context.Context, eventID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[eventID]; !exists {
		return errors.New("event not found")
	}
	delete(s.events, eventID)
	return nil
}

// ListEventsForDay retrieves events for a specific day.
func (s *Storage) ListEventsForDay(_ context.Context, date time.Time) ([]storage.Event, error) {
	return s.listEventsByFilter(func(event storage.Event) bool {
		return sameDay(event.EventTime, date)
	})
}

// ListEventsForWeek retrieves events for a specific week.
func (s *Storage) ListEventsForWeek(_ context.Context, startOfWeek time.Time) ([]storage.Event, error) {
	endOfWeek := startOfWeek.AddDate(0, 0, 7)
	return s.listEventsByFilter(func(event storage.Event) bool {
		return !event.EventTime.Before(startOfWeek) && event.EventTime.Before(endOfWeek)
	})
}

// ListEventsForMonth retrieves events for a specific month.
func (s *Storage) ListEventsForMonth(_ context.Context, startOfMonth time.Time) ([]storage.Event, error) {
	endOfMonth := startOfMonth.AddDate(0, 1, 0) // До конца месяца
	return s.listEventsByFilter(func(event storage.Event) bool {
		return !event.EventTime.Before(startOfMonth) && event.EventTime.Before(endOfMonth)
	})
}

// Helper method for filtering events by a condition.
func (s *Storage) listEventsByFilter(filter func(event storage.Event) bool) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []storage.Event
	for _, event := range s.events {
		if filter(event) {
			results = append(results, event)
		}
	}
	return results, nil
}

// Helper function to check if two dates are on the same day.
func sameDay(a, b time.Time) bool {
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
