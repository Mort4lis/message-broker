package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type statusCodeRecorder struct {
	http.ResponseWriter
	StatusCode int
}

func (rec *statusCodeRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func Log(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rec := &statusCodeRecorder{ResponseWriter: w}

		begin := time.Now()
		next.ServeHTTP(rec, req)

		logger.With(
			slog.String("path", req.URL.EscapedPath()),
			slog.String("user_agent", req.UserAgent()),
			slog.String("method", req.Method),
			slog.String("ip", req.RemoteAddr),
			slog.String("host", req.Host),
			slog.Int64("request_size", req.ContentLength),
			slog.String("duration", time.Since(begin).String()),
			slog.Int("status_code", rec.StatusCode),
		).Info("Request handled")
	})
}
