package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/devv4n/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestCreateEvent(t *testing.T) {
	ctx := context.Background()
	st := New()

	event := storage.Event{
		Title:     "Test Event",
		EventTime: time.Now(),
	}

	id, err := st.CreateEvent(ctx, event)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)

	events, _ := st.ListEventsForDay(ctx, event.EventTime)
	assert.Len(t, events, 1)
	assert.Equal(t, "Test Event", events[0].Title)
}

func TestUpdateEvent(t *testing.T) {
	ctx := context.Background()
	st := New()

	event := storage.Event{
		Title:     "Initial Event",
		EventTime: time.Now(),
	}

	id, err := st.CreateEvent(ctx, event)
	assert.NoError(t, err)

	updatedEvent := storage.Event{
		ID:        id,
		Title:     "Updated Event",
		EventTime: event.EventTime.Add(2 * time.Hour),
	}
	err = st.UpdateEvent(ctx, updatedEvent)
	assert.NoError(t, err)

	events, _ := st.ListEventsForDay(ctx, updatedEvent.EventTime)
	assert.Len(t, events, 1)
	assert.Equal(t, "Updated Event", events[0].Title)
	assert.Equal(t, id, events[0].ID)
}

func TestDeleteEvent(t *testing.T) {
	ctx := context.Background()
	st := New()

	event := storage.Event{
		Title:     "Test Event",
		EventTime: time.Now(),
	}

	id, err := st.CreateEvent(ctx, event)
	assert.NoError(t, err)

	err = st.DeleteEvent(ctx, id)
	assert.NoError(t, err)

	events, _ := st.ListEventsForDay(ctx, event.EventTime)
	assert.Empty(t, events)
}

func TestListEventsForDay(t *testing.T) {
	ctx := context.Background()
	st := New()

	date := time.Now()
	events := []storage.Event{
		{Title: "Event 1", EventTime: date},
		{Title: "Event 2", EventTime: date.Add(1 * time.Hour)},
		{Title: "Event 3", EventTime: date.AddDate(0, 0, 1)},
	}

	for _, e := range events {
		_, _ = st.CreateEvent(ctx, e)
	}

	dayEvents, _ := st.ListEventsForDay(ctx, date)
	assert.Len(t, dayEvents, 2)
}

func TestListEventsForWeek(t *testing.T) {
	ctx := context.Background()
	st := New()

	year, week := time.Now().ISOWeek()
	startOfWeek := time.Date(year, 0, 0, 0, 0, 0, 0, time.UTC).AddDate(0, 0, (week-1)*7)
	startOfWeek = startOfWeek.AddDate(0, 0, -int(startOfWeek.Weekday()-time.Monday))

	events := []storage.Event{
		{Title: "Event 1", EventTime: startOfWeek},
		{Title: "Event 2", EventTime: startOfWeek.Add(2 * time.Hour)},
		{Title: "Event 3", EventTime: startOfWeek.AddDate(0, 0, 7)},
	}

	for _, e := range events {
		_, _ = st.CreateEvent(ctx, e)
	}

	weekEvents, _ := st.ListEventsForWeek(ctx, startOfWeek)
	assert.Len(t, weekEvents, 2)
}

func TestListEventsForMonth(t *testing.T) {
	ctx := context.Background()
	st := New()

	startOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)

	events := []storage.Event{
		{Title: "Event 1", EventTime: startOfMonth},
		{Title: "Event 2", EventTime: startOfMonth.AddDate(0, 0, 15)},
		{Title: "Event 3", EventTime: startOfMonth.AddDate(0, 1, 0)},
	}

	for _, e := range events {
		_, _ = st.CreateEvent(ctx, e)
	}

	monthEvents, _ := st.ListEventsForMonth(ctx, startOfMonth)
	assert.Len(t, monthEvents, 2)
}

func TestGetEventsForNotification(t *testing.T) {
	st := New()
	now := time.Now()

	st.CreateEvent(context.TODO(), storage.Event{
		Title:     "Past Event",
		EventTime: now.Add(-1 * time.Hour),
	})
	st.CreateEvent(context.TODO(), storage.Event{
		Title:     "Upcoming Event",
		EventTime: now.Add(30 * time.Minute),
	})
	st.CreateEvent(context.TODO(), storage.Event{
		Title:     "Far Future Event",
		EventTime: now.Add(2 * time.Hour),
	})

	events, err := st.GetEventsForNotification(context.TODO())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(events) != 1 || events[0].Title != "Upcoming Event" {
		t.Fatalf("unexpected events: %+v", events)
	}
}

func TestCleanOldEvents(t *testing.T) {
	st := New()
	now := time.Now()

	st.CreateEvent(context.TODO(), storage.Event{
		Title:     "Old Event",
		EventTime: now.AddDate(-2, 0, 0),
	})
	st.CreateEvent(context.TODO(), storage.Event{
		Title:     "Recent Event",
		EventTime: now,
	})

	err := st.CleanOldEvents(context.TODO())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	events, _ := st.listEventsByFilter(func(_ storage.Event) bool {
		return true
	})
	if len(events) != 1 || events[0].Title != "Recent Event" {
		t.Fatalf("unexpected events: %+v", events)
	}
}
