package middleware

import (
	"context"
	"net/http"
	"time"
)

func Timeout(timeout time.Duration, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(req.Context(), timeout)
		defer cancel()

		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
