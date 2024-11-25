package internalhttp

import (
	"log/slog"
	"net"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		rw := &responseWrapper{
			ResponseWriter: w,
		}

		start := time.Now()

		defer func() {
			status := rw.status

			values := []any{
				"ip", ip,
				"method", r.Method,
				"path", r.URL.Path,
				"proto", r.Proto,
				"status_code", rw.status,
				"response_size", rw.size,
				"user_agent", r.UserAgent(),
				"latency_ms", float64(time.Since(start).Microseconds()) / 1000,
			}

			if status >= http.StatusBadRequest && status <= http.StatusNetworkAuthenticationRequired {
				slog.ErrorContext(ctx, http.StatusText(status), values...)
			} else {
				slog.InfoContext(ctx, "", values...)
			}
		}()

		next.ServeHTTP(rw, r)
	})
}

type responseWrapper struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWrapper) WriteHeader(code int) {
	rw.status = code

	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWrapper) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}

	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func (rw *responseWrapper) Flush() {
	rw.ResponseWriter.(http.Flusher).Flush()
}
