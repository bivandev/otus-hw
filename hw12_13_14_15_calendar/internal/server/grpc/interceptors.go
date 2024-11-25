package grpc

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor is a gRPC unary interceptor for logging request details.
func LoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	p, ok := peer.FromContext(ctx)

	ip := p.Addr.String()
	if !ok {
		ip = "unknown"
	}

	start := time.Now()

	resp, err = handler(ctx, req)

	code := status.Code(err)

	values := []any{
		"ip", ip,
		"method", info.FullMethod,
		"status_code", code.String(),
		"latency_ms", float64(time.Since(start).Microseconds()) / 1000,
	}

	if err != nil {
		slog.ErrorContext(ctx, "gRPC request error", append(values, "error", err)...)
	} else {
		slog.InfoContext(ctx, "gRPC request handled", values...)
	}

	return resp, err
}
