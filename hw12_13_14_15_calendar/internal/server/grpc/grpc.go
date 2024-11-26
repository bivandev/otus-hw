package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/devv4n/otus-hw/hw12_13_14_15_calendar/internal/storage"
	api "github.com/devv4n/otus-hw/hw12_13_14_15_calendar/pkg/calendar-api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const readHeaderTimeout = 10 * time.Second

type Server struct {
	app Application
	cfg Config

	api.UnimplementedCalendarServiceServer
}

type Application interface {
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

type Config struct {
	GRPC string
	REST string
}

func New(app Application, cfg Config) *Server {
	return &Server{
		app: app,
		cfg: cfg,
	}
}

// ServeUserAPI starts gRPC server to serve user API.
func (s *Server) ServeUserAPI(errCh chan error) {
	listener, err := net.Listen("tcp", s.cfg.GRPC)
	if err != nil {
		errCh <- fmt.Errorf("error ServeUserAPI: net.Listen: %w", err)
	}

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(LoggingInterceptor),
	)

	defer srv.GracefulStop()

	reflection.Register(srv)

	api.RegisterCalendarServiceServer(srv, s)

	if err = srv.Serve(listener); err != nil {
		errCh <- fmt.Errorf("error ServeUserAPI: net.Listen: %w", err)
	}
}
