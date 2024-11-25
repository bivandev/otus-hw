package grpc

import (
	"context"
	"fmt"
	"net/http"

	api "github.com/devv4n/otus-hw/hw12_13_14_15_calendar/pkg/calendar-api"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ServeGatewayAPI starts gRPC server to serve grpc-gateway API.
func (s *Server) ServeGatewayAPI(ctx context.Context, errCh chan error) {
	mux := runtime.NewServeMux()

	srv := &http.Server{
		Addr:              s.cfg.REST,
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
	}
	defer srv.Close()

	conn, err := grpc.NewClient(
		s.cfg.GRPC,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		errCh <- fmt.Errorf("ServeGatewayAPI: grpc.DialContext: %w", err)
	}

	defer conn.Close()

	if err = api.RegisterCalendarServiceHandler(ctx, mux, conn); err != nil {
		errCh <- fmt.Errorf("ServeGatewayAPI: apiv1.RegisterDispatchSecurityServiceHandler: %w", err)
	}

	if err = srv.ListenAndServe(); err != nil {
		errCh <- fmt.Errorf("ServeGatewayAPI: mux.ListenAndServe: %w", err)
	}
}
