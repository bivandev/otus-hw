package grpc

import (
	"database/sql"

	"github.com/devv4n/otus-hw/hw12_13_14_15_calendar/internal/storage"
	api "github.com/devv4n/otus-hw/hw12_13_14_15_calendar/pkg/calendar-api"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toAPIEvents(events []storage.Event) []*api.Event {
	apiEvents := make([]*api.Event, len(events))
	for i, e := range events {
		apiEvents[i] = toAPIEvent(e)
	}
	return apiEvents
}

func toAPIEvent(e storage.Event) *api.Event {
	return &api.Event{
		Id:           e.ID,
		Title:        e.Title,
		EventTime:    timestamppb.New(e.EventTime),
		Duration:     e.Duration.Int64,
		Description:  &e.Description.String,
		UserId:       e.UserID,
		NotifyBefore: &e.NotifyBefore.Int64,
		CreatedAt:    timestamppb.New(e.CreatedAt),
		UpdatedAt:    timestamppb.New(e.UpdatedAt),
	}
}

func fromAPIEvents(event *api.Event) storage.Event {
	return storage.Event{
		ID:        event.Id,
		Title:     event.Title,
		EventTime: event.EventTime.AsTime(),
		Duration: sql.NullInt64{
			Int64: event.Duration,
			Valid: event.Duration != 0,
		},
		Description: sql.NullString{
			String: getStringPointerValue(event.Description),
			Valid:  event.Description != nil,
		},
		UserID: event.UserId,
		NotifyBefore: sql.NullInt64{
			Int64: getInt64PointerValue(event.NotifyBefore),
			Valid: event.NotifyBefore != nil,
		},
	}
}

func getStringPointerValue(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func getInt64PointerValue(i *int64) int64 {
	if i != nil {
		return *i
	}
	return 0
}
