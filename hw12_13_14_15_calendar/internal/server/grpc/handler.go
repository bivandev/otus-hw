package grpc

import (
	"context"
	"errors"
	api "github.com/devv4n/otus-hw/hw12_13_14_15_calendar/pkg/calendar-api"
	"google.golang.org/protobuf/types/known/emptypb"
)

var ErrInvalidDate = errors.New("invalid date format")

// CreateEvent handles the gRPC request to create a new event.
func (s *Server) CreateEvent(
	ctx context.Context,
	req *api.CreateEventRequest,
) (*api.CreateEventResponse, error) {
	eventID, err := s.app.CreateEvent(ctx, fromAPIEvents(req.Event))
	if err != nil {
		return nil, err
	}

	return &api.CreateEventResponse{
		Id: eventID,
	}, nil
}

// UpdateEvent handles the gRPC request to update an existing event.
func (s *Server) UpdateEvent(
	ctx context.Context,
	req *api.UpdateEventRequest,
) (*emptypb.Empty, error) {
	err := s.app.UpdateEvent(ctx, fromAPIEvents(req.Event))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// DeleteEvent handles the gRPC request to delete an event.
func (s *Server) DeleteEvent(
	ctx context.Context,
	req *api.DeleteEventRequest,
) (*api.DeleteEventResponse, error) {
	err := s.app.DeleteEvent(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &api.DeleteEventResponse{}, nil
}

// ListEventsForDay handles the gRPC request to list events for a specific day.
func (s *Server) ListEventsForDay(
	ctx context.Context,
	req *api.ListEventsForDayRequest,
) (*api.ListEventsResponse, error) {
	date := req.GetDate().AsTime()

	events, err := s.app.ListEventsForDay(ctx, date)
	if err != nil {
		return nil, err
	}

	return &api.ListEventsResponse{
		Events: toAPIEvents(events),
	}, nil
}

// ListEventsForWeek handles the gRPC request to list events for a specific week.
func (s *Server) ListEventsForWeek(
	ctx context.Context,
	req *api.ListEventsForWeekRequest,
) (*api.ListEventsResponse, error) {
	startOfWeek := req.GetStartOfWeek().AsTime()

	events, err := s.app.ListEventsForWeek(ctx, startOfWeek)
	if err != nil {
		return nil, err
	}

	return &api.ListEventsResponse{
		Events: toAPIEvents(events),
	}, nil
}

// ListEventsForMonth handles the gRPC request to list events for a specific month.
func (s *Server) ListEventsForMonth(
	ctx context.Context,
	req *api.ListEventsForMonthRequest,
) (*api.ListEventsResponse, error) {
	startOfMonth := req.GetStartOfMonth().AsTime()

	events, err := s.app.ListEventsForMonth(ctx, startOfMonth)
	if err != nil {
		return nil, err
	}

	return &api.ListEventsResponse{
		Events: toAPIEvents(events),
	}, nil
}
