package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/devv4n/otus-hw/hw12_13_14_15_calendar/internal/storage"
	calendar_api "github.com/devv4n/otus-hw/hw12_13_14_15_calendar/pkg/calendar-api"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MockApplication struct {
	mock.Mock
}

func (m *MockApplication) CreateEvent(ctx context.Context, event storage.Event) (string, error) {
	args := m.Called(ctx, event)
	return args.String(0), args.Error(1)
}

func (m *MockApplication) UpdateEvent(ctx context.Context, event storage.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockApplication) DeleteEvent(ctx context.Context, eventID string) error {
	args := m.Called(ctx, eventID)
	return args.Error(0)
}

func (m *MockApplication) ListEventsForDay(ctx context.Context, date time.Time) ([]storage.Event, error) {
	args := m.Called(ctx, date)
	return args.Get(0).([]storage.Event), args.Error(1)
}

func (m *MockApplication) ListEventsForWeek(ctx context.Context, startOfWeek time.Time) ([]storage.Event, error) {
	args := m.Called(ctx, startOfWeek)
	return args.Get(0).([]storage.Event), args.Error(1)
}

func (m *MockApplication) ListEventsForMonth(ctx context.Context, startOfMonth time.Time) ([]storage.Event, error) {
	args := m.Called(ctx, startOfMonth)
	return args.Get(0).([]storage.Event), args.Error(1)
}

func TestCreateEvent(t *testing.T) {
	mockApp := new(MockApplication)
	server := &Server{app: mockApp}

	eventRequest := &calendar_api.CreateEventRequest{
		Event: &calendar_api.Event{
			Title: "Test Event",
		},
	}

	mockApp.On("CreateEvent", mock.Anything, mock.AnythingOfType("storage.Event")).Return("eventID123", nil)

	resp, err := server.CreateEvent(context.Background(), eventRequest)

	assert.NoError(t, err)
	assert.Equal(t, "eventID123", resp.Id)

	mockApp.AssertExpectations(t)
}

func TestUpdateEvent(t *testing.T) {
	mockApp := new(MockApplication)
	server := &Server{app: mockApp}

	eventRequest := &calendar_api.UpdateEventRequest{
		Event: &calendar_api.Event{
			Id:    "eventID123",
			Title: "Updated Event",
		},
	}

	mockApp.On("UpdateEvent", mock.Anything, mock.AnythingOfType("storage.Event")).Return(nil)

	resp, err := server.UpdateEvent(context.Background(), eventRequest)

	assert.NoError(t, err)
	assert.IsType(t, &emptypb.Empty{}, resp)

	mockApp.AssertExpectations(t)
}

func TestDeleteEvent(t *testing.T) {
	mockApp := new(MockApplication)
	server := &Server{app: mockApp}

	deleteRequest := &calendar_api.DeleteEventRequest{
		Id: "eventID123",
	}

	mockApp.On("DeleteEvent", mock.Anything, "eventID123").Return(nil)

	resp, err := server.DeleteEvent(context.Background(), deleteRequest)

	assert.NoError(t, err)
	assert.IsType(t, &calendar_api.DeleteEventResponse{}, resp)

	mockApp.AssertExpectations(t)
}

func TestListEventsForDay(t *testing.T) {
	mockApp := new(MockApplication)
	server := &Server{app: mockApp}

	date := time.Now()

	listRequest := &calendar_api.ListEventsForDayRequest{
		Date: &timestamp.Timestamp{Seconds: date.Unix()},
	}

	mockApp.On("ListEventsForDay", mock.Anything, mock.MatchedBy(func(t time.Time) bool {
		return t.Location() == time.UTC
	})).Return([]storage.Event{
		{ID: "eventID123", Title: "Test Event"},
	}, nil)

	resp, err := server.ListEventsForDay(context.Background(), listRequest)

	assert.NoError(t, err)
	assert.Len(t, resp.Events, 1)
	assert.Equal(t, "eventID123", resp.Events[0].Id)

	mockApp.AssertExpectations(t)
}

func TestListEventsForWeek(t *testing.T) {
	mockApp := new(MockApplication)
	server := &Server{app: mockApp}

	startOfWeek := time.Now().AddDate(0, 0, -7)

	listRequest := &calendar_api.ListEventsForWeekRequest{
		StartOfWeek: &timestamp.Timestamp{Seconds: startOfWeek.Unix()},
	}

	mockApp.On("ListEventsForWeek", mock.Anything, mock.MatchedBy(func(t time.Time) bool {
		return t.Location() == time.UTC
	})).Return([]storage.Event{
		{ID: "eventID123", Title: "Test Event"},
	}, nil)

	resp, err := server.ListEventsForWeek(context.Background(), listRequest)

	assert.NoError(t, err)
	assert.Len(t, resp.Events, 1)
	assert.Equal(t, "eventID123", resp.Events[0].Id)

	mockApp.AssertExpectations(t)
}

func TestListEventsForMonth(t *testing.T) {
	mockApp := new(MockApplication)
	server := &Server{app: mockApp}

	startOfMonth := time.Now().AddDate(0, 0, -30)

	listRequest := &calendar_api.ListEventsForMonthRequest{
		StartOfMonth: &timestamp.Timestamp{Seconds: startOfMonth.Unix()},
	}

	mockApp.On("ListEventsForMonth", mock.Anything, mock.MatchedBy(func(t time.Time) bool {
		return t.Location() == time.UTC
	})).Return([]storage.Event{
		{ID: "eventID123", Title: "Test Event"},
	}, nil)

	resp, err := server.ListEventsForMonth(context.Background(), listRequest)

	assert.NoError(t, err)
	assert.Len(t, resp.Events, 1)
	assert.Equal(t, "eventID123", resp.Events[0].Id)

	mockApp.AssertExpectations(t)
}
