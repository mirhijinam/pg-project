package logger

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	envLocal = "local"
	envDev   = "develop"
	envProd  = "prod"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (lw *loggingResponseWriter) WriteHeader(status int) {
	lw.status = status
	lw.ResponseWriter.WriteHeader(status)
}

func (lw *loggingResponseWriter) Write(p []byte) (int, error) {
	n, err := lw.ResponseWriter.Write(p)
	lw.length += n
	return n, err
}

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := log.With(
				slog.String("component", "middleware/logger"),
			)

			lw := &loggingResponseWriter{ResponseWriter: w}

			t1 := time.Now()

			defer func() {
				entry := log.With(
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.String("remote_addr", r.RemoteAddr),
					slog.String("user_agent", r.UserAgent()),
					slog.Int("status", lw.status),
					slog.Int("length", lw.length),
					slog.String("duration", time.Since(t1).String()),
				)
				entry.Info("request completed")
			}()

			next.ServeHTTP(lw, r)
		})
	}
}

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
